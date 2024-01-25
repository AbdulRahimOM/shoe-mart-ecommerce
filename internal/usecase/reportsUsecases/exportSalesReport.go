package reportsusecases

import (
	"MyShoo/internal/domain/entities"
	"MyShoo/internal/models/requestModels"
	repoInterface "MyShoo/internal/repository/interface"
	usecaseInterface "MyShoo/internal/usecase/interface"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/xuri/excelize/v2"
)

type ReportsUseCase struct {
	orderRepo   repoInterface.IOrderRepo
	reportsRepo repoInterface.IReportsRepo
}

func NewReportsUseCase(reportsRepo repoInterface.IReportsRepo, orderRepo repoInterface.IOrderRepo) usecaseInterface.ReportsUC {
	return &ReportsUseCase{
		reportsRepo: reportsRepo,
		orderRepo:   orderRepo,
	}
}

func (uc *ReportsUseCase) ExportSalesReportFullTime() (string, error) {

	var salesReport *[]entities.SalesReport
	var url string
	salesReport, err := uc.reportsRepo.GetSalesReportFullTime()
	if err != nil {
		fmt.Println("Error getting sales report:", err)
		return "", err
	}
	
	file := excelize.NewFile()
	_, err = file.NewSheet("Sales Report")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	file.DeleteSheet("Sheet1")

	// Set value of title cells
	titleCells := []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1", "I1", "J1", "K1"}
	titles := []string{"Date", "Order ID", "Seller ID", "Brand Name", "Model Name", "Product Name", "SKU", "Quantity", "MRP", "Sale Price", "Net Amount"}
	for i, title := range titles {
		file.SetCellValue("Sales Report", titleCells[i], title)
	}

	// Set values of cells
	for i, salesReport := range *salesReport {
		file.SetCellValue("Sales Report", fmt.Sprintf("A%d", i+2), salesReport.Date)
		file.SetCellValue("Sales Report", fmt.Sprintf("B%d", i+2), salesReport.OrderID)
		file.SetCellValue("Sales Report", fmt.Sprintf("C%d", i+2), salesReport.SellerID)
		file.SetCellValue("Sales Report", fmt.Sprintf("D%d", i+2), salesReport.BrandName)
		file.SetCellValue("Sales Report", fmt.Sprintf("E%d", i+2), salesReport.ModelName)
		file.SetCellValue("Sales Report", fmt.Sprintf("F%d", i+2), salesReport.ProductName)
		file.SetCellValue("Sales Report", fmt.Sprintf("G%d", i+2), salesReport.SKU)
		file.SetCellValue("Sales Report", fmt.Sprintf("H%d", i+2), salesReport.Quantity)
		file.SetCellValue("Sales Report", fmt.Sprintf("I%d", i+2), salesReport.MRP)
		file.SetCellValue("Sales Report", fmt.Sprintf("J%d", i+2), salesReport.SalePrice)
		file.SetCellFormula("Sales Report", fmt.Sprintf("K%d", i+2), fmt.Sprintf("H%d*J%d", i+2, i+2))
	}

	// Save the Excel file
	if err := file.SaveAs("output.xlsx"); err != nil {
		fmt.Println("Error saving Excel file:", err)
		return "", err
	}

	// Save the Excel file in the temporary directory
	tempFilePath := filepath.Join(os.TempDir(), "output.xlsx")
	err = file.SaveAs(tempFilePath)
	if err != nil {
		fmt.Println("Error saving Excel file:", err)
		return "", err
	}
	defer os.Remove(tempFilePath)

	var excelUploadReq requestModels.ExcelFileReq = requestModels.ExcelFileReq{
		File: tempFilePath,
		UploadParams: uploader.UploadParams{
			Folder:    "MyShoo/adminreports",
			PublicID:  "fulltimeReportForAdmin",
			Overwrite: true,
		},
	}

	url, err = uc.reportsRepo.UploadExcelFile(&excelUploadReq)
	if err != nil {
		fmt.Println("Error uploading Excel file:", err)
		return "", err
	}
	fmt.Println("url: ", url)
	return url, nil
}

func (*ReportsUseCase) ExportSalesReportLastMonth() error {
	panic("unimplemented")
}

// ExportSalesReportLastWeek implements usecaseInterface.ReportsUC.
func (*ReportsUseCase) ExportSalesReportLastWeek() error {
	panic("unimplemented")
}

// ExportSalesReportLastYear implements usecaseInterface.ReportsUC.
func (*ReportsUseCase) ExportSalesReportLastYear() error {
	panic("unimplemented")
}

// ExportSalesReportThisMonth implements usecaseInterface.ReportsUC.
func (*ReportsUseCase) ExportSalesReportThisMonth() error {
	panic("unimplemented")
}

// ExportSalesReportThisWeek implements usecaseInterface.ReportsUC.
func (*ReportsUseCase) ExportSalesReportThisWeek() error {
	panic("unimplemented")
}

// ExportSalesReportThisYear implements usecaseInterface.ReportsUC.
func (*ReportsUseCase) ExportSalesReportThisYear() error {
	panic("unimplemented")
}

// ExportSalesReportToday implements usecaseInterface.ReportsUC.
func (*ReportsUseCase) ExportSalesReportToday() error {
	panic("unimplemented")
}

// ExportSalesReportYesterday implements usecaseInterface.ReportsUC.
func (*ReportsUseCase) ExportSalesReportYesterday() error {
	panic("unimplemented")
}

// ExportSalesReportBetweenDates implements usecaseInterface.ReportsUC.
func (*ReportsUseCase) ExportSalesReportBetweenDates(startDate time.Time, endDate time.Time) error {
	panic("unimplemented")
}
