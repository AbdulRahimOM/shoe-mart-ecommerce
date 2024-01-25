package response

import "MyShoo/internal/domain/entities"

type GetDashBoardDataResponse struct {
	Status        string                 `json:"status"`
	Message       string                 `json:"message"`
	Error         string                 `json:"error"`
	DashboardData entities.DashboardData `json:"sales_report"` //need update
	SalePerDay    []entities.SalePerDay  `json:"sale_per_day"`
}
