package types

import "encoding/json"

// JSON alias type
type JSON = map[string]interface{}
//===========================================================

type EnumValue interface {
	Value() string
}
//===========================================================

//type AccountType string
//const (
//	AccountTypeAgency AccountType = "AGENCY"
//	AccountTypeEndUser AccountType = "EU"
//)
//func(a AccountType) Value() string {
//	return string(a)
//}
//===========================================================

type OrderStatus string
const (
	OrderStatusNew OrderStatus = "NEW"
	OrderStatusPending OrderStatus = "PENDING"
	OrderStatusRejected OrderStatus = "REJECTED"
	OrderStatusPayLater OrderStatus = "PAY_LATER"
	OrderStatusPayCompleted OrderStatus = "PAY_COMPLETED"
	OrderStatusPayIncomplete OrderStatus = "PAY_INCOMPLETE"
	OrderStatusDeployed = "DEPLOYED"
)
func (o OrderStatus) Value() string {
	return string(o)
}
//===========================================================

type OrderType string
const (
	OrderTypeBuy OrderType 		= "BUY"
	OrderTypeTrial OrderType 	= "TRIAL"
	OrderTypeUpgrade OrderType 	= "UPGRADE"
	OrderTypeExtend OrderType 	= "EXTEND"
	OrderTypeRenew OrderType	= "RENEW"
)
func (o OrderType) Value() string {
	return string(o)
}
//===========================================================

type PackageStatus string
const (
	PackageStatusActive PackageStatus 		= "ACTIVE"
	PackageStatusUpgrading PackageStatus	= "UPGRADING"
)
func (p PackageStatus) Value() string {
	return string(p)
}
//===========================================================

type PackageType string
const (
	PackageTypeStandard PackageType 	= "STANDARD"
	PackageTypeAdvanced PackageType 	= "ADVANCED"
	PackageTypeVIP PackageType 			= "VIP"
)
func (p PackageType) Value() string {
	return string(p)
}
//===========================================================

type PaymentMethod string
const (
	PaymentMethodCash PaymentMethod		= "CASH"
	PaymentMethodCredit	PaymentMethod	= "CREDIT"
	PaymentMethodVisa PaymentMethod		= "VISA"
)
func (p PaymentMethod) Value() string {
	return string(p)
}
//===========================================================

type PaymentType string
const (
	PaymentTypeCOD PaymentType			= "COD"
	PaymentTypePS	PaymentType			= "PAY_AS_YOU_GO"
	PaymentTypeMonthly PaymentType		= "MONTHLY"
)
func (p PaymentType) Value() string {
	return string(p)
}
//===========================================================

type PlatformType string
const (
	PlatformTypeOpenStack PlatformType = "OPENSTACK"
	PlatformTypePromox PlatformType = "PROMOX"
)
func(p PlatformType) Value() string {
	return string(p)
}
//===========================================================

type ProductType string
const (
	ProductTypeProcessor ProductType	= "PROCESSOR"
	ProductTypeMemory ProductType		= "MEMORY"
	ProductTypeStorage ProductType		= "STORAGE"
	ProductTypeNet ProductType			= "NET"
	ProductTypeOS ProductType			= "OS"
)
func (p ProductType) Value() string {
	return string(p)
}
//===========================================================

type ServiceType string
const (
	ServiceTypeCompute ServiceType 	= "COMPUTE"
	ServiceTypeVDC ServiceType 		= "POOL"
)
func (s ServiceType) Value() string {
	return string(s)
}
//===========================================================

type Status string
const (
	StatusNew Status 				= "NEW"
	StatusApproved Status			= "APPROVED"
	StatusPending Status 			= "PENDING"
	StatusDeleted Status 			= "DELETED"
	StatusRejected Status			= "REJECTED"
	StatusPayLater Status 			= "PAY_LATER"
	StatusPayCompleted Status 		= "PAY_COMPLETED"
	StatusPayIncomplete Status 		= "PAY_INCOMPLETE"
	StatusDeployed Status 			= "DEPLOYED"
)
func (s Status) Value() string {
	return string(s)
}
//===========================================================

type TaskStatus string
const (
	TaskStatusCreated TaskStatus = "CREATED"
	TaskStatusPending TaskStatus = "PENDING"
	TaskStatusCompleted TaskStatus = "COMPLETED"
	TaskStatusError TaskStatus = "ERROR"
	TaskStatusClosed TaskStatus = "CLOSED"
)
func(t TaskStatus) Value() string {
	return string(t)
}
//===========================================================

type TaskType string
const (
	TaskTypeOrder TaskType = "ORDER"
)
func (t TaskType) Value() string {
	return string(t)
}
//===========================================================

type TransactionStatus string
const (
	TransactionStatusNew TransactionStatus 			= "NEW"
	TransactionStatusPending TransactionStatus 		= "PENDING"
	TransactionStatusDeleted TransactionStatus 		= "DELETED"
	TransactionStatusUpgrading TransactionStatus 	= "UPGRADING"
	TransactionStatusRejected TransactionStatus 	= "REJECTED"

	TransactionStatusUserConfirmed TransactionStatus 	= "USER_CONFIRMED"
	TransactionStatusSaleApproved TransactionStatus 	= "SALE_APPROVED"
	TransactionStatusAdminApproved TransactionStatus	= "ADMIN_APPROVED"
)
func (t TransactionStatus) Value() string {
	return string(t)
}
//===========================================================

type TransactionType string
const (
	TransactionTypeTopUp TransactionType		= "TOP_UP"
	TransactionTypeRefund TransactionType		= "REFUND"
	TransactionTypeWithdraw TransactionType		= "WITHDRAW"
	TransactionTypeOrder TransactionType		= "ORDER"
)
func (t TransactionType) Value() string {
	return string(t)
}
//===========================================================

//type UserStatus string
//const (
//	UserStatusActive UserStatus = "ACTIVE"
//	UserStatusDeactivated UserStatus = "DEACTIVATED"
//	UserStatusBlocked UserStatus = "BLOCKED"
//)
//func(u UserStatus) Value() string {
//	return string(u)
//}
//===========================================================

type UnitType string
const (
	UnitKB UnitType 	= "KB"
	UnitMB UnitType 	= "MB"
	UnitGB UnitType 	= "GB"
	UnitTB UnitType 	= "TB"
	UnitVCPU UnitType 	= "vCPU"
	UnitIP UnitType 	= "IP"
	UnitMBPS UnitType 	= "Mbps"
	UnitKey UnitType 	= "KEY"
	UnitDay UnitType 	= "DAY"
	UnitMonth UnitType 	= "MONTH"
	UnitYear UnitType 	= "YEAR"
)
func(u UnitType) Value() string {
	return string(u)
}
//===========================================================

//type UserType string
//const (
//	UserTypePersonal UserType = "PERSONAL"
//	UserTypeCompany UserType = "COMPANY"
//)
//func(u UserType) Value() string {
//	return string(u)
//}
//===========================================================

func MapToEncodedJSON(m map[string]interface{}) []byte {
	data, _ := json.Marshal(m)
	return data
}