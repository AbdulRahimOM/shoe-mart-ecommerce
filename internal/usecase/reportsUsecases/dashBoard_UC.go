package reportsusecases

import (
	"time"

	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
)

func (uc *ReportsUseCase) GetDashBoardDataBetweenDates(start time.Time, end time.Time) (*entities.DashboardData, *[]entities.SalePerDay, *e.Error) {
	return uc.reportsRepo.GetDashBoardDataBetweenDates(start, end)
}

func (uc *ReportsUseCase) GetDashBoardDataFullTime() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error) {
	return uc.reportsRepo.GetDashBoardDataFullTime()
}

func (uc *ReportsUseCase) GetDashBoardDataLastMonth() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error) {
	var start, end time.Time
	if time.Now().Month() == 1 {
		start = time.Date(time.Now().Year()-1, 12, 1, 0, 0, 0, 0, time.Now().Location())
		end = time.Date(time.Now().Year()-1, 12, 31, 0, 0, 0, 0, time.Now().Location())
	} else {
		start = time.Date(time.Now().Year(), time.Now().Month()-1, 1, 0, 0, 0, 0, time.Now().Location())
		end = time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	}
	return uc.reportsRepo.GetDashBoardDataBetweenDates(start, end)
}

func (uc *ReportsUseCase) GetDashBoardDataThisMonth() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return uc.reportsRepo.GetDashBoardDataBetweenDates(start, now)
}

func (uc *ReportsUseCase) GetDashBoardDataLastWeek() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error) {
	lastWeekSundayThisTime := time.Now().AddDate(0, 0, -int(time.Now().Weekday()-7)) //need this to prevent negative day
	start := time.Date(lastWeekSundayThisTime.Year(), lastWeekSundayThisTime.Month(), lastWeekSundayThisTime.Day(), 0, 0, 0, 0, time.Now().Location())
	end := start.AddDate(0, 0, 7)
	return uc.reportsRepo.GetDashBoardDataBetweenDates(start, end)
}

func (uc *ReportsUseCase) GetDashBoardDataThisWeek() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error) {
	now := time.Now()
	thisWeekSundayThisTime := now.AddDate(0, 0, -int(now.Weekday())) //need this to prevent negative day
	start := time.Date(thisWeekSundayThisTime.Year(), thisWeekSundayThisTime.Month(), thisWeekSundayThisTime.Day(), 0, 0, 0, 0, now.Location())
	return uc.reportsRepo.GetDashBoardDataBetweenDates(start, now)
}

func (uc *ReportsUseCase) GetDashBoardDataLastYear() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error) {
	now := time.Now()
	start := time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
	end := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	return uc.reportsRepo.GetDashBoardDataBetweenDates(start, end)
}

func (uc *ReportsUseCase) GetDashBoardDataThisYear() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error) {
	now := time.Now()
	start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	return uc.reportsRepo.GetDashBoardDataBetweenDates(start, now)
}

func (uc *ReportsUseCase) GetDashBoardDataToday() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return uc.reportsRepo.GetDashBoardDataBetweenDates(start, now)
}

func (uc *ReportsUseCase) GetDashBoardDataYesterday() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error) {
	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	start := end.AddDate(0, 0, -1)
	return uc.reportsRepo.GetDashBoardDataBetweenDates(start, end)
}
