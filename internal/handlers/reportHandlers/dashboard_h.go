package reporthandlers

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	response "MyShoo/internal/models/responseModels"
	requestValidation "MyShoo/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDashBoardData returns a dashboard data for a given date range
// @Summary Get dashboard data
// @Description Get dashboard data (for a given date range)
// @Tags Admin/Analytics
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param range path string true "Range"
// @Param sd query string false "Start Date"
// @Param ed query string false "End Date"
// @Success 200 {object} response.GetDashBoardDataResponse{}
// @Failure 400 {object} response.SME{}
// @Router /admin/dashboarddata/{range} [get]
func (h *ReportsHandler) GetDashBoardData(c *gin.Context) {

	var dashBoardData *entities.DashboardData
	var salesPerDay *[]entities.SalePerDay
	var err *e.Error

	//get url param
	rangeType := c.Param("range")

	switch rangeType {
	case "full-time":
		dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataFullTime()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
			return
		} else {
			c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
				DashboardData: *dashBoardData,
				SalePerDay:    *salesPerDay,
			})
		}
	case "date-range":
		//get query params from url
		startDate := c.Query("sd")
		endDate := c.Query("ed")
		if startDate == "" && endDate == "" {
			c.JSON(http.StatusBadRequest, response.FromErrByText("no date range provided inURL"))
			return
		} else {
			// validate date params, and recieve time
			startTime, errr := requestValidation.ValidateAndParseDate(startDate)
			if errr != nil {
				c.JSON(http.StatusBadRequest, response.MsgAndError("error occured in validating and parsing start date", errr))
				return
			}
			endTime, errr := requestValidation.ValidateAndParseDate(endDate)
			if errr != nil {
				c.JSON(http.StatusBadRequest, response.MsgAndError("error occured in validating and parsing end date", errr))
				return
			}

			//get dashboard data
			dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataBetweenDates(startTime, endTime)
			if err != nil {
				c.JSON(err.StatusCode, response.FromError(err))
				return
			} else {
				c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
					DashboardData: *dashBoardData,
					SalePerDay:    *salesPerDay,
				})
				return
			}

		}
	case "this-month":
		dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataThisMonth()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
			return
		} else {
			c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
				DashboardData: *dashBoardData,
				SalePerDay:    *salesPerDay,
			})
			return
		}

	case "last-month":
		dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataLastMonth()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
			return
		} else {
			c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
				DashboardData: *dashBoardData,
				SalePerDay:    *salesPerDay,
			})
			return
		}

	case "this-week":
		dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataThisWeek()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
			return
		} else {
			c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
				DashboardData: *dashBoardData,
				SalePerDay:    *salesPerDay,
			})
			return
		}

	case "last-week":
		dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataLastWeek()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
			return
		} else {
			c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
				DashboardData: *dashBoardData,
				SalePerDay:    *salesPerDay,
			})
			return
		}
	case "this-year":
		dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataThisYear()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
			return
		} else {
			c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
				DashboardData: *dashBoardData,
				SalePerDay:    *salesPerDay,
			})
			return
		}

	case "last-year":
		dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataLastYear()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
			return
		} else {
			c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
				DashboardData: *dashBoardData,
				SalePerDay:    *salesPerDay,
			})
			return
		}

	case "today":
		dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataToday()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
			return
		} else {
			c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
				DashboardData: *dashBoardData,
				SalePerDay:    *salesPerDay,
			})
			return
		}
	case "yesterday":
		dashBoardData, salesPerDay, err = h.reportsUC.GetDashBoardDataYesterday()
		if err != nil {
			c.JSON(err.StatusCode, response.FromError(err))
			return
		} else {
			c.JSON(http.StatusOK, response.GetDashBoardDataResponse{
				DashboardData: *dashBoardData,
				SalePerDay:    *salesPerDay,
			})
			return
		}

	default:
		c.JSON(http.StatusBadRequest, response.FromErrByText("invalid url query"))
		return
	}
}
