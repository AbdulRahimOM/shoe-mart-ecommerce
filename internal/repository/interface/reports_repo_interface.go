package repo

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	response "MyShoo/internal/models/responseModels"
	"time"
)

type IReportsRepo interface {
	GetDashBoardDataBetweenDates(start time.Time, end time.Time) (*entities.DashboardData, *[]entities.SalePerDay, *e.Error)
	GetDashBoardDataFullTime() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error)

	UploadSalesReportExcel(filePath string, rangeLabel string) (string, *e.Error)

	GetSalesReportFullTime() (
		*[]entities.SalesReportOrderList,
		*[]entities.SellerWiseReport,
		*[]entities.BrandWiseReport,
		*[]entities.ModelWiseReport,
		*[]entities.SizeWiseReport,
		*[]entities.RevenueGraph,
		*e.Error)

	GetSalesReportBetweenDates(startDate time.Time, endDate time.Time) (
		*[]entities.SalesReportOrderList,
		*[]entities.SellerWiseReport,
		*[]entities.BrandWiseReport,
		*[]entities.ModelWiseReport,
		*[]entities.SizeWiseReport,
		*[]entities.RevenueGraph,
		*e.Error)

	GetTopModels(limit int) (*[]response.TopModels, *e.Error)
	GetTopProducts(limit int) (*[]response.TopProducts, *e.Error)
	GetTopBrands(limit int) (*[]response.TopBrands, *e.Error)
	GetTopSellers(limit int) (*[]response.TopSellers, *e.Error)
}
