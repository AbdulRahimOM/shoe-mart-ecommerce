package usecaseInterface

import (
	"MyShoo/internal/domain/entities"
	"time"
)

type ReportsUC interface {
	// GetSalesReport returns a sales report for a given date range
	GetSalesReport(startDate time.Time, endDate time.Time) (*entities.SalesReport, error)
}