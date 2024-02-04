package repository_interface

import (
	"MyShoo/internal/domain/entities"
	"MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	"time"
)

type IReportsRepo interface {
	GetDashBoardDataBetweenDates(start time.Time, end time.Time) (*entities.DashboardData, *[]entities.SalePerDay, error)
	GetDashBoardDataFullTime() (*entities.DashboardData, *[]entities.SalePerDay, error)

	UploadExcelFile(req *requestModels.ExcelFileReq) (string, error)

	GetSalesReportFullTime() (
		*[]entities.SalesReportOrderList,
		*[]entities.SellerWiseReport,
		*[]entities.BrandWiseReport,
		*[]entities.ModelWiseReport,
		*[]entities.SizeWiseReport,
		*[]entities.RevenueGraph,
		error)

	GetSalesReportBetweenDates(startDate time.Time, endDate time.Time) (
		*[]entities.SalesReportOrderList,
		*[]entities.SellerWiseReport,
		*[]entities.BrandWiseReport,
		*[]entities.ModelWiseReport,
		*[]entities.SizeWiseReport,
		*[]entities.RevenueGraph,
		error)

	GetTopModels(limit int) (*[]response.TopModels, error)
	GetTopProducts(limit int) (*[]response.TopProducts, error)
	GetTopBrands(limit int) (*[]response.TopBrands, error)
	GetTopSellers(limit int) (*[]response.TopSellers, error)
}
