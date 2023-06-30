package services

import (
	"casorder/db/models"
	"casorder/utils"
	"casorder/utils/mgrpc"
	"casorder/utils/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	coTime "github.com/khanhct/go-lib-core/utils/time"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"strings"
	"sync"
	"time"
)

type OrderService struct {
	CommonApi
}

var mutex = &sync.Mutex{}

func (os *OrderService) CreateNewOrder(data types.JSON, DB *gorm.DB) (*models.Order, error) {
	// validate order
	if err := os.validateNewOrder(data); err != nil {
		os.Error(400, err, err.Error())
		return nil, err
	}

	var orderType types.OrderType
	orderType = types.OrderType(data["order_type"].(string))

	serviceType := types.ServiceType(data["service_type"].(string))
	remark := ""
	if data["remark"] != nil {
		remark = data["remark"].(string)
	}

	contract := data["contract"].(types.JSON)
	startAt, err := coTime.ConvertStrToTime(contract["start_at"].(string), coTime.YYYYMMDD)
	if err != nil {
		return nil, err
	}

	endAt, err := coTime.ConvertStrToTime(contract["end_at"].(string), coTime.YYYYMMDD)
	if err != nil {
		return nil, err
	}

	if startAt.Before(endAt) == false {
		return nil, errors.New("end date must be greater than start date")
	}

	var contractCode string
	var orderCode string
	if orderType == types.OrderTypeTrial {
		prefix := fmt.Sprintf("%vTR", data["region_name"].(string))
		contractCode = generateContractCode(prefix, DB)
		orderCode = generateOrderCode(contractCode, DB)
	} else {
		contractCode = contract["code"].(string)
		orderCode = generateOrderCode(contractCode, DB)
	}

	orderDtl, err := os.createOrderDtl(DB, data)
	if err != nil {
		return nil, err
	}

	var coSale datatypes.JSON
	if data["co_sale"] != nil {
		coSale, err = json.Marshal(data["co_sale"])
		if err != nil {
			return nil, err
		}
	}

	order := models.Order{
		Code:         orderCode,
		ContractCode: contractCode,
		OrderType:    orderType,
		ServiceType:  serviceType,
		Status:       types.StatusNew,
		OrderDtlID:   orderDtl.ID,
		StartAt:      startAt,
		EndAt:        endAt,
		Duration:     int64(endAt.Sub(startAt).Hours() / 24),
		CustomerID:   uint(data["customer_id"].(float64)),
		Price:        data["price"].(float64),
		ApprovalStep: 0,
		RegionID:     uint(data["region_id"].(float64)),
		Remark:       remark,
		StaffID:      uint(data["staff_id"].(int64)),
		SaleCare:     data["sale_care"].(string),
		CoSale:       coSale,
	}

	if err := DB.Create(&order).Error; err != nil {
		DB.Delete(&orderDtl)
		return nil, err
	}

	computes, err := os.createOrderProducts(DB, order, data["items"].([]interface{}))
	if err != nil {
		return nil, err
	}

	result, err := os.Find(&order, order.ID, DB)
	if err != nil {
		os.Logger.Error(err.Error())
		os.Error(400, err, err.Error())
		return nil, err
	}

	mailStatus, err := mgrpc.SendOrderMail(result, computes, data)
	if err != nil || !mailStatus {
		os.Logger.Error(err.Error())
		os.Error(500, err, "error sending email")
		//return nil, err
	}

	return result, nil
}

