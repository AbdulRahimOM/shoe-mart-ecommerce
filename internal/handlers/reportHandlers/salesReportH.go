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

	switch rangeType {
	case "full-time":
		fileURL,err := h.reportsUC.ExportSalesReportFullTime()
		if err != nil {
			fmt.Println("error occured in exporting sales report for full time")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
		// Setting Content-Disposition header to make the file downloadable
		c.Header("Content-Disposition", "attachment; filename=SalesReportFullTimeForAdmin")
		// Redirect to the file URL
		c.Redirect(http.StatusFound, fileURL)
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

			err = h.reportsUC.ExportSalesReportBetweenDates(startTime, endTime)
			if err != nil {
				fmt.Println("error occured in exporting sales report between dates")
				c.JSON(400, response.SME{
					Status:  "failed",
					Message: "Error occured. Please recheck URL and try again",
					Error:   err.Error(),
				})
				return
			} else {
				c.JSON(http.StatusOK, response.SME{
					Status:  "success",
					Message: "Sales report exported successfully",
					Error:   "",
				})
			}
		}
	case "this-month":
		err := h.reportsUC.ExportSalesReportThisMonth()
		if err != nil {
			fmt.Println("error occured in exporting sales report for this month")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {	
			c.JSON(http.StatusOK, response.SME{
				Status:  "success",
				Message: "Sales report exported successfully",
				Error:   "",
			})
		}
	case "last-month":
		err := h.reportsUC.ExportSalesReportLastMonth()
		if err != nil {
			fmt.Println("error occured in exporting sales report for last month")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.SME{
				Status:  "success",
				Message: "Sales report exported successfully",
				Error:   "",
			})
		}

	case "this-year":
		err := h.reportsUC.ExportSalesReportThisYear()
		if err != nil {
			fmt.Println("error occured in exporting sales report for this year")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.SME{
				Status:  "success",
				Message: "Sales report exported successfully",
				Error:   "",
			})
		}

	case "last-year":
		err := h.reportsUC.ExportSalesReportLastYear()
		if err != nil {
			fmt.Println("error occured in exporting sales report for last year")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.SME{
				Status:  "success",
				Message: "Sales report exported successfully",
				Error:   "",
			})
		}

	case "this-week":
		err := h.reportsUC.ExportSalesReportThisWeek()
		if err != nil {
			fmt.Println("error occured in exporting sales report for this week")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.SME{
				Status:  "success",
				Message: "Sales report exported successfully",
				Error:   "",
			})	
		}

	case "last-week":
		err := h.reportsUC.ExportSalesReportLastWeek()
		if err != nil {
			fmt.Println("error occured in exporting sales report for last week")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.SME{
				Status:  "success",
				Message: "Sales report exported successfully",
				Error:   "",
			})
		}

	case "today":
		err := h.reportsUC.ExportSalesReportToday()
		if err != nil {
			fmt.Println("error occured in exporting sales report for today")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.SME{
				Status:  "success",
				Message: "Sales report exported successfully",
				Error:   "",
			})
		}

	case "yesterday":
		err := h.reportsUC.ExportSalesReportYesterday()
		if err != nil {
			fmt.Println("error occured in exporting sales report for yesterday")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, response.SME{
				Status:  "success",
				Message: "Sales report exported successfully",
				Error:   "",
			})
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
