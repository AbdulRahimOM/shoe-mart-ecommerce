package entities

type DashboardData struct {
	OrderCount                           uint    `json:"order_count" gorm:"column:order_count"`
	CancelledOrderCount                  uint    `json:"cancelled_order_count" gorm:"column:cancelled_order_count"`
	ReturnedOrderCount                   uint    `json:"returned_order_count" gorm:"column:returned_order_count"`
	CouponDiscounts                      float32 `json:"coupon_discounts" gorm:"column:coupon_discounts"`
	NetOriginalValue                     float32 `json:"net_original_value" gorm:"column:net_original_value"`
	NetSaleValue                         float32 `json:"net_sale_value" gorm:"column:net_sale_value"`
	SaleValueAfterCancellationAndReturns float32 `json:"sale_value_after_cancellation_and_returns" gorm:"column:sale_value_after_cancellation_and_returns"`
	UsersRegistered                      uint    `json:"users_registered" gorm:"column:users_registered"`
	// UsersRegistered     uint    `json:"users_registered" gorm:"column:users_registered"`
	// SalePerDay		  []SalePerDay `json:"sale_per_day" gorm:"column:sale_per_day"`
}
type SalePerDay struct {
	Date string  `json:"date" gorm:"column:date"`
	Sale float32 `json:"sale" gorm:"column:sale"`
}

type SalesReportOrderList struct {
	Date     string `json:"date" gorm:"column:date"`
	OrderID  uint   `json:"order_id" gorm:"column:order_id"`
	SellerID uint   `json:"seller_id" gorm:"column:seller_id"`
	// SellerName     string  `json:"seller_name" gorm:"column:seller_name"`
	BrandName   string  `json:"brand_name" gorm:"column:brand_name"`
	ModelName   string  `json:"model_name" gorm:"column:model_name"`
	ProductName string  `json:"product_name" gorm:"column:product_name"`
	SKU         string  `json:"sku" gorm:"column:sku"`
	Quantity    uint    `json:"quantity" gorm:"column:quantity"`
	MRP         float32 `json:"mrp" gorm:"column:mrp"`
	// Discount 	  float32 `json:"discount" gorm:"column:discount"`
	SalePrice float32 `json:"sale_price" gorm:"column:sale_price"`
}

type SellerWiseReport struct {
	SellerID               uint    `json:"seller_id" gorm:"column:seller_id"`
	SellerName             string  `json:"seller_name" gorm:"column:seller_name"`
	QuantityCount          uint    `json:"quantity_count" gorm:"column:quantity_count"`
	TotalSale              float32 `json:"total_sale" gorm:"column:total_sale"`
	ReturnQuantityCount    uint    `json:"return_quantity_count" gorm:"column:return_quantity_count"`
	TotalReturnValue       float32 `json:"total_return_value" gorm:"column:total_return_value"`
	CancelQuantityCount    uint    `json:"cancel_quantity_count" gorm:"column:cancel_quantity_count"`
	TotalCancelValue       float32 `json:"total_cancel_value" gorm:"column:total_cancel_value"`
	EffectiveQuantityCount uint    `json:"effective_quantity_count" gorm:"column:effective_quantity_count"`
	EffectiveSaleValue     float32 `json:"effective_sale_value" gorm:"column:effective_sale_value"`
}

type BrandWiseReport struct {
	BrandID                uint    `json:"brand_id" gorm:"column:brand_id"`
	BrandName              string  `json:"brand_name" gorm:"column:brand_name"`
	QuantityCount          uint    `json:"quantity_count" gorm:"column:quantity_count"`
	TotalSale              float32 `json:"total_sale" gorm:"column:total_sale"`
	ReturnQuantityCount    uint    `json:"return_quantity_count" gorm:"column:return_quantity_count"`
	TotalReturnValue       float32 `json:"total_return_value" gorm:"column:total_return_value"`
	CancelQuantityCount    uint    `json:"cancel_quantity_count" gorm:"column:cancel_quantity_count"`
	TotalCancelValue       float32 `json:"total_cancel_value" gorm:"column:total_cancel_value"`
	EffectiveQuantityCount uint    `json:"effective_quantity_count" gorm:"column:effective_quantity_count"`
	EffectiveSaleValue     float32 `json:"effective_sale_value" gorm:"column:effective_sale_value"`
}

type ModelWiseReport struct {
	ModelID                uint    `json:"model_id" gorm:"column:model_id"`
	ModelName              string  `json:"model_name" gorm:"column:model_name"`
	BrandName              string  `json:"brand_name" gorm:"column:brand_name"`
	QuantityCount          uint    `json:"quantity_count" gorm:"column:quantity_count"`
	TotalSale              float32 `json:"total_sale" gorm:"column:total_sale"`
	ReturnQuantityCount    uint    `json:"return_quantity_count" gorm:"column:return_quantity_count"`
	CancelQuantityCount    uint    `json:"cancel_quantity_count" gorm:"column:cancel_quantity_count"`
	EffectiveQuantityCount uint    `json:"effective_quantity_count" gorm:"column:effective_quantity_count"`
	EffectiveSaleValue     float32 `json:"effective_sale_value" gorm:"column:effective_sale_value"`
}

type SizeWiseReport struct {
	SizeIndex     uint   `json:"size_index" gorm:"column:size_index"`
	SizeName      string `json:"size_name" gorm:"column:size_name"`
	QuantityCount uint   `json:"quantity_count" gorm:"column:quantity_count"`
}

type RevenueGraph struct {
	Date     string  `json:"date" gorm:"column:date"`
	Sale     float32 `json:"sale" gorm:"column:sale"`
	Quantity uint    `json:"quantity" gorm:"column:quantity"`
}
