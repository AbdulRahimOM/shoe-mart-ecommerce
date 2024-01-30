package reporthandlers

import (
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

//export dashBoardData handler
func (h *ReportsHandler) ExportSalesReport(c *gin.Context) {
	fmt.Println("Handler ::: export sales report handler")

	//get url param
	rangeType := c.Param("range")
	fmt.Println("range type: ", rangeType)

	//print url
	fmt.Println("url: ", c.Request.URL)

	switch rangeType {
	case "full-time":
		fileURL, err := h.reportsUC.ExportSalesReportFullTime()
		if err != nil {
			fmt.Println("error occured in exporting sales report for full time")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "An Error occured!",
				Error:   err.Error(),
			})
			return
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportFullTimeForAdmin.xlsx")
			c.File(fileURL)
		}
	case "date-range":
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
					Message: "Invalid date provided. Please recheck URL and try again",
					Error:   err.Error(),
				})
				return
			}
			endTime, err := requestValidation.ValidateAndParseDate(endDate)
			if err != nil {
				fmt.Println("error occured in validating and parsing end date")
				c.JSON(http.StatusBadRequest, response.SME{
					Status:  "failed",
					Message: "Invalid date provided. Please recheck URL and try again",
					Error:   err.Error(),
				})
				return
			}

			fileURL, err := h.reportsUC.ExportSalesReportBetweenDates(startTime, endTime)
			if err != nil {
				fmt.Println("error occured in exporting sales report between dates")
				c.JSON(400, response.SME{
					Status:  "failed",
					Message: "An Error occured!",
					Error:   err.Error(),
				})
				return
			} else {
				c.Header("Content-Disposition", "attachment; filename=SalesReportBwCustomDatesForAdmin.xlsx")
				c.File(fileURL)
			}
		}
	case "this-month":
		fileURL, err := h.reportsUC.ExportSalesReportThisMonth()
		if err != nil {
			fmt.Println("error occured in exporting sales report for this month")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "An Error occured!",
				Error:   err.Error(),
			})
			return
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportThisMonthForAdmin.xlsx")
			c.File(fileURL)
		}
	case "last-month":
		fmt.Println("last month entered")
		fileURL, err := h.reportsUC.ExportSalesReportLastMonth()
		if err != nil {
			fmt.Println("error occured in exporting sales report for last month")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "An Error occured!",
				Error:   err.Error(),
			})
			return
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportLastMonth.xlsx")
			c.File(fileURL)
		}

	case "this-year":
		fileURL, err := h.reportsUC.ExportSalesReportThisYear()
		if err != nil {
			fmt.Println("error occured in exporting sales report for this year")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "An Error occured!",
				Error:   err.Error(),
			})
			return
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportThisYearForAdmin.xlsx")
			c.File(fileURL)
		}

	case "last-year":
		fileURL, err := h.reportsUC.ExportSalesReportLastYear()
		if err != nil {
			fmt.Println("error occured in exporting sales report for last year")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "An Error occured!",
				Error:   err.Error(),
			})
			return
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportLastYearForAdmin.xlsx")
			c.File(fileURL)
		}

	case "this-week":
		fileURL, err := h.reportsUC.ExportSalesReportThisWeek()
		if err != nil {
			fmt.Println("error occured in exporting sales report for this week")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "An Error occured!",
				Error:   err.Error(),
			})
			return
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportThisWeekForAdmin.xlsx")
			c.File(fileURL)
		}

	case "last-week":
		fileURL, err := h.reportsUC.ExportSalesReportLastWeek()
		if err != nil {
			fmt.Println("error occured in exporting sales report for last week")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "An Error occured!",
				Error:   err.Error(),
			})
			return
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportLastWeekForAdmin.xlsx")
			c.File(fileURL)
		}

	case "today":
		fileURL, err := h.reportsUC.ExportSalesReportToday()
		if err != nil {
			fmt.Println("error occured in exporting sales report for today")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "An Error occured!",
				Error:   err.Error(),
			})
			return
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportTodayForAdmin.xlsx")
			c.File(fileURL)
		}

	case "yesterday":
		fileURL, err := h.reportsUC.ExportSalesReportYesterday()
		if err != nil {
			fmt.Println("error occured in exporting sales report for yesterday")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "An Error occured!",
				Error:   err.Error(),
			})
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportYesterdayForAdmin.xlsx")
			c.File(fileURL)
		}

	default:
		fmt.Println("Invalid range type provided")
		c.JSON(400, response.SME{
			Status:  "failed",
			Message: "Invalid range type provided. Please recheck URL and try again",
			Error:   "Invalid range type provided in url.",
		})
		return
	}
}
