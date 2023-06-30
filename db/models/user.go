package models

import "casorder/utils/types"

// User data model
type User struct {
	BaseModel
	Username 	string			`json:"username" gorm:"size:60;unique;not null"`
	MiqID	 	int64			`json:"miq_id"`
	Role	 	string			`json:"role"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) Read(m types.JSON) {
	u.Username = m["userid"].(string)
	u.MiqID	   = m["id"].(int64)
	u.Role	   = m["role"].(string)
}
