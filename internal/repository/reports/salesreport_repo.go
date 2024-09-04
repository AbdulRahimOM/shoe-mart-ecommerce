package reportsrepo

import (
	"context"
	"errors"
	"time"

	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"

	repoInterface "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/repository/interface"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
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

func (repo *DashboardDataRepo) UploadSalesReportExcel(filePath string, rangeLabel string) (*string, *e.Error) {
	uploadParams := uploader.UploadParams{
		Folder:    "github.com/AbdulRahimOM/shoe-mart-ecommerce/adminreports",
		PublicID:  rangeLabel + "ReportForAdmin",
		Overwrite: true,
	}
	result, err := repo.Cld.Upload.Upload(context.Background(), filePath, uploadParams)
	if err != nil {
		return nil, e.SetError("error while uploading file to cloudinary", err, 500)
	}

	if result.Error.Message != "" {
		return nil, e.SetError("error while uploading file to cloudinary", errors.New(result.Error.Message), 500)
	}

	return &result.SecureURL, nil
}

func (repo *DashboardDataRepo) GetSalesReportFullTime() (
	*[]entities.SalesReportOrderList,
	*[]entities.SellerWiseReport,
	*[]entities.BrandWiseReport,
	*[]entities.ModelWiseReport,
	*[]entities.SizeWiseReport,
	*[]entities.RevenueGraph,
	*e.Error) {
	var orderList []entities.SalesReportOrderList
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
		ORDER BY "orders"."order_date_and_time" DESC
	`).Scan(&orderList)

	if query.Error != nil {
		return nil, nil, nil, nil, nil, nil, e.DBQueryError_500(&query.Error)
	}

	var sellerWiseReport []entities.SellerWiseReport
	query = repo.DB.Raw(`
			SELECT
				"brands"."sellerId" AS "seller_id",
				"sellers"."firstName" AS "seller_name",
				SUM ("order_items"."quantity") AS "quantity_count",
				SUM ("order_items"."quantity" * "colour_variants"."salePrice") AS "total_sale",
				SUM (CASE WHEN "orders"."status" IN ('returned','return requested') THEN "order_items"."quantity" ELSE 0 END) AS "return_quantity_count",
				SUM (CASE WHEN "orders"."status" IN ('returned','return requested') THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "total_return_value",
				SUM (CASE WHEN "orders"."status" = 'cancelled' THEN "order_items"."quantity" ELSE 0 END) AS "cancel_quantity_count",
				SUM (CASE WHEN "orders"."status" = 'cancelled' THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "total_cancel_value",
				SUM (CASE WHEN "orders"."status" NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" ELSE 0 END) AS "effective_quantity_count",
				SUM (CASE WHEN "orders"."status" NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "effective_sale_value"
			FROM "order_items"
			INNER JOIN "orders" ON "order_items"."order_id" = "orders"."id"
			INNER JOIN "product" ON "order_items"."product_id" = "product"."id"
			INNER JOIN "dimensional_variants" ON "product"."dimensionalVariationID" = "dimensional_variants"."id"
			INNER JOIN "colour_variants" ON "dimensional_variants"."colourVariantId"= "colour_variants"."id"
			INNER JOIN "models" ON "colour_variants"."modelId" = "models"."id"
			INNER JOIN "brands" ON "models"."brandId" = "brands"."id"
			INNER JOIN "sellers" ON "brands"."sellerId" = "sellers"."id"
			GROUP BY "brands"."sellerId", "sellers"."firstName"
			ORDER BY "brands"."sellerId" ASC
		`).Scan(&sellerWiseReport)

	if query.Error != nil {
		return nil, nil, nil, nil, nil, nil, e.DBQueryError_500(&query.Error)
	}

	var brandWiseReport []entities.BrandWiseReport
	query = repo.DB.Raw(`
			SELECT
				"brands"."id" AS "brand_id",
				"brands"."name" AS "brand_name",
				SUM ("order_items"."quantity") AS "quantity_count",
				SUM ("order_items"."quantity" * "colour_variants"."salePrice") AS "total_sale",
				SUM (CASE WHEN "orders"."status" IN ('returned','return requested') THEN "order_items"."quantity" ELSE 0 END) AS "return_quantity_count",
				SUM (CASE WHEN "orders"."status" IN ('returned','return requested') THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "total_return_value",
				SUM (CASE WHEN "orders"."status" = 'cancelled' THEN "order_items"."quantity" ELSE 0 END) AS "cancel_quantity_count",
				SUM (CASE WHEN "orders"."status" = 'cancelled' THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "total_cancel_value",
				SUM (CASE WHEN "orders"."status"  NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" ELSE 0 END) AS "effective_quantity_count",
				SUM (CASE WHEN "orders"."status" NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "effective_sale_value"
			FROM "order_items"
			INNER JOIN "orders" ON "order_items"."order_id" = "orders"."id"
			INNER JOIN "product" ON "order_items"."product_id" = "product"."id"
			INNER JOIN "dimensional_variants" ON "product"."dimensionalVariationID" = "dimensional_variants"."id"
			INNER JOIN "colour_variants" ON "dimensional_variants"."colourVariantId"= "colour_variants"."id"
			INNER JOIN "models" ON "colour_variants"."modelId" = "models"."id"
			INNER JOIN "brands" ON "models"."brandId" = "brands"."id"
			GROUP BY "brands"."id"
			ORDER BY "brands"."id" ASC
		`).Scan(&brandWiseReport)

	if query.Error != nil {
		return nil, nil, nil, nil, nil, nil, e.DBQueryError_500(&query.Error)
	}

	var modelWiseReport []entities.ModelWiseReport
	query = repo.DB.Raw(`
			SELECT
				"models"."id" AS "model_id",
				"models"."name" AS "model_name",
				"brands"."name" AS "brand_name",
				SUM ("order_items"."quantity") AS "quantity_count",
				SUM ("order_items"."quantity" * "colour_variants"."salePrice") AS "total_sale",
				SUM (CASE WHEN "orders"."status" IN ('returned','return requested') THEN "order_items"."quantity" ELSE 0 END) AS "return_quantity_count",
				SUM (CASE WHEN "orders"."status" = 'cancelled' THEN "order_items"."quantity" ELSE 0 END) AS "cancel_quantity_count",
				SUM (CASE WHEN "orders"."status"  NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" ELSE 0 END) AS "effective_quantity_count",
				SUM (CASE WHEN "orders"."status" NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "effective_sale_value"
			FROM "order_items"
			INNER JOIN "orders" ON "order_items"."order_id" = "orders"."id"
			INNER JOIN "product" ON "order_items"."product_id" = "product"."id"
			INNER JOIN "dimensional_variants" ON "product"."dimensionalVariationID" = "dimensional_variants"."id"
			INNER JOIN "colour_variants" ON "dimensional_variants"."colourVariantId"= "colour_variants"."id"
			INNER JOIN "models" ON "colour_variants"."modelId" = "models"."id"
			INNER JOIN "brands" ON "models"."brandId" = "brands"."id"
			GROUP BY "models"."id","brands"."name"
			ORDER BY "models"."id" ASC

		`).Scan(&modelWiseReport)

	if query.Error != nil {
		return nil, nil, nil, nil, nil, nil, e.DBQueryError_500(&query.Error)
	}

	var sizeWiseReport []entities.SizeWiseReport
	query = repo.DB.Raw(`
			SELECT
				"product"."sizeIndex" AS "size_index",
				SUM ("order_items"."quantity") AS "quantity_count"
			FROM "order_items"
			INNER JOIN "product" ON "order_items"."product_id" = "product"."id"
			GROUP BY "product"."sizeIndex"
			ORDER BY "product"."sizeIndex" ASC
		`).Scan(&sizeWiseReport)

	if query.Error != nil {
		return nil, nil, nil, nil, nil, nil, e.DBQueryError_500(&query.Error)
	}

	for i := range sizeWiseReport {
		sizeWiseReport[i].SizeName = entities.Size[sizeWiseReport[i].SizeIndex].Size
	}

	var revenueGraph []entities.RevenueGraph
	query = repo.DB.Raw(`
			SELECT
				DATE ("orders"."order_date_and_time") AS "date",
				SUM (CASE WHEN "orders"."status" NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "sale",
				SUM (CASE WHEN "orders"."status" NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" ELSE 0 END) AS "quantity"
			FROM "order_items"
			INNER JOIN "orders" ON "order_items"."order_id" = "orders"."id"
			INNER JOIN "product" ON "order_items"."product_id" = "product"."id"
			INNER JOIN "dimensional_variants" ON "product"."dimensionalVariationID" = "dimensional_variants"."id"
			INNER JOIN "colour_variants" ON "dimensional_variants"."colourVariantId"= "colour_variants"."id"
			GROUP BY DATE ("orders"."order_date_and_time")
			ORDER BY DATE ("orders"."order_date_and_time") ASC
		`).Scan(&revenueGraph)

	if query.Error != nil {
		return nil, nil, nil, nil, nil, nil, e.DBQueryError_500(&query.Error)
	}

	for i := range revenueGraph {
		revenueGraph[i].Date = revenueGraph[i].Date[:10]
	}

	return &orderList, &sellerWiseReport, &brandWiseReport, &modelWiseReport, &sizeWiseReport, &revenueGraph, nil
}

func (repo *DashboardDataRepo) GetSalesReportBetweenDates(startDate time.Time, endDate time.Time) (
	*[]entities.SalesReportOrderList,
	*[]entities.SellerWiseReport,
	*[]entities.BrandWiseReport,
	*[]entities.ModelWiseReport,
	*[]entities.SizeWiseReport,
	*[]entities.RevenueGraph,
	*e.Error) {
	var orderList []entities.SalesReportOrderList
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
		WHERE "orders"."order_date_and_time" BETWEEN ? AND ?
		ORDER BY "orders"."order_date_and_time" DESC
	`, startDate, endDate).Scan(&orderList)

	if query.Error != nil {
		return nil, nil, nil, nil, nil, nil, e.DBQueryError_500(&query.Error)
	}

	var sellerWiseReport []entities.SellerWiseReport
	query = repo.DB.Raw(`
			SELECT
				"brands"."sellerId" AS "seller_id",
				"sellers"."firstName" AS "seller_name",
				SUM ("order_items"."quantity") AS "quantity_count",
				SUM ("order_items"."quantity" * "colour_variants"."salePrice") AS "total_sale",
				SUM (CASE WHEN "orders"."status" IN ('returned','return requested') THEN "order_items"."quantity" ELSE 0 END) AS "return_quantity_count",
				SUM (CASE WHEN "orders"."status" IN ('returned','return requested') THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "total_return_value",
				SUM (CASE WHEN "orders"."status" = 'cancelled' THEN "order_items"."quantity" ELSE 0 END) AS "cancel_quantity_count",
				SUM (CASE WHEN "orders"."status" = 'cancelled' THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "total_cancel_value",
				SUM (CASE WHEN "orders"."status" NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" ELSE 0 END) AS "effective_quantity_count",
				SUM (CASE WHEN "orders"."status" NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "effective_sale_value"
			FROM "order_items"
			INNER JOIN "orders" ON "order_items"."order_id" = "orders"."id"
			INNER JOIN "product" ON "order_items"."product_id" = "product"."id"
			INNER JOIN "dimensional_variants" ON "product"."dimensionalVariationID" = "dimensional_variants"."id"
			INNER JOIN "colour_variants" ON "dimensional_variants"."colourVariantId"= "colour_variants"."id"
			INNER JOIN "models" ON "colour_variants"."modelId" = "models"."id"
			INNER JOIN "brands" ON "models"."brandId" = "brands"."id"
			INNER JOIN "sellers" ON "brands"."sellerId" = "sellers"."id"
			WHERE "orders"."order_date_and_time" BETWEEN ? AND ?
			GROUP BY "brands"."sellerId", "sellers"."firstName"
			ORDER BY "brands"."sellerId" ASC
		`, startDate, endDate).Scan(&sellerWiseReport)

	if query.Error != nil {
		return nil, nil, nil, nil, nil, nil, e.DBQueryError_500(&query.Error)
	}

	var brandWiseReport []entities.BrandWiseReport
	query = repo.DB.Raw(`
			SELECT
				"brands"."id" AS "brand_id",
				"brands"."name" AS "brand_name",
				SUM ("order_items"."quantity") AS "quantity_count",
				SUM ("order_items"."quantity" * "colour_variants"."salePrice") AS "total_sale",
				SUM (CASE WHEN "orders"."status" IN ('returned','return requested') THEN "order_items"."quantity" ELSE 0 END) AS "return_quantity_count",
				SUM (CASE WHEN "orders"."status" IN ('returned','return requested') THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "total_return_value",
				SUM (CASE WHEN "orders"."status" = 'cancelled' THEN "order_items"."quantity" ELSE 0 END) AS "cancel_quantity_count",
				SUM (CASE WHEN "orders"."status" = 'cancelled' THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "total_cancel_value",
				SUM (CASE WHEN "orders"."status"  NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" ELSE 0 END) AS "effective_quantity_count",
				SUM (CASE WHEN "orders"."status" NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "effective_sale_value"
			FROM "order_items"
			INNER JOIN "orders" ON "order_items"."order_id" = "orders"."id"
			INNER JOIN "product" ON "order_items"."product_id" = "product"."id"
			INNER JOIN "dimensional_variants" ON "product"."dimensionalVariationID" = "dimensional_variants"."id"
			INNER JOIN "colour_variants" ON "dimensional_variants"."colourVariantId"= "colour_variants"."id"
			INNER JOIN "models" ON "colour_variants"."modelId" = "models"."id"
			INNER JOIN "brands" ON "models"."brandId" = "brands"."id"
			WHERE "orders"."order_date_and_time" BETWEEN ? AND ?
			GROUP BY "brands"."id"
			ORDER BY "brands"."id" ASC
		`, startDate, endDate).Scan(&brandWiseReport)

	if query.Error != nil {
		return nil, nil, nil, nil, nil, nil, e.DBQueryError_500(&query.Error)
	}

	var modelWiseReport []entities.ModelWiseReport
	query = repo.DB.Raw(`
			SELECT
				"models"."id" AS "model_id",
				"models"."name" AS "model_name",
				"brands"."name" AS "brand_name",
				SUM ("order_items"."quantity") AS "quantity_count",
				SUM ("order_items"."quantity" * "colour_variants"."salePrice") AS "total_sale",
				SUM (CASE WHEN "orders"."status" IN ('returned','return requested') THEN "order_items"."quantity" ELSE 0 END) AS "return_quantity_count",
				SUM (CASE WHEN "orders"."status" = 'cancelled' THEN "order_items"."quantity" ELSE 0 END) AS "cancel_quantity_count",
				SUM (CASE WHEN "orders"."status"  NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" ELSE 0 END) AS "effective_quantity_count",
				SUM (CASE WHEN "orders"."status" NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "effective_sale_value"
			FROM "order_items"
			INNER JOIN "orders" ON "order_items"."order_id" = "orders"."id"
			INNER JOIN "product" ON "order_items"."product_id" = "product"."id"
			INNER JOIN "dimensional_variants" ON "product"."dimensionalVariationID" = "dimensional_variants"."id"
			INNER JOIN "colour_variants" ON "dimensional_variants"."colourVariantId"= "colour_variants"."id"
			INNER JOIN "models" ON "colour_variants"."modelId" = "models"."id"
			INNER JOIN "brands" ON "models"."brandId" = "brands"."id"
			WHERE "orders"."order_date_and_time" BETWEEN ? AND ?
			GROUP BY "models"."id","brands"."name"
			ORDER BY "models"."id" ASC

		`, startDate, endDate).Scan(&modelWiseReport)

	if query.Error != nil {
		return nil, nil, nil, nil, nil, nil, e.DBQueryError_500(&query.Error)
	}

	var sizeWiseReport []entities.SizeWiseReport
	query = repo.DB.Raw(`
			SELECT
				"product"."sizeIndex" AS "size_index",
				SUM ("order_items"."quantity") AS "quantity_count"
			FROM "order_items"
			INNER JOIN "product" ON "order_items"."product_id" = "product"."id"
			INNER JOIN "orders" ON "order_items"."order_id" = "orders"."id"
			WHERE "orders"."order_date_and_time" BETWEEN ? AND ?
			GROUP BY "product"."sizeIndex"
			ORDER BY "product"."sizeIndex" ASC
		`, startDate, endDate).Scan(&sizeWiseReport)

	if query.Error != nil {
		return nil, nil, nil, nil, nil, nil, e.DBQueryError_500(&query.Error)
	}

	for i := range sizeWiseReport {
		sizeWiseReport[i].SizeName = entities.Size[sizeWiseReport[i].SizeIndex].Size
	}

	var revenueGraph []entities.RevenueGraph
	query = repo.DB.Raw(`
			SELECT
				DATE ("orders"."order_date_and_time") AS "date",
				SUM (CASE WHEN "orders"."status" NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "sale",
				SUM (CASE WHEN "orders"."status" NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" ELSE 0 END) AS "quantity"
			FROM "order_items"
			INNER JOIN "orders" ON "order_items"."order_id" = "orders"."id"
			INNER JOIN "product" ON "order_items"."product_id" = "product"."id"
			INNER JOIN "dimensional_variants" ON "product"."dimensionalVariationID" = "dimensional_variants"."id"
			INNER JOIN "colour_variants" ON "dimensional_variants"."colourVariantId"= "colour_variants"."id"
			WHERE "orders"."order_date_and_time" BETWEEN ? AND ?
			GROUP BY DATE ("orders"."order_date_and_time")
			ORDER BY DATE ("orders"."order_date_and_time") ASC
		`, startDate, endDate).Scan(&revenueGraph)

	if query.Error != nil {
		return nil, nil, nil, nil, nil, nil, e.DBQueryError_500(&query.Error)
	}

	for i := range revenueGraph {
		revenueGraph[i].Date = revenueGraph[i].Date[:10]
	}

	return &orderList, &sellerWiseReport, &brandWiseReport, &modelWiseReport, &sizeWiseReport, &revenueGraph, nil
}
