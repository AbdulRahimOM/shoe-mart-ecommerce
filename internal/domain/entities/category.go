package entities

type Categories struct {
	ID   uint   `gorm:"column:id;autoIncrement;primaryKey"`
	Name string `gorm:"column:name;notNull;unique"`
}
type CategoriesByAlphabet struct {
	Alphabet   string       `json:"alphabet"`
	Categories []Categories `json:"categories"`
}
type Brands struct {
	ID       uint   `gorm:"column:id;autoIncrement;primaryKey"`
	Name     string `gorm:"column:name;notNull"`
	SellerID uint   `gorm:"column:sellerId;notNull"`
}
type BrandsByAlphabet struct {
	Alphabet string   `json:"alphabet"`
	Brands   []Brands `json:"brands"`
}

type Models struct {
	ID         uint    `gorm:"column:id;autoIncrement;primaryKey"`
	Name       string  `gorm:"column:name;notNull"`
	BrandID    uint    `gorm:"column:brandId;notNull"`
	CategoryID uint    `gorm:"column:categoryId;notNull"`
	Rating     float32 `gorm:"column:rating"`

	FkBrand    Brands     `gorm:"foreignKey:BrandID;constraint:OnDelete:CASCADE"`
	FkCategory Categories `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
}
