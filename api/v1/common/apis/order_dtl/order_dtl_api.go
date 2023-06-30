package apis

import (
	services "casorder/api/v1/common/services"
	"casorder/db/models"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type OrderDtlApi struct {
	services.CommonApi
}

// ListOrderDtls Gets the list of all order details
func (od OrderDtlApi) ListOrderDtls(c *gin.Context) {
	if err := od.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		od.Error(500, err, err.Error())
		od.Logger.Error(err.Error())
		return
	}
	var orderDtl models.OrderDtl
	var orders []*models.OrderDtl
	var pagination models.Pagination
	pagination.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	pagination.Page, _ = strconv.Atoi(c.Query("page"))
	pagination.Sort = c.Query("sort")

	result, err := od.List(orderDtl, &orders, &pagination, od.Orm)

	if err != nil {
		od.Logger.Error(err.Error())
		od.Error(404, err, err.Error())
		return
	}

	od.OK(result, "Success")
}

// GetOrderDtlByID Gets specific order detail from ID in context
func (od OrderDtlApi) GetOrderDtlByID(c *gin.Context) {
	if err := od.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		od.Error(500, err, err.Error())
		od.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var orderDtl models.OrderDtl

	result, err := od.Find(&orderDtl, id, od.Orm)
	if err != nil {
		od.Logger.Error(err.Error())
		od.Error(404, err, err.Error())
		return
	}

	od.OK(result, "Success")
}

// CreateOrderDtl Creates a new order detail
func (od OrderDtlApi) CreateOrderDtl(c *gin.Context) {
	if err := od.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		od.Error(500, err, err.Error())
		od.Logger.Error(err.Error())
		return
	}
	var model models.OrderDtl

	orderDtl, err := od.ParseObjectFromRequest(c, &model)
	if err != nil {
		od.Logger.Error(err.Error())
		od.Error(400, err, err.Error())
	}

	if err := od.Orm.Create(orderDtl).Error; err != nil {
		od.Logger.Error(err.Error())
		od.Error(400, err, err.Error())
		return
	}

	od.OK(orderDtl, "Success")
}

// UpdateOrderDtl Updates an order detail
func (od OrderDtlApi) UpdateOrderDtl(c *gin.Context) {
	if err := od.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		od.Error(500, err, err.Error())
		od.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var model *models.OrderDtl

	orderDtl, err := od.ParseObjectFromRequest(c, &model)
	if err != nil {
		od.Logger.Error(err.Error())
		od.Error(400, err, err.Error())
		return
	}

	if err := od.Orm.Model(&model).Where("id = ?", id).Updates(orderDtl).Error; err != nil {
		od.Logger.Error(err.Error())
		od.Error(400, err, err.Error())
		return
	}
	result, err := od.Find(&model, id, od.Orm)
	od.OK(result, "Success")
}

// DeleteOrderDtl Deletes an order detail
func (od OrderDtlApi) DeleteOrderDtl(c *gin.Context) {
	if err := od.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		od.Error(500, err, err.Error())
		od.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var orderDtl models.OrderDtl

	if err := od.Orm.First(&orderDtl, "id = ?", id).Delete(&orderDtl).Error; err != nil {
		od.Logger.Error(err.Error())
		od.Error(404, err, err.Error())
		return
	}
	od.OK(orderDtl, "Success")
}