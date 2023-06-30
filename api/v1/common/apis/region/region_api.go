package region

import (
	services "casorder/api/v1/common/services"
	"casorder/db/models"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type RegionApi struct {
	services.CommonApi
}

// ListRegions Gets the list of all regions
func (r RegionApi) ListRegions(c *gin.Context) {
	if err := r.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		r.Error(500, err, err.Error())
		r.Logger.Error(err.Error())
		return
	}
	
	var region models.Region
	var regions []*models.Region
	var pagination models.Pagination
	pagination.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	pagination.Page, _ = strconv.Atoi(c.Query("page"))
	pagination.Sort = c.Query("sort")

	result, err := r.List(region, &regions, &pagination, r.Orm)
	if err != nil {
		r.Logger.Errorf("Error getting regions: %v", err)
		r.Error(404, err, "Failed to get regions")
		return
	}

	r.OK(result, "Success")
}

// GetRegionByID Gets specific region from ID in context
func (r RegionApi) GetRegionByID(c *gin.Context) {
	if err := r.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		r.Error(500, err, err.Error())
		r.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var region models.Region

	result, err := r.Find(&region, id, r.Orm)
	if err != nil {
		r.Logger.Errorf("Error getting region: %v", err)
		r.Error(404, err, "Failed to get region")
		return
	}

	r.OK(result, "Success")
}

// CreateRegion Creates a new region
func (r RegionApi) CreateRegion(c *gin.Context) {
	if err := r.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		r.Error(500, err, err.Error())
		r.Logger.Error(err.Error())
		return
	}
	
	var model models.Region

	region, err := r.ParseObjectFromRequest(c, &model)
	if err != nil {
		r.Logger.Errorf("Error creating region: %v", err)
		r.Error(400, err, "Failed to create region")
	}

	if err := r.Orm.Create(region).Error; err != nil {
		r.Logger.Errorf("Error creating region: %v", err)
		r.Error(400, err, "Failed to create region")
		return
	}

	r.OK(region, "Success")
}

// UpdateRegion Updates a region
func (r RegionApi) UpdateRegion(c *gin.Context) {
	if err := r.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		r.Error(500, err, err.Error())
		r.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var model *models.Region

	region, err := r.ParseObjectFromRequest(c, &model)
	if err != nil {
		r.Logger.Errorf("Error reading data: %v", err)
		r.Error(400, err, "Failed to update region")
		return
	}

	if err := r.Orm.Model(&model).Where("id = ?", id).Updates(region).Error; err != nil {
		r.Logger.Errorf("Error creating region: %v", err)
		r.Error(400, err, "Failed to update region")
		return
	}
	result, err := r.Find(&model, id, r.Orm)
	r.OK(result, "Success")
}

// DeleteRegion Deletes a region
func (r RegionApi) DeleteRegion(c *gin.Context) {
	if err := r.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		r.Error(500, err, err.Error())
		r.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var region models.Region

	if err := r.Orm.First(&region, "id = ?", id).Delete(&region).Error; err != nil {
		r.Logger.Errorf("Error deleting region: %v", err)
		r.Error(404, err, "Failed to delete region")
		return
	}
	r.OK(region, "Success")
}