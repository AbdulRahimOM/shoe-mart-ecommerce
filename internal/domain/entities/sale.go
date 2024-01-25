package entities

import "time"

//	Table Orders{
//		ID UUID [primary key]
//		Reference_no CHAR(13)
//		OrderDate TimeStamp
//		UserID UUID
//		DeliveredDate TimeStamp
//		TotalAmount DECIMAL(10,2)
//		Discount DECIMAL(10,2)
//		AppliedCoupon UUID
//		Payment_Method ENUM('COD', 'UPI','BANK')
//		Status ENUM ('PROCESSING','DELIVERED','CANCELLED')
//		AddressID UUID
//	  }
type Order struct {
	ID               uint      `gorm:"primaryKey"`
	ReferenceNo      string    `gorm:"column:reference_no;notNull"`
	OrderDateAndTime time.Time `gorm:"column:order_date_and_time;notNull"`
	UserID           uint      `gorm:"column:user_id;notNull"`
	DeliveredDate    string    `gorm:"column:delivered_date;notNull"`
	OriginalAmount   float32   `gorm:"column:original_amount;notNull"`
	CouponDiscount   float32   `gorm:"column:coupon_discount;notNull"`
	FinalAmount      float32   `gorm:"column:final_amount;notNull"`
	CouponID         uint      `gorm:"column:applied_coupon;notNull"`
	PaymentMethod    string    `gorm:"column:payment_method;notNull"`
	Status           string    `gorm:"column:status;notNull"`
	AddressID        uint      `gorm:"column:address_id;notNull"`
	PaymentStatus    string    `gorm:"column:payment_status"` //need update to notNull
	TransactionID    string    `gorm:"column:transaction_id"` //need update to notNull

	// FkUser    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	FkAddress UserAddress `gorm:"foreignKey:AddressID;constraint:OnDelete:CASCADE"`
}

type OrderItem struct {
	OrderID          uint    `gorm:"column:order_id;notNull"`
	ProductID        uint    `gorm:"column:product_id;notNull"`
	Quantity         uint    `gorm:"column:quantity;notNull"`
	SalePriceOnOrder float32 `gorm:"column:sale_price_on_order;notNull"`

	// FkOrder   Order   `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	FkProduct Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}
type OrderInfo struct {
	OrderDetails Order `json:"orderDetails"`
	OrderItems   []PQ  `json:"orderItems"`
}
type DetailedOrderInfo struct {
	OrderDetails Order       `json:"orderDetails"`
	OrderItems   []OrderItem `json:"orderItems"`
}

var PaymentMethod = []string{"COD", "ONLINE", "WALLET"}

type PQ struct {
	ProductID uint `gorm:"column:product_id;notNull"`
	Quantity  uint `gorm:"column:quantity;notNull"`
}

// type PQS struct {
// 	ProductID uint `json:"productID" validate:"required,number"`
// 	Quantity  uint `json:"quantity" validate:"required,number"`
// 	SalePriceOnOrder float32 `json:"salePriceOnOrder" validate:"required,number"`
// }
