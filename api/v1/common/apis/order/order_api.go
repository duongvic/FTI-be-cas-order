package apis

import (
	"casorder/api/v1/common/services"
	"casorder/db/models"
	"casorder/utils"
	"casorder/utils/mgrpc"
	"casorder/utils/types"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type OrderApi struct {
	services.OrderService
}


// ListOrders Gets the list of orders by query and
// returns all if there is no query
func (o OrderApi) ListOrders(c *gin.Context) {
	if err := o.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		o.Error(500, err, err.Error())
		o.Logger.Error(err.Error())
		return
	}
	var orders []*models.Order
	var pagination models.Pagination

	filter := types.JSON{}

	for k, v := range o.Context.Request.URL.Query() {
		for _, i := range v {
			if k != "pageSize" && k != "page" && k != "sort" {
				filter[k] = i
			}
		}
	}

	pagination.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	pagination.Page, _ = strconv.Atoi(c.Query("page"))
	pagination.Sort = c.Query("sort")

	_, err := o.Lookup(orders, &pagination, filter, o.Orm)
	if err != nil {
		o.Logger.Errorf("Error getting orders: %v", err)
		o.Error(404, err, err.Error())
		return
	}

	o.OK(pagination, "Success")
}


// GetOrderByID Gets specific order from ID in context
func (o OrderApi) GetOrderByID(c *gin.Context) {
	if err := o.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		o.Error(500, err, err.Error())
		o.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var order models.Order

	result, err := o.Find(&order, id, o.Orm)
	if err != nil {
		o.Logger.Errorf(err.Error())
		o.Error(404, err, err.Error())
		return
	}

	o.OK(result, "Success")
}

// CreateOrder Creates a new order
func (o OrderApi) CreateOrder(c *gin.Context) {
	if err := o.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		o.Error(500, err, err.Error())
		o.Logger.Error(err.Error())
		return
	}

	data, err := o.ParseJsonFromRequest(c)
	if err != nil {
		o.Error(400, err, err.Error())
		return
	}

	user, _ := c.Get("user")
	data["staff_id"] = user.(*models.User).MiqID
	data["staff_username"] = user.(*models.User).Username

	var orderType types.OrderType
	var result *models.Order

	orderType = types.OrderType(data["order_type"].(string))
	if orderType == types.OrderTypeBuy || orderType	== types.OrderTypeTrial {
		result, err = o.CreateNewOrder(data, o.Orm)
		if err != nil {
			o.Error(400, err, err.Error())
			return
		}
	} else if orderType == types.OrderTypeRenew {
		result, err = o.RenewOrder(data, o.Orm)
		if err != nil {
			o.Error(400, err, err.Error())
			return
		}
	} else if orderType == types.OrderTypeExtend {
		result, err = o.ExtendOrder(data, o.Orm)
		if err != nil {
			o.Error(400, err, err.Error())
			return
		}
	}

	o.OK(result, "Success")
}

// UpdateOrder Updates an order
func (o OrderApi) UpdateOrder(c *gin.Context) {
	if err := o.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		o.Error(500, err, err.Error())
		o.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var model models.Order

	order, err := o.ParseObjectFromRequest(c, &model)
	if err != nil {
		o.Logger.Error(err.Error())
		o.Error(400, err, err.Error())
		return
	}

	_, err = o.Update(order, id, o.Orm)
	if err != nil {
		o.Logger.Error(err.Error())
		o.Error(404, err, err.Error())
		return
	}
	result, err := o.Find(&model, id, o.Orm)
	o.OK(result, "Success")
}

// UpdateOrderProducts Updates order products
func (o OrderApi) UpdateOrderProducts(c *gin.Context) {
	if err := o.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		o.Error(500, err, err.Error())
		o.Logger.Error(err.Error())
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	idx, _ := strconv.Atoi(c.Param("idx"))
	var orderProducts []models.OrderProduct

	data, err := o.ParseJsonFromRequest(c)
	if err != nil {
		o.Error(400, err, err.Error())
		return
	}

	err = utils.ParseObjectFromJson(data["items"], &orderProducts)
	if err != nil {
		o.Logger.Error(err.Error())
		o.Error(400, err, err.Error())
		return
	}

	result, err := o.UpdateOrderItems(orderProducts, uint(id), idx, o.Orm)
	if err != nil {
		o.Logger.Error(err.Error())
		o.Error(400, err, err.Error())
		return
	}

	o.OK(result, "Success")
}

// DeleteOrderProducts Deletes order products
func (o OrderApi) DeleteOrderProducts(c *gin.Context) {
	if err := o.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		o.Error(500, err, err.Error())
		o.Logger.Error(err.Error())
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	idx, _ := strconv.Atoi(c.Param("idx"))

	if err := o.DeleteOrderItems(uint(id), idx, o.Orm); err != nil {
		o.Logger.Error(err.Error())
		o.Error(400, err, err.Error())
		return
	}

	o.OK(nil, "Success")
}

// DeleteOrder Deletes an order
func (o OrderApi) DeleteOrder(c *gin.Context) {
	if err := o.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		o.Error(500, err, err.Error())
		o.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")

	if err := o.Delete(id, o.Orm); err != nil {
		o.Logger.Error(err.Error())
		o.Error(400, err, err.Error())
		return
	}

	o.OK(nil, "Success")
}

func (o *OrderApi) ApproveOrder(c *gin.Context) {
	if err := o.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		o.Error(500, err, err.Error())
		o.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var model models.Order

	order, err := o.Find(&model, id, o.Orm)
	if err != nil {
		o.Logger.Error(err.Error())
		o.Error(404, err, err.Error())
		return
	}

	approved, err := mgrpc.GetApproval(order, o.Orm)
	if err != nil {
		o.Logger.Error(err.Error())
		o.Error(404, err, err.Error())
		return
	}

	if approved {
		order.ApprovalStep += 1
		order.Status = types.StatusApproved
		_, err = o.Update(order, id, o.Orm)
		if err != nil {
			o.Logger.Errorf(err.Error())
			o.Error(404, err, err.Error())
			return
		}
	}

	o.OK(order, "Success")
}