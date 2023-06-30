package models

// User data model
type User struct {
	BaseModel
	Username 	string			`json:"username" gorm:"size:60;unique;not null"`
	MiqID	 	uint			`json:"miq_id"`
	Role	 	string			`json:"role"`
}

func (User) TableName() string {
	return "users"
}

