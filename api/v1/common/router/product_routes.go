package router

import (
	apiProduct "casorder/api/v1/common/apis/product"
	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to the gin Engine
func registerProductRoutes(r *gin.RouterGroup) {
	productAPI := apiProduct.ProductApi{}
	auth := r.Group("/product")
	{
		auth.GET("../products", productAPI.ListProducts)
		auth.GET("/:id", productAPI.GetProductByID)
		auth.POST("../products", productAPI.CreateProduct)
		auth.PATCH("/:id", productAPI.UpdateProduct)
		auth.DELETE("/:id", productAPI.DeleteProduct)
	}
}
