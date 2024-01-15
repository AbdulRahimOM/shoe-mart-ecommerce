package entities

type WishList struct {
	ID     uint   `gorm:"notNull;primaryKey"`
	Name   string `gorm:"notNull"`
	UserID uint   `gorm:"notNull"`
}

type WishListItems struct {
	WishListID uint `gorm:"notNull"`
	ProductID  uint `gorm:"notNull"`

	// FkWishList WishList `gorm:"foreignKey:WishListID;constraint:OnDelete:CASCADE"`
	FkProduct Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}
