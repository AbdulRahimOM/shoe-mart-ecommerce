package reportsusecases

import (
	"MyShoo/internal/config"
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	repoInterface "MyShoo/internal/repository/interface"
	usecase "MyShoo/internal/usecase/interface"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
)

type ReportsUseCase struct {
	orderRepo   repoInterface.IOrderRepo
	reportsRepo repoInterface.IReportsRepo
}

func NewReportsUseCase(reportsRepo repoInterface.IReportsRepo, orderRepo repoInterface.IOrderRepo) usecase.ReportsUC {
	return &ReportsUseCase{
		reportsRepo: reportsRepo,
		orderRepo:   orderRepo,
	}
}

func (uc *ReportsUseCase) ExportSalesReportFullTime() (*string, *e.Error) {
	return uc.processAdminExcelReport(uc.reportsRepo.GetSalesReportFullTime())
}

func (uc *ReportsUseCase) ExportSalesReportLastMonth() (*string, *e.Error) {
	var start, end time.Time
	if time.Now().Month() == 1 {
		start = time.Date(time.Now().Year()-1, 12, 1, 0, 0, 0, 0, time.Now().Location())
	} else {
		start = time.Date(time.Now().Year(), time.Now().Month()-1, 1, 0, 0, 0, 0, time.Now().Location())
	}
	end = start.AddDate(0, 1, 0)
	return uc.processAdminExcelReport(uc.reportsRepo.GetSalesReportBetweenDates(start, end))
}

func (uc *ReportsUseCase) ExportSalesReportLastWeek() (*string, *e.Error) {
	lastWeekSundayThisTime := time.Now().AddDate(0, 0, -int(time.Now().Weekday()-7)) //need this to prevent negative day
	start := time.Date(lastWeekSundayThisTime.Year(), lastWeekSundayThisTime.Month(), lastWeekSundayThisTime.Day(), 0, 0, 0, 0, time.Now().Location())
	end := start.AddDate(0, 0, 7)
	return uc.processAdminExcelReport(uc.reportsRepo.GetSalesReportBetweenDates(start, end))
}
func (uc *ReportsUseCase) ExportSalesReportLastYear() (*string, *e.Error) {
	now := time.Now()
	start := time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(1, 0, 0)
	return uc.processAdminExcelReport(uc.reportsRepo.GetSalesReportBetweenDates(start, end))
}

func (uc *ReportsUseCase) ExportSalesReportThisMonth() (*string, *e.Error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return uc.processAdminExcelReport(uc.reportsRepo.GetSalesReportBetweenDates(start, now))
}

func (uc *ReportsUseCase) ExportSalesReportThisWeek() (*string, *e.Error) {
	now := time.Now()
	thisWeekSundayThisTime := now.AddDate(0, 0, -int(now.Weekday())) //need this to prevent negative day
	start := time.Date(thisWeekSundayThisTime.Year(), thisWeekSundayThisTime.Month(), thisWeekSundayThisTime.Day(), 0, 0, 0, 0, now.Location())
	return uc.processAdminExcelReport(uc.reportsRepo.GetSalesReportBetweenDates(start, now))
}

func (uc *ReportsUseCase) ExportSalesReportThisYear() (*string, *e.Error) {
	now := time.Now()
	start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	return uc.processAdminExcelReport(uc.reportsRepo.GetSalesReportBetweenDates(start, now))
}

func (uc *ReportsUseCase) ExportSalesReportToday() (*string, *e.Error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return uc.processAdminExcelReport(uc.reportsRepo.GetSalesReportBetweenDates(start, now))
}

func (uc *ReportsUseCase) ExportSalesReportYesterday() (*string, *e.Error) {
	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	start := end.AddDate(0, 0, -1)
	return uc.processAdminExcelReport(uc.reportsRepo.GetSalesReportBetweenDates(start, end))
}

func (uc *ReportsUseCase) ExportSalesReportBetweenDates(startDate time.Time, endDate time.Time) (*string, *e.Error) {
	return uc.processAdminExcelReport(uc.reportsRepo.GetSalesReportBetweenDates(startDate, endDate))
}

