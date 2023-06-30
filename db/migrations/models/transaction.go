package models

import (
	"casorder/utils/types"
)

type Transaction struct {
	BaseModel
	OrderID			uint						`json:"order_id" gorm:"PRIMARY_KEY"`
	Code			string						`json:"code" gorm:"type:varchar(100)"`
	Type			types.TransactionType		`json:"type" gorm:"type:ENUM('TOP_UP', 'REFUND', 'WITHDRAW', 'ORDER');DEFAULT:'WITHDRAW'"`
	InvoiceNumber	string						`json:"invoice_number" gorm:"type:varchar(100)"`
	Discount		float64						`json:"discount" gorm:"type:Numeric(15, 2)"`
	DiscountRate	int64						`json:"discount_rate"`
	IncludedVAT		bool						`json:"included_vat" gorm:"default:true"`
	TaxCode			string						`json:"tax_code" gorm:"type:varchar(100)"`
	Quantity		float64						`json:"quantity" gorm:"type:Numeric(15, 2); not null"`
	PaymentType		types.PaymentType			`json:"payment_type" gorm:"type:ENUM('CASH', 'PAY_AS_YOU_GO', 'MONTHLY');DEFAULT:'CASH'"`
	Status			types.TransactionStatus		`json:"transaction_status" gorm:"type:ENUM('NEW', 'PENDING', 'DELETED', 'UPGRADING', 'REJECTED', 'USER_CONFIRMED', 'SALE_APPROVED', 'ADMIN_APPROVED');DEFAULT:'DELETED'"`
	Remark			string						`json:"remark" gorm:"type:text"`
	CustomerID		uint						`json:"customer_id" gorm:"not null"`
	StaffID			uint						`json:"staff_id" gorm:"not null"`
}

func (Transaction) TableName() string {
	return "transactions"
}