func (os *OrderService) RenewOrder(data types.JSON, DB *gorm.DB) (*models.Order, error) {
	refOrder, err := os.validateRenewOrder(data, DB)
	if err != nil {
		return nil, err
	}

	contract := data["contract"].(types.JSON)

	startAt, err := coTime.ConvertStrToTime(contract["start_at"].(string), coTime.YYYYMMDD)
	if err != nil {
		return nil, err
	}

	endAt, err := coTime.ConvertStrToTime(contract["end_at"].(string), coTime.YYYYMMDD)
	if err != nil {
		return nil, err
	}

	if endAt.Before(refOrder.EndAt) {
		return nil, errors.New("new end date must be after current end date")
	}

	remark := ""
	if data["remark"] != nil {
		remark = data["remark"].(string)
	}

	newOrderDtl := refOrder.OrderDtl
	newOrderDtl.ID = 0
	if err := DB.Create(&newOrderDtl).Error; err != nil {
		return nil, err
	}

	newOrder := *refOrder
	orderCode := generateOrderCode(newOrder.ContractCode, DB)

	newOrder.ID = 0
	newOrder.CreatedAt = time.Now()
	newOrder.UpdatedAt = time.Now()
	newOrder.Code = orderCode
	newOrder.OrderType = types.OrderTypeRenew
	newOrder.Status = types.StatusPayCompleted
	newOrder.StartAt = startAt
	newOrder.EndAt = endAt
	newOrder.Duration = int64(newOrder.EndAt.Sub(newOrder.StartAt).Hours() / 24)
	newOrder.OrderDtl = newOrderDtl
	newOrder.RefOrderID = refOrder.ID
	newOrder.RefOrderIdx = 0
	newOrder.OrderProducts = nil
	newOrder.Remark = remark

	if err := DB.Create(&newOrder).Error; err != nil {
		DB.Delete(&newOrderDtl)
		return nil, err
	}

	var computes [][]models.OrderProduct
	items := data["items"].([]interface{})
	for idx, item := range items {
		var opList []models.OrderProduct
		orderIdx := int(item.(float64))
		for _, op := range refOrder.OrderProducts {
			if orderIdx == op.Idx {
				opList = append(opList, *op)
			}
		}

		var compute []models.OrderProduct
		for _, newOP := range opList {
			if newOP.Idx == int(item.(float64)) {
				newOP.ID = 0
				newOP.Idx = idx + 1
				newOP.OrderID = newOrder.ID
				if err := DB.Create(&newOP).Error; err != nil {
					return nil, err
				}
				compute = append(compute, newOP)
			}
		}
		computes = append(computes, compute)

		for _, op := range refOrder.OrderProducts {
			if orderIdx == op.Idx {
				op.Disabled = true
				_, err = os.Update(op, op.ID, DB)
				if err != nil {
					os.Logger.Error(err.Error())
					os.Error(500, err, "error disabling order products")
					return nil, err
				}
			}
		}
	}

	result, err := os.Find(&newOrder, newOrder.ID, DB)
	if err != nil {
		os.Logger.Error(err.Error())
		os.Error(400, err, err.Error())
		return nil, err
	}

	mailStatus, err := mgrpc.SendOrderMail(result, computes, data)
	if err != nil || !mailStatus {
		os.Logger.Error(err.Error())
		os.Error(500, err, "error sending email")
		//return nil, err
	}

	return nil, nil
}

