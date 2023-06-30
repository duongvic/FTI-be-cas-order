package apiv1

import (
	"casorder/api/v1/common/router"

	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v2/order")
	{
		router.RegisterRoutes(v1)
	}
}
