package router

import (
	apiUnit "casorder/api/v1/common/apis/unit"
	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to the gin Engine
func registerUnitRoutes(r *gin.RouterGroup) {
	unitAPI := apiUnit.UnitApi{}
	auth := r.Group("/unit")
	{
		auth.GET("../units", unitAPI.GetUnits)
		auth.GET("/:id", unitAPI.GetUnitByID)
		auth.POST("../units", unitAPI.CreateUnit)
		auth.PATCH("/:id", unitAPI.UpdateUnit)
		auth.DELETE("/:id", unitAPI.DeleteUnit)
	}
}
