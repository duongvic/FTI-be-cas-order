package apis

import (
	services "casorder/api/v1/common/services"
	"casorder/db/models"
	"casorder/utils/types"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type ProductApi struct {
	services.ProductService
}

// ListProducts Gets the list of all products
func (p ProductApi) ListProducts(c *gin.Context) {
	if err := p.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		p.Error(500, err, err.Error())
		p.Logger.Error(err.Error())
		return
	}
	var products []*models.Product
	var pagination models.Pagination
	filter := types.JSON{}

	for k, v := range p.Context.Request.URL.Query() {
		for _, i := range v {
			if k != "pageSize" && k != "page" && k != "sort" {
				filter[k] = i
			}
		}
	}

	pagination.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	pagination.Page, _ = strconv.Atoi(c.Query("page"))
	pagination.Sort = c.Query("sort")

	result, err := p.Lookup(&products,&pagination,filter, p.Orm)
	if err != nil {
		p.Logger.Errorf("Error getting products: %v", err)
		p.Error(404, err, "Failed to get products")
		return
	}

	p.OK(result, "Success")
}

// GetProductByID Gets specific product from ID in context
func (p ProductApi) GetProductByID(c *gin.Context) {
	if err := p.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		p.Error(500, err, err.Error())
		p.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var product models.Product

	result, err := p.Find(&product, id, p.Orm)
	if err != nil {
		p.Logger.Errorf("Error getting product: %v", err)
		p.Error(404, err, "Failed to get product")
		return
	}

	p.OK(result, "Success")
}

// CreateProduct Creates a new product
func (p ProductApi) CreateProduct(c *gin.Context) {
	if err := p.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		p.Error(500, err, err.Error())
		p.Logger.Error(err.Error())
		return
	}
	var model models.Product

	product, err := p.ParseObjectFromRequest(c, &model)
	if err != nil {
		p.Logger.Errorf("Error creating product: %v", err)
		p.Error(400, err, "Failed to create product")
	}

	if err := p.Orm.Create(product).Error; err != nil {
		p.Logger.Errorf("Error creating product: %v", err)
		p.Error(400, err, "Failed to create product")
		return
	}

	p.OK(product, "Success")
}

// UpdateProduct Updates a product
func (p ProductApi) UpdateProduct(c *gin.Context) {
	if err := p.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		p.Error(500, err, err.Error())
		p.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var model *models.Product

	product, err := p.ParseObjectFromRequest(c, &model)
	if err != nil {
		p.Logger.Errorf("Error reading data: %v", err)
		p.Error(400, err, "Failed to update product")
		return
	}

	if err := p.Orm.Model(&model).Where("id = ?", id).Updates(product).Error; err != nil {
		p.Logger.Errorf("Error creating product: %v", err)
		p.Error(400, err, "Failed to update product")
		return
	}

	result, err := p.Find(&model, id, p.Orm)
	p.OK(result, "Success")
}

// DeleteProduct Deletes a product
func (p ProductApi) DeleteProduct(c *gin.Context) {
	if err := p.MakeContext(c).MakeOrm(nil).MakeLogger(nil).Errors; err != nil {
		p.Error(500, err, err.Error())
		p.Logger.Error(err.Error())
		return
	}
	id := c.Param("id")
	var product models.Product

	if err := p.Orm.First(&product, "id = ?", id).Delete(&product).Error; err != nil {
		p.Logger.Errorf("Error deleting product: %v", err)
		p.Error(404, err, "Failed to delete product")
		return
	}
	p.OK(product, "Success")
}