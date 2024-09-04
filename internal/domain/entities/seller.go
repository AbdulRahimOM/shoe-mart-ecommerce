package entities

type Seller struct {
	ID        uint   `gorm:"primaryKey"`
	FirstName string `gorm:"column:firstName;notNull"`
	LastName  string `gorm:"column:lastName"`
	Email     string `gorm:"column:email;unique;notNull"`
	Phone     string `gorm:"column:phone;notNull"`
	Password  string `gorm:"column:password"`
	Status    string `gorm:"default:Pending"`
}

type PwMaskedSeller struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	FirstName string `json:"first_name" gorm:"column:firstName;notNull"`
	LastName  string `json:"last_name" gorm:"column:lastName"`
	Email     string `json:"email" gorm:"unique;notNull"`
	Phone     string `json:"phone" gorm:"notNull"`
	Status    string `json:"status" gorm:"default:Pending"`
}
