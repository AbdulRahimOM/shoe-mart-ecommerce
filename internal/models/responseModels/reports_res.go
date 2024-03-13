package response

import "MyShoo/internal/domain/entities"

type GetDashBoardDataResponse struct {
	//Status        string                 `json:"status"`
	//Message       string                 `json:"message"`
	//Error         string                 `json:"error"`
	DashboardData entities.DashboardData `json:"sales_report"` //need update
	SalePerDay    []entities.SalePerDay  `json:"sale_per_day"`
}

type TopProductsResponse struct {
	//Status      string        `json:"status"`
	Limit 	 int           `json:"limit"`
	TopProducts []TopProducts `json:"topProducts"`
}

type TopProducts struct {
	ColourVariantID  uint    `json:"colour_variant_id" gorm:"column:colour_variant_id"`
	ModelName        string  `json:"model_name" gorm:"column:model_name"`
	BrandName        string  `json:"brand_name" gorm:"column:brand_name"`
	SellerName       string  `json:"seller_name" gorm:"column:seller_name"`
	QuantitySold     uint    `json:"quantity_sold" gorm:"column:quantity_sold"`
	NetMRPValue      float32 `json:"net_mrp_value" gorm:"column:net_mrp_value"`
	NetSaleAmount    float32 `json:"net_sale_amount" gorm:"column:net_sale_amount"`
	CurrentMRP       float32 `json:"current_mrp" gorm:"column:current_mrp"`
	CurrentSalePrice float32 `json:"current_sale_price" gorm:"column:current_sale_price"`
}

type TopSellersResponse struct {
	//Status     string         `json:"status"`
	Limit 	 int           `json:"limit"`
	TopSellers []TopSellers `json:"topSellers"`
}

type TopSellers struct {
	SellerID        uint    `json:"seller_id" gorm:"column:seller_id"`
	SellerName      string  `json:"seller_name" gorm:"column:seller_name"`
	QuantitySold    uint    `json:"quantity_sold" gorm:"column:quantity_sold"`
	NetSaleAmount   float32 `json:"net_sale_amount" gorm:"column:net_sale_amount"`
	NetMRPValueSold float32 `json:"net_mrp_value_sold" gorm:"column:net_mrp_value_sold"`
}

type TopBrandsResponse struct {
	//Status    string        `json:"status"`
	Limit 	 int           `json:"limit"`
	TopBrands []TopBrands `json:"topBrands"`
}

type TopBrands struct {
	BrandID         uint    `json:"brand_id" gorm:"column:brand_id"`
	BrandName       string  `json:"brand_name" gorm:"column:brand_name"`
	SellerName      string  `json:"seller_name" gorm:"column:seller_name"`
	QuantitySold    uint    `json:"quantity_sold" gorm:"column:quantity_sold"`
	NetSaleAmount   float32 `json:"net_sale_amount" gorm:"column:net_sale_amount"`
	NetMRPValueSold float32 `json:"net_mrp_value_sold" gorm:"column:net_mrp_value_sold"`
}

type TopModelsResponse struct {
	//Status    string        `json:"status"`
	Limit 	 int           `json:"limit"`
	TopModels []TopModels `json:"topModels"`
}

type TopModels struct {
	ModelID         uint    `json:"model_id" gorm:"column:model_id"`
	ModelName       string  `json:"model_name" gorm:"column:model_name"`
	BrandName       string  `json:"brand_name" gorm:"column:brand_name"`
	SellerName      string  `json:"seller_name" gorm:"column:seller_name"`
	QuantitySold    uint    `json:"quantity_sold" gorm:"column:quantity_sold"`
	NetSaleAmount   float32 `json:"net_sale_amount" gorm:"column:net_sale_amount"`
	NetMRPValueSold float32 `json:"net_mrp_value_sold" gorm:"column:net_mrp_value_sold"`
}
