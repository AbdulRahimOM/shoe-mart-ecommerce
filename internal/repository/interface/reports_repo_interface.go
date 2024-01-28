package repository_interface

import (
	"MyShoo/internal/domain/entities"
	"MyShoo/internal/models/requestModels"
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
}
