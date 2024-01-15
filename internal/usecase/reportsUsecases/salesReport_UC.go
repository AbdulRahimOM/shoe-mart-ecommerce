package reportsusecases

import (
	"MyShoo/internal/domain/entities"
	repoInterface "MyShoo/internal/repository/interface"
	usecaseInterface "MyShoo/internal/usecase/interface"
	"time"
	// "github.com/jinzhu/copier"
)

type ReportsUseCase struct {
	reportsRepo repoInterface.IReportsRepo
}

func NewReportsUseCase(reportsRepo repoInterface.IReportsRepo) usecaseInterface.ReportsUC {
	return &ReportsUseCase{
		reportsRepo: reportsRepo,
	}
}

// GetSalesReport returns a sales report for a given date range
func (uc *ReportsUseCase) GetSalesReport(start time.Time, end time.Time) (*entities.SalesReport, error) {
	var salesReport *entities.SalesReport
	salesReport, err := uc.reportsRepo.GetSalesReport(start, end)
	if err != nil {
		return salesReport,err
	}
	return salesReport,nil
}
