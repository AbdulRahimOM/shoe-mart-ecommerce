package reportsrepo

import (
	e "MyShoo/internal/domain/customErrors"
	response "MyShoo/internal/models/responseModels"
)

// GetTopModels implements repo.IReportsRepo.
func (repo *DashboardDataRepo) GetTopModels(limit int) (*[]response.TopModels, *e.Error) {
	var topModels = make([]response.TopModels, limit)

	query := repo.DB.Raw(`
		SELECT
			models.id as "model_id",
			models.name as "model_name",
			brands.name as "brand_name",
			CONCAT(sellers."firstName", ' ', sellers."lastName") as "seller_name",
			sum(order_items.quantity) as "quantity_sold",
			sum(order_items.sale_price_on_order*order_items.quantity) as "net_sale_amount",
			sum(colour_variants.mrp*order_items.quantity) as "net_mrp_value_sold"
		FROM
			order_items
		JOIN product ON order_items.product_id = product.id
		JOIN dimensional_variants ON product."dimensionalVariationID" = dimensional_variants.id
		JOIN colour_variants ON dimensional_variants."colourVariantId" = colour_variants.id
		JOIN models ON colour_variants."modelId" = models.id
		JOIN brands ON models."brandId" = brands.id
		JOIN sellers ON brands."sellerId" = sellers.id
		GROUP BY
			models.id,brands.name,sellers."firstName",sellers."lastName"
		ORDER BY
			net_sale_amount DESC
		LIMIT ?
	`, limit).Scan(&topModels)
	if query.Error != nil {
		return nil, e.DBQueryError_500(&query.Error)
	}

	return &topModels, nil
}

// GetTopProducts implements repo.IReportsRepo.
func (repo *DashboardDataRepo) GetTopProducts(limit int) (*[]response.TopProducts, *e.Error) {
	var topProducts = make([]response.TopProducts, limit)

	query := repo.DB.Raw(`
		SELECT
			colour_variants.id as "colour_variant_id",
			models.name as "model_name",
			brands.name as "brand_name",
			CONCAT(sellers."firstName", ' ', sellers."lastName") as "seller_name",
			sum(order_items.quantity) as "quantity_sold",
			sum(colour_variants.mrp*order_items.quantity) as "net_mrp_value",
			sum(order_items.sale_price_on_order*order_items.quantity) as "net_sale_amount",
			colour_variants.mrp as "current_mrp",
			colour_variants."salePrice" as "current_sale_price"
		FROM
			order_items
		JOIN product ON order_items.product_id = product.id
		JOIN dimensional_variants ON product."dimensionalVariationID" = dimensional_variants.id
		JOIN colour_variants ON dimensional_variants."colourVariantId" = colour_variants.id
		JOIN models ON colour_variants."modelId" = models.id
		JOIN brands ON models."brandId" = brands.id
		JOIN sellers ON brands."sellerId" = sellers.id
		GROUP BY
			colour_variants.id,models.name,brands.name,sellers."firstName",sellers."lastName"
		ORDER BY
			net_sale_amount DESC
		LIMIT ?
	`, limit).Scan(&topProducts)
	if query.Error != nil {
		return nil, e.DBQueryError_500(&query.Error)
	}

	return &topProducts, nil
}

// GetTopBrands implements repo.IReportsRepo.
func (repo *DashboardDataRepo) GetTopBrands(limit int) (*[]response.TopBrands, *e.Error) {
	var topBrands = make([]response.TopBrands, limit)

	query := repo.DB.Raw(`
		SELECT
			brands.id as "brand_id",
			brands.name as "brand_name",
			CONCAT(sellers."firstName", ' ', sellers."lastName") as "seller_name",
			sum(order_items.quantity) as "quantity_sold",
			sum(order_items.sale_price_on_order*order_items.quantity) as "net_sale_amount",
			sum(colour_variants.mrp*order_items.quantity) as "net_mrp_value_sold"
		FROM
			order_items
		JOIN product ON order_items.product_id = product.id
		JOIN dimensional_variants ON product."dimensionalVariationID" = dimensional_variants.id
		JOIN colour_variants ON dimensional_variants."colourVariantId" = colour_variants.id
		JOIN models ON colour_variants."modelId" = models.id
		JOIN brands ON models."brandId" = brands.id
		JOIN sellers ON brands."sellerId" = sellers.id
		GROUP BY
			brands.id,brands.name,sellers."firstName",sellers."lastName"
		ORDER BY
			net_sale_amount DESC
		LIMIT ?
	`, limit).Scan(&topBrands)
	if query.Error != nil {
		return nil, e.DBQueryError_500(&query.Error)
	}

	return &topBrands, nil

}

// GetTopSellers implements repo.IReportsRepo.
func (repo *DashboardDataRepo) GetTopSellers(limit int) (*[]response.TopSellers, *e.Error) {
	var topSellers = make([]response.TopSellers, limit)

	query := repo.DB.Raw(`
		SELECT
			sellers.id as "seller_id",
			CONCAT(sellers."firstName", ' ', sellers."lastName") as "seller_name",
			sum(order_items.quantity) as "quantity_sold",
			sum(order_items.sale_price_on_order*order_items.quantity) as "net_sale_amount",
			sum(colour_variants.mrp*order_items.quantity) as "net_mrp_value_sold"
		FROM
			order_items
		JOIN product ON order_items.product_id = product.id
		JOIN dimensional_variants ON product."dimensionalVariationID" = dimensional_variants.id
		JOIN colour_variants ON dimensional_variants."colourVariantId" = colour_variants.id
		JOIN models ON colour_variants."modelId" = models.id
		JOIN brands ON models."brandId" = brands.id
		JOIN sellers ON brands."sellerId" = sellers.id
		GROUP BY
			sellers.id,sellers."firstName",sellers."lastName"
		ORDER BY
			net_sale_amount DESC
		LIMIT ?
	`, limit).Scan(&topSellers)
	if query.Error != nil {
		return nil, e.DBQueryError_500(&query.Error)
	}

	return &topSellers, nil
}
