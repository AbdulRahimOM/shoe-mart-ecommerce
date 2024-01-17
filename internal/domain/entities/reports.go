package entities

type SalesReport struct {
	OrderCount          uint    `json:"order_count" gorm:"column:order_count"`
	CancelledOrderCount uint    `json:"cancelled_order_count" gorm:"column:cancelled_order_count"`
	ReturnedOrderCount  uint    `json:"returned_order_count" gorm:"column:returned_order_count"`
	CouponDiscounts     float32 `json:"coupon_discounts" gorm:"column:coupon_discounts"`
	NetOriginalValue    float32 `json:"net_original_value" gorm:"column:net_original_value"`
	NetSaleValue        float32 `json:"net_sale_value" gorm:"column:net_sale_value"`
	SaleValueAfterCancellationAndReturns float32 `json:"sale_value_after_cancellation_and_returns" gorm:"column:sale_value_after_cancellation_and_returns"`
	UsersRegistered	 uint    `json:"users_registered" gorm:"column:users_registered"`
	// UsersRegistered     uint    `json:"users_registered" gorm:"column:users_registered"`
	// SalePerDay		  []SalePerDay `json:"sale_per_day" gorm:"column:sale_per_day"`
}
type SalePerDay struct {
	Date        string  `json:"date" gorm:"column:date"`
	Sale  float32 `json:"sale" gorm:"column:sale"`
}