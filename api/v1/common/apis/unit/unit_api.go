package apis

import (
	services "casorder/api/v1/common/services"
	"casorder/db/models"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type UnitApi struct {
	services.CommonApi
}

// GetUnits Gets the list of all units
func (u UnitApi) GetUnits(c *gin.Context) {
	err := u.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors
	if err != nil {
		if err != nil {
			u.Error(500, err, err.Error())
			u.Logger.Error(err.Error())
			return
		}
	}
	//pageSize := c.Param("PageSize")
	var unit models.Unit
	var units []*models.Unit
	var pagination models.Pagination
	pagination.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	pagination.Page, _ = strconv.Atoi(c.Query("page"))
	pagination.Sort = c.Query("sort")

	result, err := u.List(unit, &units, &pagination, u.Orm)

	if err != nil {
		u.Logger.Errorf("Error getting units: %v", err)
		u.Error(404, err, "Failed to get units")
		return
	}

	u.OK(result, "Success")
}

// GetUnitByID Gets specific unit from ID in context
func (u UnitApi) GetUnitByID(c *gin.Context) {
	err := u.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors
	if err != nil {
		if err != nil {
			u.Error(500, err, err.Error())
			u.Logger.Error(err.Error())
			return
		}
	}
	id := c.Param("id")
	var unit models.Unit

	result, err := u.Find(&unit, id, u.Orm)
	if err != nil {
		u.Logger.Errorf("Error getting unit: %v", err)
		u.Error(404, err, "Failed to get unit")
		return
	}

	u.OK(result, "Success")
}

// CreateUnit Creates a new unit
func (u UnitApi) CreateUnit(c *gin.Context) {
	err := u.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors
	if err != nil {
		if err != nil {
			u.Error(500, err, err.Error())
			u.Logger.Error(err.Error())
			return
		}
	}

	var model models.Unit

	unit, err := u.ParseObjectFromRequest(c, &model)
	if err != nil {
		u.Logger.Errorf("Error creating unit: %v", err)
		u.Error(400, err, "Failed to create unit")
	}

	if err := u.Orm.Create(unit).Error; err != nil {
		u.Logger.Errorf("Error creating unit: %v", err)
		u.Error(400, err, "Failed to create unit")
		return
	}

	u.OK(unit, "Success")
}

// UpdateUnit Updates a unit
func (u UnitApi) UpdateUnit(c *gin.Context) {
	err := u.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors
	if err != nil {
		if err != nil {
			u.Error(500, err, err.Error())
			u.Logger.Error(err.Error())
			return
		}
	}
	id := c.Param("id")
	var model *models.Unit

	unit, err := u.ParseObjectFromRequest(c, &model)
	if err != nil {
		u.Logger.Errorf("Error reading data: %v", err)
		u.Error(400, err, "Failed to update unit")
		return
	}

	if err := u.Orm.Model(&model).Where("id = ?", id).Updates(unit).Error; err != nil {
		u.Logger.Errorf("Error creating unit: %v", err)
		u.Error(400, err, "Failed to update unit")
		return
	}
	result, err := u.Find(&model, id, u.Orm)
	u.OK(result, "Success")
}

// DeleteUnit Deletes a unit
func (u UnitApi) DeleteUnit(c *gin.Context) {
	err := u.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors
	if err != nil {
		if err != nil {
			u.Error(500, err, err.Error())
			u.Logger.Error(err.Error())
			return
		}
	}
	id := c.Param("id")
	var unit models.Unit

	if err := u.Orm.First(&unit, "id = ?", id).Delete(&unit).Error; err != nil {
		u.Logger.Errorf("Error deleting unit: %v", err)
		u.Error(404, err, "Failed to delete unit")
		return
	}
	u.OK(unit, "Success")
}