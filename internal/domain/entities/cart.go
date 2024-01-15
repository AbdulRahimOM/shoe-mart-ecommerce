package entities

type Cart struct {
	UserID    uint `gorm:"notNull"`
	ProductID uint `gorm:"column:productId;notNull"`
	Quantity  uint `gorm:"column:quantity;notNull"`

	// FkUser    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`		//user details need not be provided along with cart items
	FkProduct Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

type CartProductPrice struct {
	ProductID uint    `json:"productID"`
	Quantity  uint    `json:"quantity"`
	Price     float32 `json:"price"`
}