func (os *OrderService) ExtendOrder(data types.JSON, DB *gorm.DB) (*models.Order, error) {
	refOrder, err := os.validateExtendOrder(data, DB)
	if err != nil {
		return nil, err
	}

	items := data["items"].([]interface{})
	for _, item := range items {
		contractCode := refOrder.ContractCode
		orderCode := generateOrderCode(contractCode, DB)
		newOrderDtl := refOrder.OrderDtl
		newOrderDtl.ID = 0

		if err := DB.Create(&newOrderDtl).Error; err != nil {
			return nil, err
		}
		itemDict := item.(types.JSON)

		newOrder := *refOrder
		newOrder.ID = 0
		newOrder.CreatedAt = time.Now()
		newOrder.UpdatedAt = time.Now()
		newOrder.OrderDtl = newOrderDtl
		newOrder.Code = orderCode
		newOrder.OrderType = types.OrderTypeExtend
		newOrder.Status = types.StatusPayCompleted
		newOrder.ApprovalStep = 0
		newOrder.ContractCode = contractCode
		newOrder.RefOrderID = refOrder.ID
		newOrder.RefOrderIdx = int64(itemDict["order_idx"].(float64))
		newOrder.OrderProducts = nil
		//newOrder.Remark = data["remark"].(string)
		if data["remark"] != nil {
			newOrder.Remark = data["remark"].(string)
		} else {
			data["remark"] = nil
		}
		newOrder.Price = data["price"].(float64)

		if err := DB.Create(&newOrder).Error; err != nil {
			DB.Delete(&newOrderDtl)
			return nil, err
		}
		computes, err := os.extendOrderProducts(DB, *refOrder, newOrder, itemDict)
		result, err := os.Find(&newOrder, newOrder.ID, DB)
		if err != nil {
			os.Logger.Error(err.Error())
			os.Error(400, err, err.Error())
			return nil, err
		}

		for _, op := range refOrder.OrderProducts {
			if op.Idx == int(itemDict["order_idx"].(float64)) {
				op.Disabled = true
				_, err = os.Update(op, op.ID, DB)
				if err != nil {
					os.Logger.Error(err.Error())
					os.Error(500, err, "error disabling order products")
					return nil, err
				}
			}
		}

		mailStatus, err := mgrpc.SendOrderMail(result, computes, data)
		if err != nil || !mailStatus {
			os.Logger.Error(err.Error())
			os.Error(500, err, "error sending email")
			//return nil, err
		}
	}

	return nil, nil
}

func (os *OrderService) createOrderDtl(DB *gorm.DB, data types.JSON) (*models.OrderDtl, error) {
	customer := data["customer"].(types.JSON)
	orderDtl := models.OrderDtl{}
	if err := utils.ParseObjectFromJson(customer, &orderDtl); err != nil {
		return nil, err
	}

	if err := DB.Create(&orderDtl).Error; err != nil {
		return nil, err
	}

	return &orderDtl, nil
}

func (os *OrderService) createOrderProducts(DB *gorm.DB, order models.Order, items []interface{}) ([][]models.OrderProduct, error) {
	var computes [][]models.OrderProduct
	for idx, item := range items {
		itemDict := item.(types.JSON)

		var products []models.OrderProduct
		if err := utils.ParseObjectFromJson(itemDict["products"].([]interface{}), &products); err != nil {
			return nil, err
		}

		productGroup := utils.Group(products)
		products = make([]models.OrderProduct, 0)
		for _, pds := range productGroup {
			p := utils.Reduce(nil, pds, nil)
			products = append(products, p)
		}

		var compute []models.OrderProduct
		for _, product := range products {
			product.ID = 0
			product.OrderID = order.ID
			product.Idx = idx + 1
			if err := DB.Create(&product).Error; err != nil {
				return nil, err
			}
			compute = append(compute, product)
		}
		computes = append(computes, compute)
	}
	return computes, nil
}

func (os *OrderService) extendOrderProducts(DB *gorm.DB, refOrder models.Order, newOrder models.Order, item types.JSON) ([][]models.OrderProduct, error) {
	var computes [][]models.OrderProduct
	var products []models.OrderProduct

	if err := utils.ParseObjectFromJson(item["products"].([]interface{}), &products); err != nil {
		return nil, err
	}
	var opList []models.OrderProduct

	orderIdx := int(item["order_idx"].(float64))
	for _, op := range refOrder.OrderProducts {
		if orderIdx == op.Idx {
			opList = append(opList, *op)
		}
	}

	opList = append(opList, products...)

	productGroup := utils.Group(opList)
	products = make([]models.OrderProduct, 0)
	for _, pds := range productGroup {
		p := utils.Reduce(nil, pds, nil)
		products = append(products, p)
	}

	var compute []models.OrderProduct
	for _, product := range products {
		product.ID = 0
		product.OrderID = newOrder.ID
		product.Idx = 1
		if err := DB.Create(&product).Error; err != nil {
			return nil, err
		}
		compute = append(compute, product)
	}

	computes = append(computes, compute)

	return computes, nil
}

