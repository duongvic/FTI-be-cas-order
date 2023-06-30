package apis

import (
	services "casorder/api/v1/common/services"
	"casorder/db/models"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type OrderProductApi struct {
	services.OrderProductService
}

// ListOrderProducts Gets the list of all order products
func (op OrderProductApi) ListOrderProducts(c *gin.Context) {
	if err := op.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		op.Error(500, err, err.Error())
		op.Logger.Error(err.Error())
		return
	}
	var orderProduct models.OrderProduct
	var orderProducts []*models.OrderProduct
	var pagination models.Pagination
	pagination.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	pagination.Page, _ = strconv.Atoi(c.Query("page"))
	pagination.Sort = c.Query("sort")

	result, err := op.List(orderProduct, &orderProducts, &pagination, op.Orm)
	if err != nil {
		op.Logger.Error(err.Error())
		op.Error(404, err, err.Error())
		return
	}

	op.OK(result, "Success")
}

// GetOrderProductByID Gets specific order product from ID in context
func (op OrderProductApi) GetOrderProductByID(c *gin.Context) {
	if err := op.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		op.Error(500, err, err.Error())
		op.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var orderProduct models.OrderProduct

	result, err := op.Find(&orderProduct, id, op.Orm)
	if err != nil {
		op.Logger.Error(err.Error())
		op.Error(404, err, err.Error())
		return
	}

	op.OK(result, "Success")
}

