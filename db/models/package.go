package models

import (
	"casorder/utils/types"
	"math"
)

type Package struct {
	BaseModel
	Name 				string 					`json:"name" gorm:"type:varchar(100);not null"`
	Description 		string 					`json:"description" gorm:"type:text"`
	Type 				types.PackageType 		`json:"package_type" gorm:"type:ENUM('STANDARD', 'ADVANCED', 'VIP');DEFAULT:'STANDARD'"`
	Status 				types.PackageStatus 	`json:"package_status" gorm:"type:ENUM('ACTIVE', 'UPGRADING');DEFAULT:'ACTIVE'"`
	ServiceType			types.ServiceType		`json:"service_type" gorm:"type:ENUM('COMPUTE', 'POOL');DEFAULT:'COMPUTE'"`
	Platform 			string 					`json:"platform" gorm:"type:varchar(100)"`
	TrialTime 			int 					`json:"trial_time" gorm:"default:0"`
	AllowCustom 		bool 					`json:"allow_custom"`
	InitFee 			float64 				`json:"init_fee" gorm:"type:Numeric(15, 2);DEFAULT:0"`
	MaintenanceFee 		float64 				`json:"maintenance_fee" gorm:"type:Numeric(15, 2);DEFAULT:0"`
	PaymentPeriod 		int 					`json:"payment_period"`
	UnitID 				uint 					`json:"unit_id" gorm:"not null"`
	Unit 				Unit					`json:"unit"`
	Priority 			int 					`json:"priority"`
}

func (Package) TableName() string {
	return "packages"
}

func (p Package) Price() float64 {
	return math.Round((p.InitFee + p.MaintenanceFee)*100)/100
}