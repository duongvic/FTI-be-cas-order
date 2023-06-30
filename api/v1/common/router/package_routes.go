package router

import (
	apiPackage "casorder/api/v1/common/apis/package"
	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to the gin Engine
func registerPackageRoutes(r *gin.RouterGroup) {
	packageAPI := apiPackage.PackageApi{}
	auth := r.Group("/package")
	{
		auth.GET("../packages", packageAPI.ListPackages)
		auth.GET("/:id", packageAPI.GetPackageByID)
		auth.POST("../packages", packageAPI.CreatePackage)
		auth.PATCH("/:id", packageAPI.UpdatePackage)
		auth.DELETE("/:id", packageAPI.DeletePackage)
	}
}
