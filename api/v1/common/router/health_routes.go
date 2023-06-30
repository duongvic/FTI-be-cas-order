package router

import (
	"github.com/gin-gonic/gin"

	apiHealth "casorder/api/v1/common/apis"
)

// ApplyRoutes applies router to the gin Engine
func registerHealthRoutes(r *gin.RouterGroup) {
	healthAPI := apiHealth.HealthApi{}
	auth := r.Group("/health")
	{
		auth.GET("/check", healthAPI.Check)
	}
}
