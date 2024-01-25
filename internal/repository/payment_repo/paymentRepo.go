package paymentrepo

import (
	repository_interface "MyShoo/internal/repository/interface"

	"gorm.io/gorm"
)

type PaymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) repository_interface.IPaymentRepo {
	return &PaymentRepo{
		db: db,
	}
}
