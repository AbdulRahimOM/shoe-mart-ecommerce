package reporthandlers

import (
	"MyShoo/internal/domain/entities"
	response "MyShoo/internal/models/responseModels"
	usecaseInterface "MyShoo/internal/usecase/interface"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportsHandler struct {
	reportsUC usecaseInterface.ReportsUC
}

func NewReportsHandler(reportsUseCase usecaseInterface.ReportsUC) *ReportsHandler {
	return &ReportsHandler{reportsUC: reportsUseCase}
}

// GetSalesReport returns a sales report for a given date range
func (h *ReportsHandler) GetSalesReport(c *gin.Context) {
	fmt.Println("Handler ::: get sales report handler")
	var salesReport *entities.SalesReport
	var salesPerDay *[]entities.SalePerDay
	var err error

	//get url param
	rangeType := c.Param("range")

	switch rangeType {
	case "full-time":
		salesReport, salesPerDay, err = h.reportsUC.GetSalesReportFullTime()
		if err != nil {
			fmt.Println("error occured in getting sales report")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.GetSalesReportResponse{
				Status:      "success",
				Message:     "Sales report - full time",
				Error:       "",
				SalesReport: *salesReport,
				SalePerDay:  *salesPerDay,
			})
		}
	case "date-range":
		//get query params from url
		startDate := c.Query("sd")
		endDate := c.Query("ed")
		if startDate == "" && endDate == "" {
			fmt.Println("No date range provided")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "No date range provided. Please recheck URL and try again",
				Error:   "No date range provided in url.",
			})
			return
		} else {
			// validate date params, and recieve time
			startTime, err := requestValidation.ValidateAndParseDate(startDate)
			if err != nil {
				fmt.Println("error occured in validating and parsing start date")
				c.JSON(http.StatusBadRequest, response.SME{
					Status:  "failed",
					Message: "Error occured. Please recheck URL and try again",
					Error:   err.Error(),
				})
				return
			}
			endTime, err := requestValidation.ValidateAndParseDate(endDate)
			if err != nil {
				fmt.Println("error occured in validating and parsing end date")
				c.JSON(400, response.SME{
					Status:  "failed",
					Message: "Error occured. Please recheck URL and try again",
					Error:   err.Error(),
				})
				return
			}

			//get sales report
			salesReport, salesPerDay, err = h.reportsUC.GetSalesReportBetweenDates(startTime, endTime)
			if err != nil {
				fmt.Println("error occured in getting sales report")
				c.JSON(400, response.SME{
					Status:  "failed",
					Message: "Error occured. Please recheck URL and try again", //need update
					Error:   err.Error(),
				})
				return
			} else {
				c.JSON(http.StatusOK, response.GetSalesReportResponse{
					Status:      "success",
					Message:     fmt.Sprint("Sales report between ", startDate, " and ", endDate),
					Error:       "",
					SalesReport: *salesReport,
					SalePerDay:  *salesPerDay,
				})
				return
			}

		}
	case "this-month":
		salesReport, salesPerDay, err = h.reportsUC.GetSalesReportThisMonth()
		if err != nil {
			fmt.Println("error occured in getting sales report for this month")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.GetSalesReportResponse{
				Status:      "success",
				Message:     "Sales report - this month",
				Error:       "",
				SalesReport: *salesReport,
				SalePerDay:  *salesPerDay,
			})
		}

	case "last-month":
		salesReport, salesPerDay, err = h.reportsUC.GetSalesReportLastMonth()
		if err != nil {
			fmt.Println("error occured in getting sales report for last month")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.GetSalesReportResponse{
				Status:      "success",
				Message:     "Sales report - last month",
				Error:       "",
				SalesReport: *salesReport,
				SalePerDay:  *salesPerDay,
			})
		}

	case "this-week":
		salesReport, salesPerDay, err = h.reportsUC.GetSalesReportThisWeek()
		if err != nil {
			fmt.Println("error occured in getting sales report for this week")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.GetSalesReportResponse{
				Status:      "success",
				Message:     "Sales report - this week",
				Error:       "",
				SalesReport: *salesReport,
				SalePerDay:  *salesPerDay,
			})
		}

	case "last-week":
		salesReport, salesPerDay, err = h.reportsUC.GetSalesReportLastWeek()
		if err != nil {
			fmt.Println("error occured in getting sales report for last week")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.GetSalesReportResponse{
				Status:      "success",
				Message:     "Sales report - last week",
				Error:       "",
				SalesReport: *salesReport,
				SalePerDay:  *salesPerDay,
			})
		}
	case "this-year":
		salesReport, salesPerDay, err = h.reportsUC.GetSalesReportThisYear()
		if err != nil {
			fmt.Println("error occured in getting sales report for this year")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.GetSalesReportResponse{
				Status:      "success",
				Message:     "Sales report - this year",
				Error:       "",
				SalesReport: *salesReport,
				SalePerDay:  *salesPerDay,
			})
		}

	case "last-year":
		salesReport, salesPerDay, err = h.reportsUC.GetSalesReportLastYear()
		if err != nil {
			fmt.Println("error occured in getting sales report for last year")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.GetSalesReportResponse{
				Status:      "success",
				Message:     "Sales report - last year",
				Error:       "",
				SalesReport: *salesReport,
				SalePerDay:  *salesPerDay,
			})
		}

	case "today":
		salesReport, salesPerDay, err = h.reportsUC.GetSalesReportToday()
		if err != nil {
			fmt.Println("error occured in getting sales report for today")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.GetSalesReportResponse{
				Status:      "success",
				Message:     "Sales report - today",
				Error:       "",
				SalesReport: *salesReport,
				SalePerDay:  *salesPerDay,
			})
		}
	case "yesterday":
		salesReport, salesPerDay, err = h.reportsUC.GetSalesReportYesterday()
		if err != nil {
			fmt.Println("error occured in getting sales report for yesterday")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.GetSalesReportResponse{
				Status:      "success",
				Message:     "Sales report - yesterday",
				Error:       "",
				SalesReport: *salesReport,
				SalePerDay:  *salesPerDay,
			})
		}

	default:
		c.JSON(400, response.SME{
			Status:  "failed",
			Message: "Please recheck URL and try again",
			Error:   "Invalid url query.",
		})
		return
	}
}
