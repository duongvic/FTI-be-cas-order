package router

import (
	apiRegion "casorder/api/v1/common/apis/region"
	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to the gin Engine
func registerRegionRoutes(r *gin.RouterGroup) {
	regionAPI := apiRegion.RegionApi{}
	auth := r.Group("/region")
	{
		auth.GET("../regions", regionAPI.ListRegions)
		auth.GET("/:id", regionAPI.GetRegionByID)
		auth.POST("../regions", regionAPI.CreateRegion)
		auth.PATCH("/:id", regionAPI.UpdateRegion)
		auth.DELETE("/:id", regionAPI.DeleteRegion)
	}
}
