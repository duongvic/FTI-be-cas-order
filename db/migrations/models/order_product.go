package models

import "gorm.io/datatypes"

type OrderProduct struct {
	BaseModel
	OrderID			uint 					`json:"order_id" gorm:"primarykey"`
	Idx				int 					`json:"idx" gorm:"primarykey;not null, default:1"`
	ProductID 		uint 					`json:"product_id" gorm:"primarykey;not null"`
	Product			Product					`json:"product"`
	IsPackage 		bool 					`json:"is_package" gorm:"primarykey;default:false"`
	Quantity 		int 					`json:"quantity" gorm:"not null"`
	Price 			float64 				`json:"price" gorm:"type:Numeric(15, 2)"`
	UnitID			uint 					`json:"unit_id" gorm:"not null"`
	Unit 			Unit					`json:"unit"`
	Disabled		bool					`json:"disabled"`
	Data			datatypes.JSON			`json:"data"`
}

func (OrderProduct) TableName() string {
	return "order_products"
}