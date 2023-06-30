package mgrpc

import (
	"casorder/db"
	"casorder/db/models"
	"casorder/utils"
	"casorder/utils/types"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net"
	"time"

	"casorder/utils/logging"

	mailer "casorder/taskmanager/grpc/build/mails"
	od "casorder/taskmanager/grpc/build/order"
	us "casorder/taskmanager/grpc/build/user"
	"google.golang.org/grpc"
)

var LOG = logging.GetLogger()

func getMiqGrpc() string {
	host := viper.GetString("miq_grpc.host")
	port := viper.GetString("miq_grpc.port")
	miqGrpc := fmt.Sprintf("%v:%v", host, port)
	return miqGrpc
}

func getMailGrpc() string {
	host := viper.GetString("mail_grpc.host")
	port := viper.GetString("mail_grpc.port")
	mailGrpc := fmt.Sprintf("%v:%v", host, port)
	return mailGrpc
}

// server is used to implement Order Server.
type server struct {
	od.Order
	od.OrderServiceServer
}

func (s *server) GetOrder(ctx context.Context, request *od.OrderRequest) (*od.Order, error) {
	LOG.Info("Order information: ", request.OrderId)
	DB := db.GetDB()
	orderId := request.OrderId
	orderIdx := request.OrderIdx

	var order models.Order

	err := DB.Model(&order).Preload(clause.Associations).First(&order, "id = ?", orderId).Error
	if err != nil {
		LOG.Errorf("Error getting order: %v", err)
		return nil, err
	} else {
		var ops []*models.OrderProduct
		if orderIdx != 0 {
			if err := DB.Model(&ops).Preload(clause.Associations).Preload("Product.Unit").Find(&ops, "order_id = ? AND idx = ?", orderId, orderIdx).Error; err != nil {
				LOG.Errorf("Error getting order: %v", err)
				return nil, err
			}
		} else {
			if err := DB.Model(&ops).Preload(clause.Associations).Preload("Product.Unit").Find(&ops, "order_id = ?", orderId).Error; err != nil {
				LOG.Errorf("Error getting order: %v", err)
				return nil, err
			}
		}
		var computes []*od.Compute
		var compute *od.Compute

		idx := 0

		if len(ops) > 0 {
			idx = ops[0].Idx
		}

		if idx > 0 {
			compute = &od.Compute{Idx: int32(idx)}
		}

		for _, op := range ops {
			if op.Idx != idx {
				idx = op.Idx
				computes = append(computes, compute)
				compute = &od.Compute{Idx: int32(idx)}
			}

			if op.IsPackage {
				var pps []*models.PackageProduct
				if err := DB.Model(&pps).Preload(clause.Associations).Preload("Product.Unit").Find(&pps, "package_id = ?", op.ProductID).Error; err != nil {
					LOG.Errorf("Error getting order: %v", err)
					return nil, err
				}
				for _, pp := range pps {
					product := od.Product{
						Id:       int32(pp.ID),
						Cn:       pp.Product.CN,
						Quantity: int32(pp.Quantity),
					}
					compute.Products = append(compute.Products, &product)
				}
			} else {
				product := od.Product{
					Id: int32(op.ID),
					Cn: op.Product.CN,
					Quantity: int32(op.Quantity),
					Type: op.Product.Type.Value(),
					Version: "v1",
				}
				compute.Products = append(compute.Products, &product)
			}
		}

		computes = append(computes, compute)

		var result *od.Order
		result = &od.Order {
			RegionId: int32(order.RegionID),
			Duration: int32(order.Duration),
			ApprovalStep: int32(order.ApprovalStep),
			ContractNo: order.ContractCode,
			StartAt: timestamppb.New(order.StartAt),
			EndAt: timestamppb.New(order.EndAt),
			UserId: int64(order.CustomerID),
			Computes: computes,
		}
		return result, nil
	}
}

type GrpcServer struct {
	host string
	port string
	grpcServer *grpc.Server
}

func New(host string, port string) GrpcServer {
	gs := GrpcServer{
		host: host,
		port: port,
	}
	gs.grpcServer = grpc.NewServer()

	return gs
}

func (gs *GrpcServer) Start() error {
	grpcAddr := net.JoinHostPort(gs.host, gs.port)
	LOG.Info("Starting Grpc server: ", grpcAddr)
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return err
	}

	od.RegisterOrderServiceServer(gs.grpcServer, &server{})

	if err := gs.grpcServer.Serve(listener); err != nil {
		LOG.Error("failed to serve: ", err)
		return err
	}
	return nil
}

func VerifyToken(token string) (types.JSON, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, getMiqGrpc(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		LOG.Error("Error connecting: ", err)
		return nil, err
	}
	defer conn.Close()
	c := us.NewUserServiceClient(conn)

	resp, err := c.VerifyToken(ctx, &us.TokenRequest{Token: token})

	if err != nil {
		fmt.Printf("Token verification failed: %v\n", err)
		LOG.Error("Token verification failed: ", err)
		return nil, err
	}

	userData := types.JSON {
		"id": resp.Id,
		"userid": resp.Userid,
		"role": resp.Role,
	}

	return userData, err
}

