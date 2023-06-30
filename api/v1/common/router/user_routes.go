package router

import (
	apiUser "casorder/api/v1/common/apis/user"
	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to the gin Engine
func registerUserRoutes(r *gin.RouterGroup) {
	userAPI := apiUser.UserApi{}
	auth := r.Group("/user")
	{
		auth.GET("../users", userAPI.ListMiqUsers)
		auth.GET("/:id", userAPI.GetUserByID)
		auth.GET("", userAPI.ListUsers)
		auth.POST("../users", userAPI.CreateUser)
		auth.PATCH("/:id", userAPI.UpdateUser)
		auth.DELETE("/:id", userAPI.DeleteUser)
	}
}
