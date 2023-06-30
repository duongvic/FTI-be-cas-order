package apis

import (
	services "casorder/api/v1/common/services"
	"casorder/db/models"
	"casorder/utils/mgrpc"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)


type UserApi struct {
	services.CommonApi
}

//ListUsers Gets the list of all users
func (u UserApi) ListUsers(c *gin.Context) {
	if err := u.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		u.Error(500, err, err.Error())
		u.Logger.Error(err.Error())
		return
	}
	
	var user models.User
	var users []*models.User
	var pagination models.Pagination
	pagination.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	pagination.Page, _ = strconv.Atoi(c.Query("page"))
	pagination.Sort = c.Query("sort")

	result, err := u.List(user, &users,&pagination, u.Orm)
	if err != nil {
		u.Logger.Errorf("Error getting users: %v", err)
		u.Error(404, err, "Failed to get user")
		return
	}

	u.OK(result, "Success")
}

// GetUserByID Gets a specific user from ID in context
func (u UserApi) GetUserByID(c *gin.Context) {
	if err := u.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		u.Error(500, err, err.Error())
		u.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var user models.User

	result, err := u.Find(&user, id, u.Orm)
	if err != nil {
		u.Logger.Errorf("Error getting user: %v", err)
		u.Error(404, err, "Failed to get user")
		return
	}

	u.OK(result, "Success")
}

// CreateUser Creates a new user
func (u UserApi) CreateUser(c *gin.Context) {
	if err := u.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		u.Error(500, err, err.Error())
		u.Logger.Error(err.Error())
		return
	}
	
	var model models.User

	user, err := u.ParseObjectFromRequest(c, &model)
	if err != nil {
		u.Logger.Errorf("Error creating user: %v", err)
		u.Error(400, err, "Failed to create user")
	}

	if err := u.Orm.Create(user).Error; err != nil {
		u.Logger.Errorf("Error creating user: %v", err)
		u.Error(400, err, "Failed to create user")
		return
	}

	u.OK(user, "Success")
}

//UpdateUser Updates a user
func (u UserApi) UpdateUser(c *gin.Context) {
	if err := u.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		u.Error(500, err, err.Error())
		u.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var model *models.User

	user, err := u.ParseObjectFromRequest(c, &model)
	if err != nil {
		u.Logger.Errorf("Error reading data: %v", err)
		u.Error(400, err, "Failed to update user")
		return
	}

	if err := u.Orm.Model(&model).Where("id = ?", id).Updates(user).Error; err != nil {
		u.Logger.Errorf("Error creating order: %v", err)
		u.Error(400, err, "Failed to update user")
		return
	}
	result, err := u.Find(&model, id, u.Orm)
	u.OK(result, "Success")
}

//DeleteUser Deletes a user
func (u UserApi) DeleteUser(c *gin.Context) {
	if err := u.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		u.Error(500, err, err.Error())
		u.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var user models.User

	if err := u.Orm.First(&user, "id = ?", id).Delete(&user).Error; err != nil {
		u.Logger.Errorf("Error deleting user: %v", err)
		u.Error(404, err, "Failed to delete user")
		return
	}
	u.OK(user, "Success")
}

// ListMiqUsers Gets the list users from MIQ via GRPC using API token
func(u UserApi) ListMiqUsers(c *gin.Context) {
	if err := u.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		u.Error(500, err, err.Error())
		u.Logger.Error(err.Error())
		return
	}
	token := c.Request.Header.Get("Authorization")
	users, err := mgrpc.GetAllUsers(token)
	if err != nil {
		u.Logger.Errorf("Error getting users: %v", err)
		u.Error(401, err, "Error: Failed to get users")
		return
	}
	u.OK(users, "Success")
}