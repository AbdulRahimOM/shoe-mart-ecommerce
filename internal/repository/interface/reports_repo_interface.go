package repository_interface

import (
	"MyShoo/internal/domain/entities"
	"time"
)

type IReportsRepo interface{
	GetSalesReport(start time.Time, end time.Time) (*entities.SalesReport, error)
}