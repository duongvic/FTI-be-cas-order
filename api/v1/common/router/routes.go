package router

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes ApplyRoutes applies router to the gin Engine
func RegisterRoutes(r *gin.RouterGroup) {

	v1 := r.Group("")
	{
		registerHealthRoutes(v1)
		registerOrderDtlRoutes(v1)
		registerOrderProductRoutes(v1)
		registerOrderRoutes(v1)
		registerPackageRoutes(v1)
		registerProductRoutes(v1)
		registerRegionRoutes(v1)
		registerTransactionRoutes(v1)
		registerUserRoutes(v1)
		registerUnitRoutes(v1)
	}
}