func (os *OrderService) RawFind(model interface{}, value interface{}, DB *gorm.DB) (interface{}, error) {
	query := types.JSON{
		"id": value,
	}
	if err := DB.Model(model).Preload(clause.Associations).First(model, query).Error; err != nil {
		return model, err
	}
	return model, nil
}

func (os *OrderService) Find(model *models.Order, value interface{}, DB *gorm.DB) (*models.Order, error) {
	query := types.JSON{
		"id": value,
	}
	if err := DB.Model(model).Preload(clause.Associations).Preload(
		"OrderProducts.Unit").Preload(
		"OrderProducts.Product.Unit").First(model, query).Error; err != nil {
		return model, err
	}
	return model, nil
}

func (os *OrderService) FindBy(model *models.Order, field string, value interface{}, DB *gorm.DB) (*models.Order, error) {
	query := types.JSON{
		field: value,
	}
	if err := DB.Model(model).
		Preload(clause.Associations).
		Preload("OrderProducts.Unit").
		Preload("OrderProducts.Product.Unit").
		First(model, query).Error; err != nil {
		return model, err
	}
	return model, nil
}

func (os *OrderService) Lookup(list []*models.Order, pagination *models.Pagination, filter types.JSON, DB *gorm.DB) (*models.Pagination, error) {
	var queryString string
	max := len(filter)
	i := 1
	for k, v := range filter {
		switch v.(type) {
		case string:
			if strings.Contains(k, "order_dtls") {
				queryBits := strings.Split(k, ".")
				queryString += fmt.Sprintf("order_dtl_id in (select id from %v where %v LIKE \"%%%v%%\")", queryBits[0], queryBits[1], v)
			} else if strings.Contains(k, "customer_id") || strings.Contains(k, "contract_code") || strings.Contains(k, "order_id") {
				queryString += fmt.Sprintf("%v = \"%v\"", k, v)
			} else {
				queryString += fmt.Sprintf("%v LIKE \"%%%v%%\"", k, v)
			}
		default:
			queryString += fmt.Sprintf("%v = %v", k, v)
		}
		if i < max {
			queryString += " AND "
		}
		i++
	}
	if err := DB.Scopes(models.Paginate(list, queryString, pagination, DB)).
		Where(queryString).
		Preload(clause.Associations).
		Find(&list).Error; err != nil {
		return nil, err
	}
	pagination.Rows = list
	return pagination, nil
}

func (os *OrderService) UpdateOrderItems(newOrderProducts []models.OrderProduct, id uint, idx int, DB *gorm.DB) (*models.Order, error) {
	if err := os.DeleteOrderItems(id, idx, DB); err != nil {
		os.Logger.Error(err.Error())
		os.Error(400, err, err.Error())
		return nil, err
	}

	for _, newOP := range newOrderProducts {
		newOP.OrderID = id
		newOP.Idx = idx
		if err := DB.Create(&newOP).Error; err != nil {
			os.Logger.Error(err.Error())
			os.Error(400, err, err.Error())
			return nil, err
		}
	}

	var order models.Order
	result, err := os.Find(&order, id, DB)
	if err != nil {
		os.Logger.Error(err.Error())
		os.Error(400, err, err.Error())
		return nil, err
	}

	return result, nil
}

func (os *OrderService) DeleteOrderItems(id uint, idx int, DB *gorm.DB) error {
	var targetOrderProducts []*models.OrderProduct
	queryString := fmt.Sprintf("order_id = %v AND idx = %v", id, idx)
	if err := DB.Model(targetOrderProducts).
		Where(queryString).
		Preload(clause.Associations).
		Preload("Product.Unit").
		Find(&targetOrderProducts).Error; err != nil {
		os.Logger.Error(err.Error())
		os.Error(400, err, err.Error())
		return err
	}
	for _, op := range targetOrderProducts {
		if err := DB.Model(op).Delete(op).Error; err != nil {
			os.Logger.Error(err.Error())
			os.Error(400, err, err.Error())
			return err
		}
	}

	return nil
}

