package usecaseInterface

import (
	"MyShoo/internal/domain/entities"
	"time"
)

type ReportsUC interface {
	// GetSalesReport returns a sales report for a given date range
	GetSalesReportBetweenDates(startDate time.Time, endDate time.Time) (*entities.SalesReport,*[]entities.SalePerDay, error)
	GetSalesReportLastYear() (*entities.SalesReport,*[]entities.SalePerDay, error)
	GetSalesReportThisYear() (*entities.SalesReport,*[]entities.SalePerDay, error)
	GetSalesReportLastMonth() (*entities.SalesReport,*[]entities.SalePerDay, error)
	GetSalesReportThisMonth() (*entities.SalesReport,*[]entities.SalePerDay, error)
	GetSalesReportLastWeek() (*entities.SalesReport,*[]entities.SalePerDay, error)
	GetSalesReportThisWeek() (*entities.SalesReport,*[]entities.SalePerDay, error)
	GetSalesReportYesterday() (*entities.SalesReport,*[]entities.SalePerDay, error)
	GetSalesReportToday() (*entities.SalesReport,*[]entities.SalePerDay, error)

	GetSalesReportFullTime() (*entities.SalesReport,*[]entities.SalePerDay, error)
}