package reportsrepo

import (
	"MyShoo/internal/domain/entities"

	"gorm.io/gorm"
)

type SalesReportRepo struct {
	DB *gorm.DB
}

func NewSalesReportRepository(db *gorm.DB) *SalesReportRepo {
	return &SalesReportRepo{DB: db}
}

func (repo *SalesReportRepo) GetSalesReport(start string, end string) (*entities.SalesReport, error) {

	// var orderCount uint
	// var netOrderValue float32
	var salesReport entities.SalesReport

	// err := repo.DB.Raw("SELECT COUNT(*) AS order_count, SUM(total) AS net_order_value FROM orders WHERE created_at BETWEEN ? AND ?", start, end).Scan(&orderCount, &netOrderValue).Error
	// if err != nil {
	// 	return &salesReport, err
	// }

	return &salesReport, nil
}
