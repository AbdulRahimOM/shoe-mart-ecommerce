package usecase

import (
	"time"

	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	response "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/responseModels"
)

type ReportsUC interface {
	// GetDashBoardData returns a dashboard data for a given date range
	GetDashBoardDataFullTime() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error)
	GetDashBoardDataBetweenDates(startDate time.Time, endDate time.Time) (*entities.DashboardData, *[]entities.SalePerDay, *e.Error)
	GetDashBoardDataLastYear() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error)
	GetDashBoardDataThisYear() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error)
	GetDashBoardDataLastMonth() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error)
	GetDashBoardDataThisMonth() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error)
	GetDashBoardDataLastWeek() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error)
	GetDashBoardDataThisWeek() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error)
	GetDashBoardDataYesterday() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error)
	GetDashBoardDataToday() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error)

	ExportSalesReportFullTime() (*string, *e.Error)
	ExportSalesReportBetweenDates(startDate time.Time, endDate time.Time) (*string, *e.Error)
	ExportSalesReportThisMonth() (*string, *e.Error)
	ExportSalesReportLastMonth() (*string, *e.Error)
	ExportSalesReportThisYear() (*string, *e.Error)
	ExportSalesReportLastYear() (*string, *e.Error)
	ExportSalesReportThisWeek() (*string, *e.Error)
	ExportSalesReportLastWeek() (*string, *e.Error)
	ExportSalesReportToday() (*string, *e.Error)
	ExportSalesReportYesterday() (*string, *e.Error)

	GetTopProducts(limit int) (*[]response.TopProducts, *e.Error)
	GetTopSellers(limit int) (*[]response.TopSellers, *e.Error)
	GetTopBrands(limit int) (*[]response.TopBrands, *e.Error)
	GetTopModels(limit int) (*[]response.TopModels, *e.Error)
}
