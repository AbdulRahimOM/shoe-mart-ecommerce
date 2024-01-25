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

	var orderList *[]entities.SalesReportOrderList
	var sellerWiseReport *[]entities.SellerWiseReport
	var url string
	orderList, sellerWiseReport, err := uc.reportsRepo.GetSalesReportFullTime()
	if err != nil {
		fmt.Println("Error getting sales report:", err)
		return "", err
	}
	filePath := "output.xlsx"
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Println("Error opening template:", err)
	}
	// file := excelize.NewFile()
	_, err = file.NewSheet("All Orders")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// file.DeleteSheet("Sheet1")

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
		file.SetCellValue("Seller Data", fmt.Sprintf("A%d", i+3), sellerData.SellerID)
		file.SetCellValue("Seller Data", fmt.Sprintf("B%d", i+3), sellerData.SellerName)
		file.SetCellValue("Seller Data", fmt.Sprintf("C%d", i+3), sellerData.QuantityCount)
		file.SetCellValue("Seller Data", fmt.Sprintf("D%d", i+3), sellerData.TotalSale)
		file.SetCellValue("Seller Data", fmt.Sprintf("E%d", i+3), sellerData.ReturnQuantityCount)
		file.SetCellValue("Seller Data", fmt.Sprintf("F%d", i+3), sellerData.TotalReturnValue)
		file.SetCellValue("Seller Data", fmt.Sprintf("G%d", i+3), sellerData.CancelQuantityCount)
		file.SetCellValue("Seller Data", fmt.Sprintf("H%d", i+3), sellerData.TotalCancelValue)
		file.SetCellValue("Seller Data", fmt.Sprintf("I%d", i+3), sellerData.EffectiveQuantityCount)
		file.SetCellValue("Seller Data", fmt.Sprintf("J%d", i+3), sellerData.EffectiveSaleValue)
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
	if os.Getenv("UPLOAD_EXCEL") == "true" {
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
	}
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
