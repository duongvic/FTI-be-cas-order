package services

import (
	"casorder/db/models"
	"casorder/utils/types"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderProductService struct {
	CommonApi
}


func (ops *OrderProductService) Find(model interface{}, value interface{}, DB *gorm.DB) (interface{}, error) {
	query := types.JSON {
		"id" : value,
	}
	if err := DB.Model(model).Preload(clause.Associations).Preload("Product.Unit").First(model, query).Error; err != nil {
		return model, err
	}
	return model, nil
}

func (ops *OrderProductService) FindBy(model interface{}, field string, value interface{}, DB *gorm.DB) (interface{}, error) {
	query := types.JSON {
		field : value,
	}
	if err := DB.Model(model).Preload(clause.Associations).Preload("Product.Unit").First(model, query).Error; err != nil {
		return model, err
	}
	return model, nil
}

func (ops *OrderProductService) List(model interface{}, list interface{}, pagination *models.Pagination, DB *gorm.DB) (*models.Pagination, error) {
	if err := DB.Scopes(models.Paginate(model, list, pagination, DB)).Preload(clause.Associations).Preload("Product.Unit").Find(list).Error; err != nil {
		return nil, err
	}
	pagination.Rows = list
	return pagination, nil
}

func (ops *OrderProductService) Lookup(list interface{}, pagination *models.Pagination, filter types.JSON, DB *gorm.DB) (*models.Pagination, error) {
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
	if err := DB.Scopes(models.Paginate(list, queryString, pagination, DB)).Where(queryString).Preload(clause.Associations).Preload("Product.Unit").Find(list).Error; err != nil {
		return nil, err
	}
	pagination.Rows = list
	return pagination, nil
}