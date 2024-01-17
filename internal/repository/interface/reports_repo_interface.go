package repository_interface

import (
	"MyShoo/internal/domain/entities"
	"time"
)

type IReportsRepo interface{
	GetSalesReportBetweenDates(start time.Time, end time.Time) (*entities.SalesReport,*[]entities.SalePerDay, error)
	GetSalesReportFullTime() (*entities.SalesReport,*[]entities.SalePerDay, error)
}