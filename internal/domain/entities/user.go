package entities

import "time"

type User struct {
	ID            uint      `gorm:"primaryKey"`
	FirstName     string    `gorm:"column:firstName;notNull"`
	LastName      string    `gorm:"column:lastName"`
	Email         string    `gorm:"column:email;unique;notNull"`
	Phone         string    `gorm:"column:phone;notNull"`
	Password      string    `gorm:"column:password;notNull"`
	Status        string    `gorm:"default:Pending;notNull"`
	CreatedAt     time.Time `gorm:"column:created_at;notNull"`
	WalletBalance float32   `gorm:"column:wallet_balance;notNull;default:0"`
}

type UserDetails struct {
    ID        uint   `json:"id" gorm:"primaryKey"`
    FirstName string `json:"first_name" gorm:"column:firstName;notNull"`
    LastName  string `json:"last_name" gorm:"column:lastName"`
    Email     string `json:"email" gorm:"column:email;unique;notNull"`
    Phone     string `json:"phone" gorm:"column:phone;notNull"`
    Status    string `json:"status" gorm:"default:Pending"`
}


type UserAddress struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"column:userId;notNull"`
	AddressName string `gorm:"column:addressName;notNull"`
	FirstName   string `gorm:"column:firstName;notNull"`
	LastName    string `gorm:"column:lastName"`
	Email       string `gorm:"column:email;notNull"`
	Phone       string `gorm:"column:phone;notNull"` //with country code
	Street      string `gorm:"column:street;notNull"`
	LandMark    string `gorm:"column:landmark"`
	City        string `gorm:"column:city;notNull"`
	State       string `gorm:"column:state;notNull"`
	Pincode     uint   `gorm:"column:pincode;notNull"`

	// Country     string `gorm:"column:country;notNull"`	India only

	// FkUser UserDetails `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
