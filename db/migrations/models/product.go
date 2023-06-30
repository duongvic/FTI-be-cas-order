package models

import (
	"casorder/utils/types"
	"gorm.io/datatypes"
)

type Product struct {
	BaseModel
	Name 				string 					`json:"name" gorm:"type:varchar(100);not null"`
	CN 					string 					`json:"cn" gorm:"type:varchar(50);not null"`
	Description 		string 					`json:"description" gorm:"type:text"`
	Type 				types.ProductType 		`json:"type" gorm:"type:ENUM('PROCESSOR', 'MEMORY', 'STORAGE', 'NET', 'OS')"`
	ServiceType			types.ServiceType		`json:"service_type" gorm:"type:ENUM('COMPUTE', 'POOL');DEFAULT:'COMPUTE'"`
	Platform 			string 					`json:"platform" gorm:"type:varchar(100)"`
	InitFee 			float64 				`json:"init_fee" gorm:"type:Numeric(15, 2);DEFAULT:0"`
	MaintenanceFee 		float64 				`json:"maintenance_fee" gorm:"type:Numeric(15, 2);DEFAULT:0"`
	UnitID 				uint 					`json:"unit_id" gorm:"not null"`
	Unit				Unit					`json:"unit"`
	IsBase				bool					`json:"is_base" gorm:"not null"`
	Data				datatypes.JSON			`json:"data"`
}

func (Product) TableName() string {
	return "products"
}
