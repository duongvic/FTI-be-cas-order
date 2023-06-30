package models

type PackageProduct struct {
	BaseModel
	Package				Package
	PackageID 			uint 		`json:"package_id" gorm:"primarykey;not null"`
	Product				Product		`json:"product"`
	ProductID 			uint 		`json:"product_id" gorm:"primarykey;not null"`
	Quantity			int 		`json:"quantity" gorm:"default:0"`
	InitFee			 	float64 	`json:"init_fee" gorm:"type:Numeric(15, 2);DEFAULT:0"`
	MaintenanceFee 		float64 	`json:"maintenance_fee" gorm:"type:Numeric(15, 2);DEFAULT:0"`
	UnitID 				uint 		`json:"unit_id" gorm:"not null"`
	Unit 				Unit		`json:"unit"`
}

func (PackageProduct) TableName() string {
	return "package_products"
}
