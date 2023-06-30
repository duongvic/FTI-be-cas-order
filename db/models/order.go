package models

import (
	"casorder/utils/types"
	"gorm.io/datatypes"
	"math"
	"time"
)

// Order data model
type Order struct {
	BaseModel
	Code          string            `json:"code" gorm:"type:varchar(100);index;not null;unique"`
	Remark        string            `json:"remark"`
	OrderType     types.OrderType   `json:"order_type" gorm:"type:ENUM('BUY', 'TRIAL', 'UPGRADE', 'EXTEND', 'RENEW');DEFAULT:'TRIAL'"`
	ServiceType   types.ServiceType `json:"service_type" gorm:"type:ENUM('COMPUTE', 'POOL');DEFAULT:'COMPUTE'"`
	Platform      string            `json:"platform"`
	Status        types.Status      `json:"status" gorm:"type:ENUM('NEW', 'APPROVED', 'PENDING', 'DELETED', 'REJECTED', 'PAY_LATER', 'PAY_COMPLETED', 'PAY_INCOMPLETE', 'DEPLOYED');DEFAULT:'PENDING'"`
	ApprovalStep  int               `json:"approval_step" gorm:"default:0"`
	Duration      int64             `json:"duration" gorm:"not null; default:0"`
	RegionID      uint              `json:"region_id" gorm:"not null"`
	Region        Region            `json:"region"`
	Vouchers      datatypes.JSON    `json:"vouchers"`
	Price         float64           `json:"price" gorm:"type:Numeric(15, 2);DEFAULT:0;not null"`
	Discount      float64           `json:"discount" gorm:"type:Numeric(15, 2)"`
	VATFee        float64           `json:"vat_fee" gorm:"type:Numeric(15, 2)"`
	PaymentType   types.PaymentType `json:"pmt_type" gorm:"type:ENUM('COD', 'PAY_AS_YOU_GO', 'MONTHLY');DEFAULT:'COD'"`
	ContractCode  string            `json:"contract_code" gorm:"type:varchar(100)"`
	StartAt       time.Time         `json:"start_at"`
	EndAt         time.Time         `json:"end_at"`
	OrderDtlID    uint              `json:"order_dtl_id" gorm:"not null"`
	OrderDtl      OrderDtl          `json:"order_dtl"`
	RefOrderID    uint              `json:"ref_order_id"`
	RefOrderIdx   int64             `json:"ref_order_idx"`
	Lock          bool              `json:"lock"`
	LockReason    string            `json:"lock_reason"`
	CustomerID    uint              `json:"customer_id" gorm:"not null"`
	SaleCare      string            `json:"sale_care" gorm:"type:varchar(100)"`
	CoSale        datatypes.JSON    `json:"co_sale"`
	StaffID       uint              `json:"staff_id" gorm:"not null"`
	OrderProducts []*OrderProduct   `json:"order_products"`
	Transactions  []*Transaction    `json:"transactions"`
	Liquidated    bool              `json:"liquidated"`
}

func (Order) TableName() string {
	return "orders"
}

func (o Order) Total() float64 {
	return math.Round((o.Price-o.VATFee-o.Discount)*100) / 100
}
