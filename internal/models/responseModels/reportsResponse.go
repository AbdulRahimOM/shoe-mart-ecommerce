package response

import "MyShoo/internal/domain/entities"

type GetSalesReportResponse struct {
	Status      string               `json:"status"`
	Message     string               `json:"message"`
	Error       string               `json:"error"`
	SalesReport entities.SalesReport `json:"sales_report"` //need update
	SalePerDay  []entities.SalePerDay `json:"sale_per_day"`
}
