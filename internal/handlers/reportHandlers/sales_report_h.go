package reporthandlers

import (
	response "MyShoo/internal/models/responseModels"
	usecase "MyShoo/internal/usecase/interface"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportsHandler struct {
	reportsUC usecase.ReportsUC
}

func NewReportsHandler(reportsUseCase usecase.ReportsUC) *ReportsHandler {
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
			c.JSON(err.StatusCode, response.FromError(err))
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportFullTimeForAdmin.xlsx")
			c.File(*fileURL)
		}
	case "date-range":
		startDate := c.Query("sd")
		endDate := c.Query("ed")
		if startDate == "" || endDate == "" {
			c.JSON(http.StatusBadRequest, response.FromErrByText("start and end date ranges not provided inURL"))
		} else {
			// validate date params, and recieve time
			startTime, errr := requestValidation.ValidateAndParseDate(startDate)
			if errr != nil {
				c.JSON(http.StatusBadRequest, response.FromErrByTextCumError("error occured in validating and parsing start date", errr))
				return
			}
			endTime, errr := requestValidation.ValidateAndParseDate(endDate)
			if errr != nil {
				c.JSON(http.StatusBadRequest, response.FromErrByTextCumError("error occured in validating and parsing end date", errr))
			}

			fileURL, err := h.reportsUC.ExportSalesReportBetweenDates(startTime, endTime)
			if err != nil {
				c.JSON(err.StatusCode, response.FromError(err))
				return
			} else {
				c.Header("Content-Disposition", "attachment; filename=SalesReportBwCustomDatesForAdmin.xlsx")
				c.File(*fileURL)
			}
		}
	case "this-month":
		fileURL, err := h.reportsUC.ExportSalesReportThisMonth()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportThisMonthForAdmin.xlsx")
			c.File(*fileURL)
		}
		
	case "last-month":
		fmt.Println("last month entered")
		fileURL, err := h.reportsUC.ExportSalesReportLastMonth()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportLastMonth.xlsx")
			c.File(*fileURL)
		}

	case "this-year":
		fileURL, err := h.reportsUC.ExportSalesReportThisYear()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportThisYearForAdmin.xlsx")
			c.File(*fileURL)
		}

	case "last-year":
		fileURL, err := h.reportsUC.ExportSalesReportLastYear()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportLastYearForAdmin.xlsx")
			c.File(*fileURL)
		}

	case "this-week":
		fileURL, err := h.reportsUC.ExportSalesReportThisWeek()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportThisWeekForAdmin.xlsx")
			c.File(*fileURL)
		}

	case "last-week":
		fileURL, err := h.reportsUC.ExportSalesReportLastWeek()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportLastWeekForAdmin.xlsx")
			c.File(*fileURL)
		}

	case "today":
		fileURL, err := h.reportsUC.ExportSalesReportToday()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportTodayForAdmin.xlsx")
			c.File(*fileURL)
		}

	case "yesterday":
		fileURL, err := h.reportsUC.ExportSalesReportYesterday()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
		} else {
			c.Header("Content-Disposition", "attachment; filename=SalesReportYesterdayForAdmin.xlsx")
			c.File(*fileURL)
		}

	default:
		c.JSON(http.StatusBadRequest, response.FromErrByText("invalid range type provided inURL"))
	}
}
