package entities

// type ProductImages struct {
// 	VariantID string `gorm:"column:variantId;notNull"`
// 	ImageURL  string `gorm:"column:imageUrl"`
// }

type ColourVariant struct {
	ID        uint    `gorm:"column:id;autoIncrement;primaryKey"`
	Colour    string  `gorm:"column:colour;notNull"`
	ModelID   uint    `gorm:"column:modelId;notNull"`
	MRP       float32 `gorm:"column:mrp;notNull"`
	SalePrice float32 `gorm:"column:salePrice;notNull"`

	FkModel Models `gorm:"foreignKey:ModelID;constraint:OnDelete:CASCADE"`
}

type DimensionalVariant struct {
	ID              uint `json:"id" gorm:"column:id;autoIncrement;primaryKey"`
	ColourVariantID uint `json:"colourVariantId" gorm:"column:colourVariantId;notNull"`
	DVIndex         uint `json:"dvIndex" gorm:"column:dvIndex;notNull"`
	// ColourVariant  ColourVariant `gorm:"-"`
	FkColourVariant ColourVariant `gorm:"foreignKey:ColourVariantID;constraint:OnDelete:CASCADE"`
}

type Product struct {
	ID                     uint   `gorm:"column:id;autoIncrement;primaryKey"`
	SKUCode                string `gorm:"column:skuCode"`
	Name                   string `gorm:"column:name;notNull"`
	SizeIndex              uint   `gorm:"column:sizeIndex;notNull"`
	DimensionalVariationID uint   `gorm:"column:dimensionalVariationID;notNull"`
	Stock                  uint   `gorm:"column:stock;notNull"`

	FkDimensionalVariation DimensionalVariant `gorm:"foreignKey:DimensionalVariationID;constraint:OnDelete:CASCADE"`
}

func (Product) TableName() string {
	return "product"
}

type SizeVariations struct {
	Index int
	Size  string //`gorm:"column:size;notNull"`
}

// "6", "6.5", "7", "7.5", "8", "8.5", "9", "9.5", "10", "10.5", "11", "11.5", "12", "12.5", "13", "13.5", "14", "14.5", "15"}
var Size = []SizeVariations{
	{0, "6"},
	{1, "6.5"},
	{2, "7"},
	{3, "7.5"},
	{4, "8"},
	{5, "8.5"},
	{6, "9"},
	{7, "9.5"},
	{8, "10"},
	{9, "10.5"},
	{10, "11"},
	{11, "11.5"},
	{12, "12"},
	{13, "12.5"},
	{14, "13"},
	{15, "13.5"},
	{16, "14"},
	{17, "14.5"},
	{18, "15"},
}

type DimensionalVariation struct {
	Index uint   //`gorm:"column:id;autoIncrement;primaryKey"`
	Name  string //`gorm:"column:name;notNull"`
	Code  string //`gorm:"column:code;notNull;primaryKey"`
}

// "Extra-Extra-Narrow","Extra-Narrow","Narrow","Medium","Wide","Extra-Wide","Extra-Extra-Wide"
var DimensionalVariations = []DimensionalVariation{
	{0, "Extra-Extra-Narrow", "XXN"},
	{1, "Extra-Narrow", "XN"},
	{2, "Narrow", "N"},
	{3, "Medium", "M"},
	{4, "Wide", "W"},
	{5, "Extra-Wide", "XW"},
	{6, "Extra-Extra-Wide", "XXW"},
}