func SendOrderMail(order *models.Order, items [][]models.OrderProduct, data types.JSON) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, getMailGrpc(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		LOG.Error("Error connecting: ", err)
		return false, err
	}
	defer conn.Close()
	c := mailer.NewMailServiceClient(conn)

	computeAttr := []string{"vcpu", "ram"}
	diskAttr := []string{"disk", "root_disk", "data_disk"}
	var computes []*mailer.Compute

	for _, item := range items {
		var compute mailer.Compute
		for _, op := range item {

			if utils.Contains(computeAttr, op.Product.CN) {
				if op.Product.CN == "vcpu" {
					compute.Cpu = int32(op.Quantity)
				}
				if op.Product.CN == "ram" {
					compute.Ram = int32(op.Quantity)
				}
			} else if utils.Contains(diskAttr, op.Product.CN) {
				compute.Disk += int32(op.Quantity)
			}
		}
		compute.Email = order.OrderDtl.Email
		compute.OrderId = int32(order.ID)
		compute.OrderCode = order.Code
		compute.OrderCreator = data["staff_username"].(string)

		computes = append(computes, &compute)
	}

	customer := data["customer"].(types.JSON)

	username := ""
	if customer["username"] != nil {
		username = customer["username"].(string)
	}

	password := ""
	if customer["password"] != nil {
		password = customer["password"].(string)
	}

	orderInfo := mailer.Order{
		ContractNo: order.ContractCode,
		Code: order.Code,
		Username: username,
		Password: password,
		CustomerName: order.OrderDtl.Name,
		CustomerEmail: order.OrderDtl.Email,
		Computes: computes,
	}

	reply, err := c.SendOrderInfo(ctx, &orderInfo)
	if err != nil {
		fmt.Printf("Error sending order info: %v\n", err)
		LOG.Error("Error sending order info: ", err)
		return false, err
	}

	return reply.Status, nil
}

func GetApproval(order *models.Order, DB *gorm.DB) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, getMiqGrpc(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		LOG.Error("Error connecting: ", err)
		return false, err
	}
	defer conn.Close()
	c := od.NewOrderServiceClient(conn)

	var computes []*od.Compute
	ops := order.OrderProducts

	for _, op := range ops {
		DB.Model(&op).Preload(clause.Associations).First(&op)
		compute := od.Compute{Idx: int32(op.Idx)}
		product := od.Product{
				Id: int32(op.ID),
				Cn: op.Product.CN,
				Quantity: int32(op.Quantity),
				Type: op.Product.Type.Value(),
				Version: "v1",
		}
		compute.Products = append(compute.Products, &product)
		computes = append(computes, &compute)
		fmt.Printf("%v\n", &compute)
	}

	ap, err := c.ApproveOrder(ctx, &od.ApprovalRequest{
		UserId: int64(order.CustomerID),
		ContractCode: order.ContractCode,
		Computes: computes,
	})

	if err != nil {
		fmt.Printf("Order Approval failed: %v\n", err)
		LOG.Error("Order Approval failed: ", err)
		return false, err
	}

	return ap.Approved, nil
}

func GetAllUsers(token string) ([]types.JSON, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, getMiqGrpc(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		LOG.Error("Error connecting: ", err)
		return nil, err
	}
	defer conn.Close()
	c := us.NewUserServiceClient(conn)

	resp, err := c.GetAllUsers(ctx, &us.TokenRequest{Token: token})

	if err != nil {
		fmt.Printf("Failed to get users: %v\n", err)
		LOG.Error("Failed to get users: ", err)
		return nil, err
	}
	var users []types.JSON

	for _, v := range resp.UserProfiles {
		user := types.JSON{
			"id": v.User.Id,
			"userid": v.User.Userid,
			"name": v.User.Name,
			"email": v.User.Email,
			"phone_number": v.User.PhoneNumber,
			"role": v.User.Role,
			"status": v.User.Status,
			"user_type": v.UserType,
			"account_type": v.AccountType,
			"id_number": v.IdNumber,
			"id_issue_date": v.IdIssueDate,
			"id_issue_location": v.IdIssueLocation,
			"tax_number": v.TaxNumber,
			"address": v.Address,
			"date_of_birth": v.DateOfBirth,
			"company": v.Company,
			"rep_name": v.RepName,
			"rep_phone": v.RepPhone,
			"rep_email": v.RepEmail,
			"ref_name": v.RefName,
			"ref_phone": v.RefPhone,
			"ref_email": v.RefEmail,
		}
		users = append(users, user)
	}

	return users, nil
}

func TestClient(address string) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("did not connect: %v", err)
		LOG.Error("did not connect: ", err)
	}
	defer conn.Close()
	c := od.NewOrderServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetOrder(ctx, &od.OrderRequest{OrderId: 1, OrderIdx: 1})
	if err != nil {
		fmt.Printf("could not get order: %v", err)
		LOG.Error("could not get order: ", err)
	}
	LOG.Info("Order: ", r)
}

func (s *server) ApproveOrder (ctx context.Context, request *od.ApprovalRequest) (*od.Approval, error) {
	fmt.Printf("%v\n", request)
	ap := od.Approval{Approved: true}
	return &ap, nil
}