package entities

type Seller struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	FirstName string `json:"firstName" gorm:"column:firstName;notNull"`
	LastName  string `json:"lastName" gorm:"column:lastName"`
	Email     string `json:"email" gorm:"column:email;unique;notNull"`
	Phone     string `json:"phone" gorm:"column:phone;notNull"`
	Password  string `json:"-" gorm:"column:password"`
	Status    string `json:"status" gorm:"default:Pending"`
}

type PwMaskedSeller struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	FirstName string `json:"first_name" gorm:"column:firstName;notNull"`
	LastName  string `json:"last_name" gorm:"column:lastName"`
	Email     string `json:"email" gorm:"unique;notNull"`
	Phone     string `json:"phone" gorm:"notNull"`
	Status    string `json:"status" gorm:"default:Pending"`
}
