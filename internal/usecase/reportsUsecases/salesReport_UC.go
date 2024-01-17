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

func (uc *ReportsUseCase) GetSalesReportBetweenDates(start time.Time, end time.Time) (*entities.SalesReport, *[]entities.SalePerDay, error) {
	salesReport, salePerDay, err := uc.reportsRepo.GetSalesReportBetweenDates(start, end)
	if err != nil {
		return salesReport, salePerDay, err
	}
	return salesReport, salePerDay, nil
}

func (uc *ReportsUseCase) GetSalesReportFullTime() (*entities.SalesReport, *[]entities.SalePerDay, error) {
	salesReport, salePerDay, err := uc.reportsRepo.GetSalesReportFullTime()
	if err != nil {
		return salesReport, salePerDay, err
	}
	return salesReport, salePerDay, nil
}

func (uc *ReportsUseCase) GetSalesReportLastMonth() (*entities.SalesReport, *[]entities.SalePerDay, error) {
	var start, end time.Time
	if time.Now().Month() == 1 {
		start = time.Date(time.Now().Year()-1, 12, 1, 0, 0, 0, 0, time.Now().Location())
		end = time.Date(time.Now().Year()-1, 12, 31, 0, 0, 0, 0, time.Now().Location())
	} else {
		start = time.Date(time.Now().Year(), time.Now().Month()-1, 1, 0, 0, 0, 0, time.Now().Location())
		end = time.Date(time.Now().Year(), time.Now().Month(),1, 0, 0, 0, 0, time.Now().Location())
	}

	salesReport, salePerDay, err := uc.reportsRepo.GetSalesReportBetweenDates(start, end)
	if err != nil {
		return salesReport, salePerDay, err
	}
	return salesReport, salePerDay, nil
}

func (uc *ReportsUseCase) GetSalesReportThisMonth() (*entities.SalesReport, *[]entities.SalePerDay, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	salesReport, salePerDay, err := uc.reportsRepo.GetSalesReportBetweenDates(start, now)
	if err != nil {
		return salesReport, salePerDay, err
	}
	return salesReport, salePerDay, nil
}

func (uc *ReportsUseCase) GetSalesReportLastWeek() (*entities.SalesReport, *[]entities.SalePerDay, error) {
	lastWeekSundayThisTime:= time.Now().AddDate(0, 0, -int(time.Now().Weekday()-7))	//need this to prevent negative day
	start := time.Date(lastWeekSundayThisTime.Year(), lastWeekSundayThisTime.Month(), lastWeekSundayThisTime.Day(), 0, 0, 0, 0, time.Now().Location())
	end := start.AddDate(0, 0, 7)

	salesReport, salePerDay, err := uc.reportsRepo.GetSalesReportBetweenDates(start, end)
	if err != nil {
		return salesReport, salePerDay, err
	}
	return salesReport, salePerDay, nil
}

func (uc *ReportsUseCase) GetSalesReportThisWeek() (*entities.SalesReport, *[]entities.SalePerDay, error) {
	now := time.Now()
	thisWeekSundayThisTime:= now.AddDate(0, 0, -int(now.Weekday()))	//need this to prevent negative day
	start := time.Date(thisWeekSundayThisTime.Year(), thisWeekSundayThisTime.Month(), thisWeekSundayThisTime.Day(), 0, 0, 0, 0, now.Location())

	salesReport, salePerDay, err := uc.reportsRepo.GetSalesReportBetweenDates(start, now)
	if err != nil {
		return salesReport, salePerDay, err
	}
	return salesReport, salePerDay, nil
}

func (uc *ReportsUseCase) GetSalesReportLastYear() (*entities.SalesReport, *[]entities.SalePerDay, error) {
	now:=time.Now()
	start := time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
	end := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	salesReport, salePerDay, err := uc.reportsRepo.GetSalesReportBetweenDates(start, end)
	if err != nil {
		return salesReport, salePerDay, err
	}
	return salesReport, salePerDay, nil
}

func (uc *ReportsUseCase) GetSalesReportThisYear() (*entities.SalesReport, *[]entities.SalePerDay, error) {
	now := time.Now()
	start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	salesReport, salePerDay, err := uc.reportsRepo.GetSalesReportBetweenDates(start, now)
	if err != nil {
		return salesReport, salePerDay, err
	}
	return salesReport, salePerDay, nil
}

func (uc *ReportsUseCase) GetSalesReportToday() (*entities.SalesReport, *[]entities.SalePerDay, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	salesReport, salePerDay, err := uc.reportsRepo.GetSalesReportBetweenDates(start, now)
	if err != nil {
		return salesReport, salePerDay, err
	}
	return salesReport, salePerDay, nil
}

func (uc *ReportsUseCase) GetSalesReportYesterday() (*entities.SalesReport, *[]entities.SalePerDay, error) {
	now:=time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	start := end.AddDate(0, 0, -1)

	salesReport, salePerDay, err := uc.reportsRepo.GetSalesReportBetweenDates(start, end)
	if err != nil {
		return salesReport, salePerDay, err
	}
	return salesReport, salePerDay, nil
}