func (uc *ReportsUseCase) processAdminExcelReport(
	orderList *[]entities.SalesReportOrderList,
	sellerWiseReport *[]entities.SellerWiseReport,
	brandWiseReport *[]entities.BrandWiseReport,
	modelWiseReport *[]entities.ModelWiseReport,
	sizeWiseReport *[]entities.SizeWiseReport,
	revenueGraph *[]entities.RevenueGraph,
	err *e.Error,
) (*string, *e.Error) {
	if err != nil {
		return nil, err
	}
	filePath := filepath.Join(config.ExecutableDir, "internal/templates/TemplateExcel.xlsx")
	file, errr := excelize.OpenFile(filePath)
	if errr != nil {
		return nil, e.TextCumError("Error opening template:", errr, 500)
	}

	_, errr = file.NewSheet("All Orders")
	if errr != nil {
		return nil, e.TextCumError("Error creating new sheet:", errr, 500)
	}

	// Set value of title cells
	titleCells := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1", "I1", "J1", "K1"}
	titles := []string{"Date", "Order ID", "Seller ID", "Brand Name", "Model Name", "Product Name", "SKU", "Quantity", "MRP", "Sale Price", "Net Amount"}
	for i, title := range titles {
		file.SetCellValue("All Orders", titleCells[i], title)
	}

	// Set values of cells
	for i, order := range *orderList {
		file.SetCellValue("All Orders", fmt.Sprintf("A%d", i+2), order.Date)
		file.SetCellValue("All Orders", fmt.Sprintf("B%d", i+2), order.OrderID)
		file.SetCellValue("All Orders", fmt.Sprintf("C%d", i+2), order.SellerID)
		file.SetCellValue("All Orders", fmt.Sprintf("D%d", i+2), order.BrandName)
		file.SetCellValue("All Orders", fmt.Sprintf("E%d", i+2), order.ModelName)
		file.SetCellValue("All Orders", fmt.Sprintf("F%d", i+2), order.ProductName)
		file.SetCellValue("All Orders", fmt.Sprintf("G%d", i+2), order.SKU)
		file.SetCellValue("All Orders", fmt.Sprintf("H%d", i+2), order.Quantity)
		file.SetCellValue("All Orders", fmt.Sprintf("I%d", i+2), order.MRP)
		file.SetCellValue("All Orders", fmt.Sprintf("J%d", i+2), order.SalePrice)
		file.SetCellFormula("All Orders", fmt.Sprintf("K%d", i+2), fmt.Sprintf("H%d*J%d", i+2, i+2))
	}

	for i, sellerData := range *sellerWiseReport {
		file.SetCellValue("Seller Data", fmt.Sprintf("A%d", i+2), sellerData.SellerID)
		file.SetCellValue("Seller Data", fmt.Sprintf("B%d", i+2), sellerData.SellerName)
		file.SetCellValue("Seller Data", fmt.Sprintf("C%d", i+2), sellerData.QuantityCount)
		file.SetCellValue("Seller Data", fmt.Sprintf("D%d", i+2), sellerData.TotalSale)
		file.SetCellValue("Seller Data", fmt.Sprintf("E%d", i+2), sellerData.ReturnQuantityCount)
		file.SetCellValue("Seller Data", fmt.Sprintf("F%d", i+2), sellerData.TotalReturnValue)
		file.SetCellValue("Seller Data", fmt.Sprintf("G%d", i+2), sellerData.CancelQuantityCount)
		file.SetCellValue("Seller Data", fmt.Sprintf("H%d", i+2), sellerData.TotalCancelValue)
		file.SetCellValue("Seller Data", fmt.Sprintf("I%d", i+2), sellerData.EffectiveQuantityCount)
		file.SetCellValue("Seller Data", fmt.Sprintf("J%d", i+2), sellerData.EffectiveSaleValue)
		file.SetCellFormula("Seller Data", fmt.Sprintf("K%d", i+2), fmt.Sprintf("ROUND(E%d/C%d%%,1)", i+2, i+2))
	}

	for i, brandData := range *brandWiseReport {
		file.SetCellValue("Brand Data", fmt.Sprintf("A%d", i+2), brandData.BrandID)
		file.SetCellValue("Brand Data", fmt.Sprintf("B%d", i+2), brandData.BrandName)
		file.SetCellValue("Brand Data", fmt.Sprintf("C%d", i+2), brandData.QuantityCount)
		file.SetCellValue("Brand Data", fmt.Sprintf("D%d", i+2), brandData.TotalSale)
		file.SetCellValue("Brand Data", fmt.Sprintf("E%d", i+2), brandData.ReturnQuantityCount)
		file.SetCellValue("Brand Data", fmt.Sprintf("F%d", i+2), brandData.TotalReturnValue)
		file.SetCellValue("Brand Data", fmt.Sprintf("G%d", i+2), brandData.CancelQuantityCount)
		file.SetCellValue("Brand Data", fmt.Sprintf("H%d", i+2), brandData.TotalCancelValue)
		file.SetCellValue("Brand Data", fmt.Sprintf("I%d", i+2), brandData.EffectiveQuantityCount)
		file.SetCellValue("Brand Data", fmt.Sprintf("J%d", i+2), brandData.EffectiveSaleValue)
		file.SetCellFormula("Seller Data", fmt.Sprintf("K%d", i+2), fmt.Sprintf("ROUND(E%d/C%d%%,1)", i+2, i+2))
	}

	for i, modelData := range *modelWiseReport {
		file.SetCellValue("Model Data", fmt.Sprintf("A%d", i+2), modelData.ModelID)
		file.SetCellValue("Model Data", fmt.Sprintf("B%d", i+2), modelData.ModelName)
		file.SetCellValue("Model Data", fmt.Sprintf("C%d", i+2), modelData.BrandName)
		file.SetCellValue("Model Data", fmt.Sprintf("D%d", i+2), modelData.QuantityCount)
		file.SetCellValue("Model Data", fmt.Sprintf("E%d", i+2), modelData.TotalSale)
		file.SetCellValue("Model Data", fmt.Sprintf("F%d", i+2), modelData.ReturnQuantityCount)
		file.SetCellValue("Model Data", fmt.Sprintf("G%d", i+2), modelData.CancelQuantityCount)
		file.SetCellValue("Model Data", fmt.Sprintf("H%d", i+2), modelData.EffectiveQuantityCount)
		file.SetCellValue("Model Data", fmt.Sprintf("I%d", i+2), modelData.EffectiveSaleValue)
		file.SetCellFormula("Seller Data", fmt.Sprintf("J%d", i+2), fmt.Sprintf("ROUND(F%d/D%d%%,1)", i+2, i+2))
	}

	for _, sizeData := range *sizeWiseReport {
		file.SetCellValue("Size Data", fmt.Sprintf("B%d", sizeData.SizeIndex+3), sizeData.QuantityCount)
	}

	if err := file.AddChart("Size Data", "C1", &excelize.Chart{
		Type: excelize.Col,
		Series: func() []excelize.ChartSeries {
			var series []excelize.ChartSeries
			for i := 0; i < 19; i++ {
				series = append(series, excelize.ChartSeries{
					Name:       fmt.Sprintf("'Size Data'!$A$%d", i+2),
					Categories: "'Size Data'!$B$1:$B$1",
					Values:     fmt.Sprintf("'Size Data'!$B$%d:$B$%d", i+2, i+2),
				})
			}
			return series
		}(),
		Title: []excelize.RichTextRun{
			{
				Text: "Size Data",
			},
		},
	}); err != nil {
		return nil, e.TextCumError("error occured while creating chart:", err, 500)
	}

	for i, revenueData := range *revenueGraph {
		file.SetCellValue("Sale Graph", fmt.Sprintf("A%d", i+3), revenueData.Date)
		file.SetCellValue("Sale Graph", fmt.Sprintf("B%d", i+3), revenueData.Sale)
		file.SetCellValue("Sale Graph", fmt.Sprintf("C%d", i+3), revenueData.Quantity)
	}

	// Save the Excel file in the temporary directory
	tempFilePath := filepath.Join(os.TempDir(), "output.xlsx")
	errr = file.SaveAs(tempFilePath)
	defer os.Remove(tempFilePath)
	if errr != nil {
		return nil, e.TextCumError("error saving excel file:", errr, 500)
	}

	if config.ShouldUploadExcel {
		return uc.reportsRepo.UploadSalesReportExcel(tempFilePath, "myrange")
	} else {
		// Saving the Excel file locally (for dev/testing purposes)
		localUrl := filepath.Join(config.ExecutableDir, "testKit/salesReportOutput.xlsx")
		if errr = file.SaveAs(localUrl); errr != nil {
			return nil, e.TextCumError("error saving excel file:", errr, 500)
		} else {
			return &localUrl, nil
		}
	}
}
