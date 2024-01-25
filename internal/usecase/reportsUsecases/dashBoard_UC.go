package reportsusecases

import (
	"MyShoo/internal/domain/entities"
	"time"
	// "github.com/jinzhu/copier"
)


func (uc *ReportsUseCase) GetDashBoardDataBetweenDates(start time.Time, end time.Time) (*entities.DashboardData, *[]entities.SalePerDay, error) {
	dashBoardData, salePerDay, err := uc.reportsRepo.GetDashBoardDataBetweenDates(start, end)
	if err != nil {
		return dashBoardData, salePerDay, err
	}
	return dashBoardData, salePerDay, nil
}

func (uc *ReportsUseCase) GetDashBoardDataFullTime() (*entities.DashboardData, *[]entities.SalePerDay, error) {
	dashBoardData, salePerDay, err := uc.reportsRepo.GetDashBoardDataFullTime()
	if err != nil {
		return dashBoardData, salePerDay, err
	}
	return dashBoardData, salePerDay, nil
}

func (uc *ReportsUseCase) GetDashBoardDataLastMonth() (*entities.DashboardData, *[]entities.SalePerDay, error) {
	var start, end time.Time
	if time.Now().Month() == 1 {
		start = time.Date(time.Now().Year()-1, 12, 1, 0, 0, 0, 0, time.Now().Location())
		end = time.Date(time.Now().Year()-1, 12, 31, 0, 0, 0, 0, time.Now().Location())
	} else {
		start = time.Date(time.Now().Year(), time.Now().Month()-1, 1, 0, 0, 0, 0, time.Now().Location())
		end = time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	}

	dashBoardData, salePerDay, err := uc.reportsRepo.GetDashBoardDataBetweenDates(start, end)
	if err != nil {
		return dashBoardData, salePerDay, err
	}
	return dashBoardData, salePerDay, nil
}

func (uc *ReportsUseCase) GetDashBoardDataThisMonth() (*entities.DashboardData, *[]entities.SalePerDay, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	dashBoardData, salePerDay, err := uc.reportsRepo.GetDashBoardDataBetweenDates(start, now)
	if err != nil {
		return dashBoardData, salePerDay, err
	}
	return dashBoardData, salePerDay, nil
}

func (uc *ReportsUseCase) GetDashBoardDataLastWeek() (*entities.DashboardData, *[]entities.SalePerDay, error) {
	lastWeekSundayThisTime := time.Now().AddDate(0, 0, -int(time.Now().Weekday()-7)) //need this to prevent negative day
	start := time.Date(lastWeekSundayThisTime.Year(), lastWeekSundayThisTime.Month(), lastWeekSundayThisTime.Day(), 0, 0, 0, 0, time.Now().Location())
	end := start.AddDate(0, 0, 7)

	dashBoardData, salePerDay, err := uc.reportsRepo.GetDashBoardDataBetweenDates(start, end)
	if err != nil {
		return dashBoardData, salePerDay, err
	}
	return dashBoardData, salePerDay, nil
}

func (uc *ReportsUseCase) GetDashBoardDataThisWeek() (*entities.DashboardData, *[]entities.SalePerDay, error) {
	now := time.Now()
	thisWeekSundayThisTime := now.AddDate(0, 0, -int(now.Weekday())) //need this to prevent negative day
	start := time.Date(thisWeekSundayThisTime.Year(), thisWeekSundayThisTime.Month(), thisWeekSundayThisTime.Day(), 0, 0, 0, 0, now.Location())

	dashBoardData, salePerDay, err := uc.reportsRepo.GetDashBoardDataBetweenDates(start, now)
	if err != nil {
		return dashBoardData, salePerDay, err
	}
	return dashBoardData, salePerDay, nil
}

func (uc *ReportsUseCase) GetDashBoardDataLastYear() (*entities.DashboardData, *[]entities.SalePerDay, error) {
	now := time.Now()
	start := time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
	end := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	dashBoardData, salePerDay, err := uc.reportsRepo.GetDashBoardDataBetweenDates(start, end)
	if err != nil {
		return dashBoardData, salePerDay, err
	}
	return dashBoardData, salePerDay, nil
}

func (uc *ReportsUseCase) GetDashBoardDataThisYear() (*entities.DashboardData, *[]entities.SalePerDay, error) {
	now := time.Now()
	start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	dashBoardData, salePerDay, err := uc.reportsRepo.GetDashBoardDataBetweenDates(start, now)
	if err != nil {
		return dashBoardData, salePerDay, err
	}
	return dashBoardData, salePerDay, nil
}

func (uc *ReportsUseCase) GetDashBoardDataToday() (*entities.DashboardData, *[]entities.SalePerDay, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	dashBoardData, salePerDay, err := uc.reportsRepo.GetDashBoardDataBetweenDates(start, now)
	if err != nil {
		return dashBoardData, salePerDay, err
	}
	return dashBoardData, salePerDay, nil
}

func (uc *ReportsUseCase) GetDashBoardDataYesterday() (*entities.DashboardData, *[]entities.SalePerDay, error) {
	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	start := end.AddDate(0, 0, -1)

	dashBoardData, salePerDay, err := uc.reportsRepo.GetDashBoardDataBetweenDates(start, end)
	if err != nil {
		return dashBoardData, salePerDay, err
	}
	return dashBoardData, salePerDay, nil
}
