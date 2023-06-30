package router

import (
	apiOrderDtl "casorder/api/v1/common/apis/order_dtl"
	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to the gin Engine
func registerOrderDtlRoutes(r *gin.RouterGroup) {
	orderDtlAPI := apiOrderDtl.OrderDtlApi{}
	auth := r.Group("/order_dtl")
	{
		auth.GET("../order_dtls", orderDtlAPI.ListOrderDtls)
		auth.GET("/:id", orderDtlAPI.GetOrderDtlByID)
		auth.POST("../order_dtls", orderDtlAPI.CreateOrderDtl)
		auth.PATCH("/:id", orderDtlAPI.UpdateOrderDtl)
		auth.DELETE("/:id", orderDtlAPI.DeleteOrderDtl)
	}
}