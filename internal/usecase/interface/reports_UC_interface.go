package usecaseInterface

import (
	"MyShoo/internal/domain/entities"
	"time"
)

type ReportsUC interface {
	// GetDashBoardData returns a dashboard data for a given date range
	GetDashBoardDataFullTime()(*entities.DashboardData, *[]entities.SalePerDay, error) 
	GetDashBoardDataBetweenDates(startDate time.Time, endDate time.Time) (*entities.DashboardData, *[]entities.SalePerDay, error)
	GetDashBoardDataLastYear() (*entities.DashboardData, *[]entities.SalePerDay, error)
	GetDashBoardDataThisYear() (*entities.DashboardData, *[]entities.SalePerDay, error)
	GetDashBoardDataLastMonth() (*entities.DashboardData, *[]entities.SalePerDay, error)
	GetDashBoardDataThisMonth() (*entities.DashboardData, *[]entities.SalePerDay, error)
	GetDashBoardDataLastWeek() (*entities.DashboardData, *[]entities.SalePerDay, error)
	GetDashBoardDataThisWeek() (*entities.DashboardData, *[]entities.SalePerDay, error)
	GetDashBoardDataYesterday() (*entities.DashboardData, *[]entities.SalePerDay, error)
	GetDashBoardDataToday() (*entities.DashboardData, *[]entities.SalePerDay, error)

	ExportSalesReportFullTime() (string,error)
	ExportSalesReportBetweenDates(startDate time.Time, endDate time.Time) error
	ExportSalesReportThisMonth() error
	ExportSalesReportLastMonth() error
	ExportSalesReportThisYear() error
	ExportSalesReportLastYear() error
	ExportSalesReportThisWeek() error
	ExportSalesReportLastWeek() error
	ExportSalesReportToday() error
	ExportSalesReportYesterday() error
}
