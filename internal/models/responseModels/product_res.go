package response

import "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"

// get products response
type GetProductsResponse struct {
	//Status   string            `json:"status"`
	//Message  string            `json:"message"`
	//Error    string            `json:"error"`
	Products []ResponseProduct `json:"products"`
}

// responseProduct
type ResponseProduct struct {
	ID                     uint   `gorm:"column:id;autoIncrement;primaryKey"`
	SKUCode                string `gorm:"column:skuCode"`
	Name                   string `gorm:"column:name;notNull"`
	SizeIndex              uint   `gorm:"column:sizeIndex;notNull"`
	DimensionalVariationID uint   `gorm:"column:dimensionalVariationID;notNull"`
	Stock                  uint   `gorm:"column:stock;notNull"`

	FkDimensionalVariation struct {
		FkColourVariant struct {
			MRP       float32 `gorm:"column:mrp;notNull"`
			SalePrice float32 `gorm:"column:salePrice;notNull"`
		} `gorm:"column:fk_colourVariantID;notNull"`
	} `gorm:"column:fk_dimensionalVariationID;notNull"`
}

// get categories response
type GetCategoriesResponse struct {
	//Status     string                `json:"status"`
	//Message    string                `json:"message"`
	//Error      string                `json:"error"`
	Categories []entities.Categories `json:"categories"`
}

// get models response
type GetModelsResponse struct {
	//Status  string            `json:"status"`
	//Message string            `json:"message"`
	//Error   string            `json:"error"`
	Models []entities.Models `json:"models"`
}

// get brands response
type GetBrandsResponse struct {
	//Status           string                        `json:"status"`
	//Message          string                        `json:"message"`
	//Error            string                        `json:"error"`
	BrandsByAlphabet [26]entities.BrandsByAlphabet `json:"brandsByAlphabet"`
}

// get dimentional variations response
type GetDimensionalVariationsResponse struct { //Not used, may be SMED is being used instead
	//Status                string                          `json:"status"`
	//Message               string                          `json:"message"`
	//Error                 string                          `json:"error"`
	DimentionalVariations []entities.DimensionalVariation `json:"dimentionalVariations"`
}

// GetColourVariantsUnderModelResponse
type GetColourVariantsUnderModelResponse struct {
	//Status         string                  `json:"status"`
	//Message        string                  `json:"message"`
	//Error          string                  `json:"error"`
	ColourVariants []ResponseColourVarient `json:"colourVariants"`
}

type ResponseColourVarient struct {
	ID        uint    `gorm:"column:id;autoIncrement;primaryKey"`
	Colour    string  `gorm:"column:colour;notNull"`
	ModelID   uint    `gorm:"column:modelId;notNull"`
	MRP       float32 `gorm:"column:mrp;notNull"`
	SalePrice float32 `gorm:"column:salePrice;notNull"`
}

// PQR //R for sale price
type PQR struct {
	ProductID        uint    `gorm:"column:product_id;notNull"`
	Quantity         uint    `gorm:"column:quantity;notNull"`
	SalePriceOnOrder float32 `gorm:"column:sale_price_on_order;notNull"`
}
