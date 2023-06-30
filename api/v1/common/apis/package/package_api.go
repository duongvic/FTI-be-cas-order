package apis

import (
	services "casorder/api/v1/common/services"
	"casorder/db/models"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type PackageApi struct {
	services.CommonApi
}

// ListPackages Gets the list of all packages
func (p PackageApi) ListPackages(c *gin.Context) {
	if err := p.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		p.Error(500, err, err.Error())
		p.Logger.Error(err.Error())
		return
	}
	var pkg models.Package
	var pkgs []*models.Package
	var pagination models.Pagination
	pagination.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	pagination.Page, _ = strconv.Atoi(c.Query("page"))
	pagination.Sort = c.Query("sort")

	result, err := p.List(pkg, &pkgs, &pagination, p.Orm)
	if err != nil {
		p.Logger.Errorf(err.Error())
		p.Error(404, err, err.Error())
		return
	}

	p.OK(result, "Success")
}

// GetPackageByID Gets specific package from ID in context
func (p PackageApi) GetPackageByID(c *gin.Context) {
	if err := p.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		p.Error(500, err, err.Error())
		p.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var pkg models.Package

	result, err := p.Find(&pkg, id, p.Orm)
	if err != nil {
		p.Logger.Errorf(err.Error())
		p.Error(404, err, err.Error())
		return
	}

	p.OK(result, "Success")
}

// CreatePackage Creates a new package
func (p PackageApi) CreatePackage(c *gin.Context) {
	if err := p.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		p.Error(500, err, err.Error())
		p.Logger.Error(err.Error())
		return
	}
	var model models.Package

	pkg, err := p.ParseObjectFromRequest(c, &model)
	if err != nil {
		p.Logger.Errorf(err.Error())
		p.Error(400, err, err.Error())
	}

	if err := p.Orm.Create(pkg).Error; err != nil {
		p.Logger.Errorf(err.Error())
		p.Error(400, err, err.Error())
		return
	}

	p.OK(pkg, "Success")
}

// UpdatePackage Updates a package
func (p PackageApi) UpdatePackage(c *gin.Context) {
	if err := p.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		p.Error(500, err, err.Error())
		p.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var model *models.Package

	pkg, err := p.ParseObjectFromRequest(c, &model)
	if err != nil {
		p.Logger.Errorf(err.Error())
		p.Error(400, err, err.Error())
		return
	}

	if err := p.Orm.Model(&model).Where("id = ?", id).Updates(pkg).Error; err != nil {
		p.Logger.Errorf(err.Error())
		p.Error(400, err, err.Error())
		return
	}
	result, err := p.Find(&model, id, p.Orm)
	p.OK(result, "Success")
}

// DeletePackage Deletes a package
func (p PackageApi) DeletePackage(c *gin.Context) {
	if err := p.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		p.Error(500, err, err.Error())
		p.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var pkg models.Package

	if err := p.Orm.First(&pkg, "id = ?", id).Delete(&pkg).Error; err != nil {
		p.Logger.Errorf(err.Error())
		p.Error(404, err, err.Error())
		return
	}
	p.OK(pkg, "Success")
}