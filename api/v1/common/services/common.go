package services

import (
	"casorder/db/models"
	"casorder/utils/types"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/khanhct/go-lib-core/sdk/api"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"io/ioutil"
)

type CommonApi struct {
	api.Api
}

func (ca *CommonApi) ParseObjectFromRequest(c *gin.Context, model interface{}) (interface{}, error) {
	reader := io.Reader(c.Request.Body)
	decoder := json.NewDecoder(reader)

	err := decoder.Decode(model)
	if err != nil {
		ca.Logger.Errorf("Error reading data: %v", err)
	}
	return model, err
}

func (ca *CommonApi) ParseJsonFromRequest(c *gin.Context) (map[string]interface{}, error) {
	value, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(value), &data)

	return data, err
}

func (ca *CommonApi) Find(model interface{}, value interface{}, DB *gorm.DB) (interface{}, error) {
	query := types.JSON {
		"id" : value,
	}
	if err := DB.Model(model).Preload(clause.Associations).First(model, query).Error; err != nil {
		return model, err
	}
	return model, nil
}

func (ca *CommonApi) FindBy(model interface{}, field string, value interface{}, DB *gorm.DB) (interface{}, error) {
	query := types.JSON {
		field : value,
	}
	if err := DB.Model(model).Preload(clause.Associations).First(model, query).Error; err != nil {
		return model, err
	}
	return model, nil
}

func (ca *CommonApi) List(model interface{}, list interface{}, pagination *models.Pagination, DB *gorm.DB) (*models.Pagination, error) {
	if err := DB.Scopes(models.Paginate(model, list, pagination, DB)).Preload(clause.Associations).Find(list).Error; err != nil {
		return nil, err
	}
	pagination.Rows = list
	return pagination, nil
}

func (ca *CommonApi) Lookup(list interface{}, pagination *models.Pagination, filter types.JSON, DB *gorm.DB) (*models.Pagination, error) {
	var queryString string
	max := len(filter)
	i := 1
	for k, v := range filter {
		switch v.(type) {
		case string:
			queryString += fmt.Sprintf("%v LIKE \"%%%v%%\"", k, v)
		default:
			queryString += fmt.Sprintf("%v = %v", k, v)
		}
		if i < max {
			queryString += " AND "
		}
		i++
	}
	if err := DB.Scopes(models.Paginate(list, queryString, pagination, DB)).Where(queryString).Preload(clause.Associations).Find(list).Error; err != nil {
		return nil, err
	}
	pagination.Rows = list
	return pagination, nil
}

func (ca *CommonApi) Update(model interface{}, id interface{}, DB *gorm.DB) (interface{}, error) {
	query := types.JSON {
		"id" : id,
	}

	if err := DB.Model(model).Where(query).Updates(model).Error; err != nil {
		return model, err
	}
	return model, nil
}

func (ca *CommonApi) Delete(model interface{}, id interface{}, DB *gorm.DB) (interface{}, error) {
	query := types.JSON {
		"id" : id,
	}
	if err := DB.Model(model).Where(query).Delete(&model).Error; err != nil {
		return model, err
	}
	return model, nil
}