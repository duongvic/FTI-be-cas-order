package services

import (
	"casorder/db/models"
	"casorder/utils/types"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductService struct {
	CommonApi
}

func (p *ProductService) Lookup(list interface{}, pagination *models.Pagination, filter types.JSON, DB *gorm.DB) (*models.Pagination, error) {
	var queryString = "cn <> 'no_os'"
	max := len(filter)
	if max >= 1 {
		queryString += " AND "
	}
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