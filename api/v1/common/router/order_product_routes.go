package router

import (
	apiOrderProduct "casorder/api/v1/common/apis/order_product"
	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to the gin Engine
func registerOrderProductRoutes(r *gin.RouterGroup) {
	orderProductAPI := apiOrderProduct.OrderProductApi{}
	auth := r.Group("/")
	{
		auth.GET("/order_products", orderProductAPI.ListOrderProducts)
		auth.GET("/order_product/:id", orderProductAPI.GetOrderProductByID)
	}
}
