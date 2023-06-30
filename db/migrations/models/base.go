package models

import (
	"gorm.io/gorm"
	"math"
	"time"
)

type BaseModel struct {
	ID        		uint 				`json:"id" gorm:"primarykey"`
	CreatedAt 		time.Time 			`json:"created_at"`
	UpdatedAt 		time.Time 			`json:"updated_at"`
	DeletedAt 		gorm.DeletedAt		`json:"deleted_at" gorm:"index"`
	Deleted   		bool		      	`json:"deleted"`
}


type ModelMixin interface {
}

type Pagination struct {
	PageSize   int         `json:"pageSize,omitempty;query:pageSize"`
	Page       int         `json:"page,omitempty;query:page"`
	Sort       string      `json:"sort,omitempty;query:sort"`
	TotalRows  int64       `json:"totalRows"`
	TotalPages int         `json:"totalPages"`
	Rows       interface{} `json:"rows"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

func (p *Pagination) GetPageSize() int {
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	return p.PageSize
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "id DESC"
	}
	return p.Sort
}

func Paginate(model interface{}, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(model).Count(&totalRows)
	pagination.TotalRows = totalRows
	pagination.PageSize = pagination.GetPageSize()
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.PageSize)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.PageSize).Order(pagination.GetSort())
	}
}
