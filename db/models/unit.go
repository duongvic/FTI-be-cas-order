package models

type Unit struct {
	BaseModel
	Name 			string 		`json:"name" gorm:"type:varchar(100);not null"`
	Code 			string 		`json:"code" gorm:"type:varchar(100);unique"`
	Description 	string 		`json:"description" gorm:"type:text"`
}

func (Unit) TableName() string {
	return "units"
}
