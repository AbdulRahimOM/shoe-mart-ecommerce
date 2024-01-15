package reporthandlers

import (
	// "MyShoo/internal/domain/entities"
	// requestModels "MyShoo/internal/models/requestModels"
	response "MyShoo/internal/models/responseModels"
	"net/http"
	// "MyShoo/internal/tools"
	usecaseInterface "MyShoo/internal/usecase/interface"
	requestValidation "MyShoo/pkg/validation"
	"fmt"

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

	//get query params from url
	startDate := c.Query("sd")
	endDate := c.Query("ed")

	//validate request
	if startDate == "" && endDate == "" {
		//allow
	} else 	if startDate == "" || endDate == "" {
		c.JSON(400, response.SME{
			Status:  "failed",
			Message: "Error validating request. Try Again",
			Error:   "start date or end date is empty",
		})
		return
	}else {
		//validate date params, and recieve time
		startTime,err:=requestValidation.ValidateAndParseDate(startDate)
		if err != nil {
			fmt.Println("error occured in validating and parsing start date")
			c.JSON(http.StatusBadRequest,response.SME{
				Status: "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error: err.Error(),
			})
		}
		endTime,err:=requestValidation.ValidateAndParseDate(endDate)
		if err != nil {
			fmt.Println("error occured in validating and parsing end date")
			c.JSON(400,response.SME{
				Status: "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error: err.Error(),
			})
		}

		//get sales report
		salesReport,err:=h.reportsUC.GetSalesReport(startTime,endTime)
		if err != nil {
			fmt.Println("error occured in getting sales report")
			c.JSON(400,response.SME{
				Status: "failed",
				Message: "Error occured. Please recheck URL and try again",
				Error: err.Error(),
			})
		}

		c.JSON(http.StatusOK,response.GetSalesReportResponse{
			Status: "success",
			Message: fmt.Sprint("Sales report between ",startDate," and ",endDate),
			Error: "",
			SalesReport: *salesReport,
		})
		
	}


}