func (os *OrderService) Delete(id string, DB *gorm.DB) error {
	var order models.Order

	result, err := os.Find(&order, id, DB)
	if err != nil {
		os.Logger.Error(err.Error())
		os.Error(404, err, err.Error())
		return err
	}

	for _, op := range result.OrderProducts {
		if err := DB.Model(op).Where("id = ?", op.ID).Delete(&op).Error; err != nil {
			os.Logger.Error(err.Error())
			os.Error(404, err, err.Error())
			return err
		}
	}

	if err := DB.Model(order.OrderDtl).Where("id = ?", order.OrderDtl.ID).Delete(&order.OrderDtl).Error; err != nil {
		os.Logger.Error(err.Error())
		os.Error(404, err, err.Error())
		return err
	}

	if err := DB.Model(order).Where("id = ?", order.ID).Delete(&order).Error; err != nil {
		os.Logger.Error(err.Error())
		os.Error(404, err, err.Error())
		return err
	}

	return nil
}

func (os *OrderService) validateNewOrder(data types.JSON) error {
	requiredParams := []string{"order_type", "service_type", "customer", "items"}

	if err := utils.ValidateArgs(data, requiredParams); err != nil {
		return err
	}

	return nil
}

func (os *OrderService) validateRenewOrder(data types.JSON, DB *gorm.DB) (*models.Order, error) {
	requiredParams := []string{"order_id", "contract", "duration"}

	if err := utils.ValidateArgs(data, requiredParams); err != nil {
		return nil, err
	}

	var refOrder *models.Order

	if err := DB.Model(&refOrder).
		Where("id = ?", data["order_id"].(float64)).Preload(clause.Associations).
		First(&refOrder).Error; err != nil {
		return nil, err
	}

	return refOrder, nil
}

func (os *OrderService) validateExtendOrder(data types.JSON, DB *gorm.DB) (*models.Order, error) {
	requiredParams := []string{"order_id", "contract", "order_type", "items"}
	if err := utils.ValidateArgs(data, requiredParams); err != nil {
		return nil, err
	}

	var refOrder *models.Order
	if err := DB.Model(&refOrder).
		Where("id = ?", data["order_id"].(float64)).Preload(clause.Associations).
		First(&refOrder).Error; err != nil {
		return nil, err
	}

	if refOrder.EndAt.Before(time.Now()) {
		err := errors.New("order expired")
		return nil, err
	}

	return refOrder, nil
}

func generateContractCode(prefix string, DB *gorm.DB) string {
	var lastOrder *models.Order
	var idx = 0

	mutex.Lock()
	err := DB.Model(&lastOrder).Unscoped().Order("id DESC").First(&lastOrder).Error
	if err != nil {
		// idx = 1 if not found records
		idx = 1
	} else {
		idx = int(lastOrder.ID + 1)
	}

	mutex.Unlock()

	return fmt.Sprintf("%v%05d", prefix, idx)
}

func generateOrderCode(prefix string, DB *gorm.DB) string {
	var lastOrder *models.Order
	var idx = 0

	// Make sure that the latest order will be got
	mutex.Lock()
	err := DB.Model(&lastOrder).Unscoped().Order("id DESC").First(&lastOrder).Error
	if err != nil {
		// idx = 1 if not found records
		idx = 1
	} else {
		idx = int(lastOrder.ID + 1)
	}

	mutex.Unlock()
	return fmt.Sprintf("%v-%05d", prefix, idx)
}

func (os *OrderService) ParseObjectFromRequest(c *gin.Context, order *models.Order) (models.Order, error) {
	reader := io.Reader(c.Request.Body)
	decoder := json.NewDecoder(reader)

	err := decoder.Decode(&order)
	if err != nil {
		os.Logger.Error(err.Error())
	}
	return *order, err
}
