package reportsrepo

import (
	"MyShoo/internal/domain/entities"
	"MyShoo/internal/models/requestModels"
	"MyShoo/internal/services"
	"fmt"

	repoInterface "MyShoo/internal/repository/interface"

	"github.com/cloudinary/cloudinary-go"
	"gorm.io/gorm"
)

type DashboardDataRepo struct {
	DB  *gorm.DB
	Cld *cloudinary.Cloudinary
}

func NewReportRepository(db *gorm.DB, cloudinary *cloudinary.Cloudinary) repoInterface.IReportsRepo {
	return &DashboardDataRepo{
		DB:  db,
		Cld: cloudinary,
	}
}

// UploadExcelFile
func (repo *DashboardDataRepo) UploadExcelFile(req *requestModels.ExcelFileReq) (string, error) {
	fileUploadService := services.NewFileUploadService(repo.Cld)

	var err error
	url, err := fileUploadService.UploadExcelFile(req)
	if err != nil {
		return "", err
	}
	fmt.Println("url=", url)

	return url, nil
}

// GetSalesReportFullTime
func (repo *DashboardDataRepo) GetSalesReportFullTime() (*[]entities.SalesReport, error) {
	var salesReport []entities.SalesReport
	query := repo.DB.Raw(`
		SELECT
			"orders"."id" AS "order_id",
			"orders"."order_date_and_time" AS "date",
			"brands"."sellerId" AS "seller_id",
			"brands"."name" AS "brand_name",
			"models"."name" AS "model_name",
			"product"."name" AS "product_name",
			"product"."skuCode" AS "sku",
			"order_items"."quantity" AS "quantity",
			"colour_variants"."mrp" AS "mrp",
			"colour_variants"."salePrice" AS "sale_price"
		FROM "orders"
		INNER JOIN "order_items" ON "orders"."id" = "order_items"."order_id"
		INNER JOIN "product" ON "order_items"."product_id" = "product"."id"
		INNER JOIN "dimensional_variants" ON "product"."dimensionalVariationID" = "dimensional_variants"."id"
		INNER JOIN "colour_variants" ON "dimensional_variants"."colourVariantId"= "colour_variants"."id"
		INNER JOIN "models" ON "colour_variants"."modelId" = "models"."id"
		INNER JOIN "brands" ON "models"."brandId" = "brands"."id"
		WHERE "orders"."order_date_and_time" >= '2019-01-01 00:00:00'
		AND "orders"."order_date_and_time" <= '2024-12-31 23:59:59'
		ORDER BY "orders"."order_date_and_time" DESC
	`).Scan(&salesReport)
	


	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return nil, query.Error
	}
	// fmt.Println("salesReport=", salesReport)

	return &salesReport, nil
}
