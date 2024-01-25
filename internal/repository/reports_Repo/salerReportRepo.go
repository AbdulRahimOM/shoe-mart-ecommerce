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

	return url, nil
}

// GetSalesReportFullTime
func (repo *DashboardDataRepo) GetSalesReportFullTime() (*[]entities.SalesReportOrderList, *[]entities.SellerWiseReport, error) {
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

	//

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return nil, nil, query.Error
	}
	// fmt.Println("orderList=", orderList)

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
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return nil, nil, query.Error
	}

	fmt.Println("sellerWiseReport=", sellerWiseReport)

	var qsum *uint
	query = repo.DB.Raw(`
			SELECT
				SUM (CASE WHEN "orders"."status" NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" ELSE 0 END) AS "effective_quantity_count"
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
		`).Scan(&qsum)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return nil, nil, query.Error
	}

	fmt.Println("qsum=",*qsum)
	if qsum == nil {
		fmt.Println("qsum is nil")
	}else{
		fmt.Println("qsum is not nil")
	}

	// 	var brandWiseReport []entities.BrandWiseReport
	// 	query = repo.DB.Raw(`
	// 		SELECT
	// 			"brands"."id" AS "brand_id",
	// 			"brands"."name" AS "brand_name",
	// 			SUM ("order_items"."quantity") AS "quantity_count",
	// 			SUM ("order_items"."quantity" * "colour_variants"."salePrice") AS "total_sale",
	// 			SUM (CASE WHEN "orders"."status" IN ('returned','return requested') THEN "order_items"."quantity" ELSE 0 END) AS "return_quantity_count",
	// 			SUM (CASE WHEN "orders"."status" IN ('returned','return requested') THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "total_return_value",
	// 			SUM (CASE WHEN "orders"."status" = 'cancelled' THEN "order_items"."quantity" ELSE 0 END) AS "cancel_quantity_count",
	// 			SUM (CASE WHEN "orders"."status" = 'cancelled' THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "total_cancel_value",
	// 			SUM (CASE WHEN "orders"."order_status"  NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" ELSE 0 END) AS "effective_quantity_count",
	// 			SUM (CASE WHEN "orders"."order_status" NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "effective_sale_value"
	// 		FROM "order_items"
	// 		INNER JOIN "orders" ON "order_items"."order_id" = "orders"."id"
	// 		INNER JOIN "product" ON "order_items"."product_id" = "product"."id"
	// 		INNER JOIN "dimensional_variants" ON "product"."dimensionalVariationID" = "dimensional_variants"."id"
	// 		INNER JOIN "colour_variants" ON "dimensional_variants"."colourVariantId"= "colour_variants"."id"
	// 		INNER JOIN "models" ON "colour_variants"."modelId" = "models"."id"
	// 		INNER JOIN "brands" ON "models"."brandId" = "brands"."id"
	// 		GROUP BY "brands"."id"
	// 		ORDER BY "brands"."id" ASC
	// 	`).Scan(&brandWiseReport)

	// 	if query.Error != nil {
	// 		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
	// 		return nil, query.Error
	// 	}

	// 	var modelWiseReport []entities.ModelWiseReport
	// 	query = repo.DB.Raw(`
	// 		SELECT
	// 			"models"."id" AS "model_id",
	// 			"models"."name" AS "model_name",
	// 			"brands"."name" AS "brand_name",
	// 			SUM ("order_items"."quantity") AS "quantity_count",
	// 			SUM ("order_items"."quantity" * "colour_variants"."salePrice") AS "total_sale",
	// 			SUM (CASE WHEN "orders"."status" IN ('returned','return requested') THEN "order_items"."quantity" ELSE 0 END) AS "return_quantity_count",
	// 			SUM (CASE WHEN "orders"."status" = 'cancelled' THEN "order_items"."quantity" ELSE 0 END) AS "cancel_quantity_count",
	// 			SUM (CASE WHEN "orders"."order_status"  NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" ELSE 0 END) AS "effective_quantity_count",
	// 			SUM (CASE WHEN "orders"."order_status" NOT IN ('cancelled','returned','return requested') THEN "order_items"."quantity" * "colour_variants"."salePrice" ELSE 0 END) AS "effective_sale_value"
	// 		FROM "order_items"
	// 		INNER JOIN "orders" ON "order_items"."order_id" = "orders"."id"
	// 		INNER JOIN "product" ON "order_items"."product_id" = "product"."id"
	// 		INNER JOIN "dimensional_variants" ON "product"."dimensionalVariationID" = "dimensional_variants"."id"
	// 		INNER JOIN "colour_variants" ON "dimensional_variants"."colourVariantId"= "colour_variants"."id"
	// 		INNER JOIN "models" ON "colour_variants"."modelId" = "models"."id"
	// 		INNER JOIN "brands" ON "models"."brandId" = "brands"."id"
	// 		GROUP BY "models"."id" ASC
	// 	`).Scan(&modelWiseReport)

	// 	if query.Error != nil {
	// 		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
	// 		return nil, query.Error
	// 	}

	// var sizeWiseReport []entities.SizeWiseReport
	// 	query = repo.DB.Raw(`
	// 		SELECT
	// 			"product"."sizeIndex" AS "size_index",
	// 			SUM ("order_items"."quantity") AS "quantity_count"
	// 		FROM "order_items"
	// 		INNER JOIN "product" ON "order_items"."product_id" = "product"."id"
	// 		GROUP BY "product"."sizeIndex"
	// 		ORDER BY "product"."sizeIndex" ASC
	// 	`).Scan(&sizeWiseReport)

	// 	if query.Error != nil {
	// 		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
	// 		return nil, query.Error
	// 	}

	return &orderList, &sellerWiseReport, nil
}
