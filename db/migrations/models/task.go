package models

import (
	"fmt"
	"gorm.io/datatypes"
)

type Task struct {
	BaseModel
	Type 			string 			`json:"type" gorm:"type:varchar(100);index"`
	Name 			string 			`json:"name" gorm:"type:varchar(100);index"`
	UserID 			uint 			`json:"user_id" gorm:"index"`
	User 			User			`json:"user"`
	TargetID 		uint 			`json:"target_id" gorm:"index"`
	TargetEntity 	string 			`json:"target_entity" gorm:"type:text;index"`
	TargetTime 		string			`json:"target_time" gorm:"type:text;index"`
	Status			string 			`json:"status" gorm:"type:varchar(100);index"`
	Description 	string 			`json:"description" gorm:"type:text"`
	Data 			datatypes.JSON 	`json:"data"`
	Extra 			datatypes.JSON 	`json:"extra"`
}

func (Task) TableName() string {
	return "tasks"
}

func (t Task) Repr() string {
	return fmt.Sprintf("<TaskJob %v user = %v", t.ID, t.UserID)
}