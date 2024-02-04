package usecaseInterface

import (
	"MyShoo/internal/domain/entities"
	response "MyShoo/internal/models/responseModels"
	"time"
)

type ReportsUC interface {
	// GetDashBoardData returns a dashboard data for a given date range
	GetDashBoardDataFullTime() (*entities.DashboardData, *[]entities.SalePerDay, error)
	GetDashBoardDataBetweenDates(startDate time.Time, endDate time.Time) (*entities.DashboardData, *[]entities.SalePerDay, error)
	GetDashBoardDataLastYear() (*entities.DashboardData, *[]entities.SalePerDay, error)
	GetDashBoardDataThisYear() (*entities.DashboardData, *[]entities.SalePerDay, error)
	GetDashBoardDataLastMonth() (*entities.DashboardData, *[]entities.SalePerDay, error)
	GetDashBoardDataThisMonth() (*entities.DashboardData, *[]entities.SalePerDay, error)
	GetDashBoardDataLastWeek() (*entities.DashboardData, *[]entities.SalePerDay, error)
	GetDashBoardDataThisWeek() (*entities.DashboardData, *[]entities.SalePerDay, error)
	GetDashBoardDataYesterday() (*entities.DashboardData, *[]entities.SalePerDay, error)
	GetDashBoardDataToday() (*entities.DashboardData, *[]entities.SalePerDay, error)

	ExportSalesReportFullTime() (string, error)
	ExportSalesReportBetweenDates(startDate time.Time, endDate time.Time) (string, error)
	ExportSalesReportThisMonth() (string, error)
	ExportSalesReportLastMonth() (string, error)
	ExportSalesReportThisYear() (string, error)
	ExportSalesReportLastYear() (string, error)
	ExportSalesReportThisWeek() (string, error)
	ExportSalesReportLastWeek() (string, error)
	ExportSalesReportToday() (string, error)
	ExportSalesReportYesterday() (string, error)

	GetTopProducts(limit int) (*[]response.TopProducts, error)
	GetTopSellers(limit int) (*[]response.TopSellers, error)
	GetTopBrands(limit int) (*[]response.TopBrands, error)
	GetTopModels(limit int) (*[]response.TopModels, error)
}
