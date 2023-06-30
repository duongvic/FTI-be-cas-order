package models

import (
	"time"
)

type OrderDtl struct {
	BaseModel
	Name				string		`json:"name" gorm:"type:varchar(100)"`
	Email				string		`json:"email" gorm:"type:varchar(100)"`
	PhoneNumber			string		`json:"phone_number" gorm:"type:varchar(100)"`
	TaxNumber			string		`json:"tax_number" gorm:"type:varchar(100)"`
	IdNumber			string		`json:"id_number" gorm:"type:varchar(100)"`
	IdIssueDate			time.Time	`json:"id_issue_date"`
	IdIssueLocation		string		`json:"id_issue_location" gorm:"type:varchar(100)"`
	Address				string		`json:"address" gorm:"type:text"`
	RepName				string		`json:"rep_name" gorm:"type:varchar(100)"`
	RepPhone			string		`json:"rep_phone" gorm:"type:varchar(100)"`
	RepEmail			string		`json:"rep_email" gorm:"type:varchar(100)"`
	RefName				string		`json:"ref_name" gorm:"type:varchar(100)"`
	RefPhone			string		`json:"ref_phone" gorm:"type:varchar(100)"`
	RefEmail			string		`json:"ref_email" gorm:"type:varchar(100)"`
}

func (OrderDtl) TableName() string {
	return "order_dtls"
}
