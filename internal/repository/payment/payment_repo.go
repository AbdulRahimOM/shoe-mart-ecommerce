package paymentrepo

import (
	repo "MyShoo/internal/repository/interface"

	"gorm.io/gorm"
)

type PaymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) repo.IPaymentRepo {
	return &PaymentRepo{
		db: db,
	}
}
