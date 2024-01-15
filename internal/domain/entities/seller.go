package entities

type Seller struct {
	ID        uint   `gorm:"primaryKey"`
	FirstName string `gorm:"column:firstName;notNull"`
	LastName  string `gorm:"column:lastName"`
	Email     string `gorm:"column:email;unique;notNull"`
	Phone     string `gorm:"column:phone;notNull"`
	Password  string `gorm:"-"`
	Status    string `gorm:"default:Pending"`
}

type SellerDetails struct {
	ID        uint   `gorm:"primaryKey"`
	FirstName string `gorm:"column:firstName;notNull"`
	LastName  string `gorm:"column:lastName"`
	Email     string `gorm:"unique;notNull"`
	Phone     string `gorm:"notNull"`
	Status    string `gorm:"default:Pending"`
}
