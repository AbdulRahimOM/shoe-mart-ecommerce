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

//export sales report handler
// @Summary Export sales report
// @Description Export sales report
// @Tags Admin/Analytics
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param range path string true "Range type"
// @Param sd query string false "Start Date"
// @Param ed query string false "End Date"
// @Success 200 {file} application/octet-stream
// @Failure 400 {object} response.SME{} "Error details"
// @Router /admin/exportsalesreport/{range} [get]

func (h *ReportsHandler) ExportSalesReport(c *gin.Context) {
	rangeType := c.Param("range")

	switch rangeType {
	case "full-time":
		fileURL, err := h.reportsUC.ExportSalesReportFullTime()
		if err != nil {
			c.JSON(500, response.FailedSME("An Error occured!", err))
			return
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportFullTimeForAdmin.xlsx")
			c.File(fileURL)
		}
	case "date-range":
		startDate := c.Query("sd")
		endDate := c.Query("ed")
		if startDate == "" && endDate == "" {
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
				c.JSON(500, response.FailedSME("An Error occured!", err))
				return
			} else {
				c.Header("Content-Disposition", "attachment; filename=SalesReportBwCustomDatesForAdmin.xlsx")
				c.File(fileURL)
			}
		}
	case "this-month":
		fileURL, err := h.reportsUC.ExportSalesReportThisMonth()
		if err != nil {
			c.JSON(500, response.FailedSME("An Error occured!", err))
			return
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportThisMonthForAdmin.xlsx")
			c.File(fileURL)
		}
	case "last-month":
		fmt.Println("last month entered")
		fileURL, err := h.reportsUC.ExportSalesReportLastMonth()
		if err != nil {
			c.JSON(500, response.FailedSME("An Error occured!", err))
			return
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportLastMonth.xlsx")
			c.File(fileURL)
		}

	case "this-year":
		fileURL, err := h.reportsUC.ExportSalesReportThisYear()
		if err != nil {
			c.JSON(500, response.FailedSME("An Error occured!", err))
			return
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportThisYearForAdmin.xlsx")
			c.File(fileURL)
		}

	case "last-year":
		fileURL, err := h.reportsUC.ExportSalesReportLastYear()
		if err != nil {
			c.JSON(500, response.FailedSME("An Error occured!", err))
			return
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportLastYearForAdmin.xlsx")
			c.File(fileURL)
		}

	case "this-week":
		fileURL, err := h.reportsUC.ExportSalesReportThisWeek()
		if err != nil {
			c.JSON(500, response.FailedSME("An Error occured!", err))
			return
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportThisWeekForAdmin.xlsx")
			c.File(fileURL)
		}

	case "last-week":
		fileURL, err := h.reportsUC.ExportSalesReportLastWeek()
		if err != nil {
			c.JSON(500, response.FailedSME("An Error occured!", err))
			return
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportLastWeekForAdmin.xlsx")
			c.File(fileURL)
		}

	case "today":
		fileURL, err := h.reportsUC.ExportSalesReportToday()
		if err != nil {
			c.JSON(500, response.FailedSME("An Error occured!", err))
			return
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportTodayForAdmin.xlsx")
			c.File(fileURL)
		}

	case "yesterday":
		fileURL, err := h.reportsUC.ExportSalesReportYesterday()
		if err != nil {
			c.JSON(500, response.FailedSME("An Error occured!", err))
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
