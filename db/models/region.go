package models

type Region struct {
	BaseModel
	Type 			string 		`json:"type" gorm:"type:varchar(100)"`
	Name 			string 		`json:"name" gorm:"type:varchar(100);not null"`
	Status 			bool 		`json:"status"`
	Description 	string 		`json:"description" gorm:"type:text"`
	Address 		string 		`json:"address" gorm:"type:text"`
	City 			string 		`json:"city" gorm:"type:varchar(100)"`
	CountryCode 	string 		`json:"country_code" gorm:"type:varchar(50)"`
}

func (Region) TableName() string {
	return "regions"
}
