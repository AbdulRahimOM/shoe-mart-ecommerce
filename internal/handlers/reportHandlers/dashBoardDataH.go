package reporthandlers

import (
	"MyShoo/internal/domain/entities"
	response "MyShoo/internal/models/responseModels"
	requestValidation "MyShoo/pkg/validation"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDashBoardData returns a dashboard data for a given date range
func (h *ReportsHandler) GetDashBoardData(c *gin.Context) {
	fmt.Println("Handler ::: get dashboard data handler")
	var dashBoardData *entities.DashboardData
	var salesPerDay *[]entities.SalePerDay
	var err error

	//get url param
	rangeType := c.Param("range")

	switch rangeType {
	case "full-time":
		dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataFullTime()
		if err != nil {
			fmt.Println("error occured in getting dashboard data")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
				Status:        "success",
				Message:       "DashBoard Data - full time",
				Error:         "",
				DashboardData: *dashBoardData,
				SalePerDay:    *salesPerDay,
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

			//get dashboard data
			dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataBetweenDates(startTime, endTime)
			if err != nil {
				fmt.Println("error occured in getting dashboard data")
				c.JSON(400, response.SME{
					Status:  "failed",
					Message: "Error occured. Please recheck URL and try again", //need update
					Error:   err.Error(),
				})
				return
			} else {
				c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
					Status:        "success",
					Message:       fmt.Sprint("DashBoard Data between ", startDate, " and ", endDate),
					Error:         "",
					DashboardData: *dashBoardData,
					SalePerDay:    *salesPerDay,
				})
				return
			}

		}
	case "this-month":
		dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataThisMonth()
		if err != nil {
			fmt.Println("error occured in getting dashboard data for this month")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
				Status:        "success",
				Message:       "DashBoard Data - this month",
				Error:         "",
				DashboardData: *dashBoardData,
				SalePerDay:    *salesPerDay,
			})
		}

	case "last-month":
		dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataLastMonth()
		if err != nil {
			fmt.Println("error occured in getting dashboard data for last month")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
				Status:        "success",
				Message:       "DashBoard Data - last month",
				Error:         "",
				DashboardData: *dashBoardData,
				SalePerDay:    *salesPerDay,
			})
		}

	case "this-week":
		dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataThisWeek()
		if err != nil {
			fmt.Println("error occured in getting dashboard data for this week")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
				Status:        "success",
				Message:       "DashBoard Data - this week",
				Error:         "",
				DashboardData: *dashBoardData,
				SalePerDay:    *salesPerDay,
			})
		}

	case "last-week":
		dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataLastWeek()
		if err != nil {
			fmt.Println("error occured in getting dashboard data for last week")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
				Status:        "success",
				Message:       "DashBoard Data - last week",
				Error:         "",
				DashboardData: *dashBoardData,
				SalePerDay:    *salesPerDay,
			})
		}
	case "this-year":
		dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataThisYear()
		if err != nil {
			fmt.Println("error occured in getting dashboard data for this year")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
				Status:        "success",
				Message:       "DashBoard Data - this year",
				Error:         "",
				DashboardData: *dashBoardData,
				SalePerDay:    *salesPerDay,
			})
		}

	case "last-year":
		dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataLastYear()
		if err != nil {
			fmt.Println("error occured in getting dashboard data for last year")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
				Status:        "success",
				Message:       "DashBoard Data - last year",
				Error:         "",
				DashboardData: *dashBoardData,
				SalePerDay:    *salesPerDay,
			})
		}

	case "today":
		dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataToday()
		if err != nil {
			fmt.Println("error occured in getting dashboard data for today")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
				Status:        "success",
				Message:       "DashBoard Data - today",
				Error:         "",
				DashboardData: *dashBoardData,
				SalePerDay:    *salesPerDay,
			})
		}
	case "yesterday":
		dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataYesterday()
		if err != nil {
			fmt.Println("error occured in getting dashboard data for yesterday")
			c.JSON(400, response.SME{
				Status:  "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error:   err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
				Status:        "success",
				Message:       "DashBoard Data - yesterday",
				Error:         "",
				DashboardData: *dashBoardData,
				SalePerDay:    *salesPerDay,
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
