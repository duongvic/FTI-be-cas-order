package router

import (
	apiOrder "casorder/api/v1/common/apis/order"
	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to the gin Engine
func registerOrderRoutes(r *gin.RouterGroup) {
	orderAPI := apiOrder.OrderApi{}
	auth := r.Group("")
	{
		auth.GET("/orders", orderAPI.ListOrders)
		auth.GET("/order/:id", orderAPI.GetOrderByID)
		auth.POST("/orders", orderAPI.CreateOrder)
		auth.PUT("/order/:id", orderAPI.ApproveOrder)
		auth.PUT("/order/:id/order_products/:idx", orderAPI.UpdateOrderProducts)
		auth.PATCH("/order/:id", orderAPI.UpdateOrder)
		auth.DELETE("/order/:id", orderAPI.DeleteOrder)
		auth.DELETE("/order/:id/order_products/:idx", orderAPI.DeleteOrderProducts)
	}
}
