package reportsrepo

import (
	"time"

	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
)

func (repo *DashboardDataRepo) GetDashBoardDataBetweenDates(start time.Time, end time.Time) (*entities.DashboardData, *[]entities.SalePerDay, *e.Error) {
	var dashBoardData entities.DashboardData
	var salePerDay []entities.SalePerDay

	err := repo.DB.Raw(`
		SELECT COUNT(*) AS order_count, 
		SUM(final_amount) AS net_sale_value,
		SUM(coupon_discount) AS coupon_discounts,
		SUM(original_amount) AS net_original_value,
		SUM(CASE WHEN status = 'cancelled' THEN 1 ELSE 0 END) AS cancelled_order_count,
		SUM(CASE WHEN status = 'returned' THEN 1 ELSE 0 END) AS returned_order_count,
		SUM(CASE WHEN status = 'cancelled' THEN 0 ELSE final_amount END) AS sale_value_after_cancellation_and_returns,
		(SELECT COUNT(*) FROM users WHERE created_at BETWEEN ? AND ?) AS users_registered
		FROM orders 
		WHERE order_date_and_time BETWEEN ? AND ?`,
		start, end, start, end).Scan(&dashBoardData).Error
	if err != nil {
		return nil, nil, e.DBQueryError_500(&err)
	}

	err = repo.DB.Raw(`
    SELECT TO_CHAR(order_date_and_time, 'YYYY-MM-DD') AS date, SUM(final_amount) AS sale
    FROM orders
    WHERE order_date_and_time BETWEEN ? AND ?
    GROUP BY TO_CHAR(order_date_and_time, 'YYYY-MM-DD')`,
		start, end).Scan(&salePerDay).Error
	if err != nil {
		return nil, nil, e.DBQueryError_500(&err)
	}

	return &dashBoardData, &salePerDay, nil
}

func (repo *DashboardDataRepo) GetDashBoardDataFullTime() (*entities.DashboardData, *[]entities.SalePerDay, *e.Error) {
	var dashBoardData entities.DashboardData
	var salePerDay []entities.SalePerDay

	err := repo.DB.Raw(`
		SELECT COUNT(*) AS order_count, 
		SUM(final_amount) AS net_sale_value,
		SUM(coupon_discount) AS coupon_discounts,
		SUM(original_amount) AS net_original_value,
		SUM(CASE WHEN status = 'cancelled' THEN 1 ELSE 0 END) AS cancelled_order_count,
		SUM(CASE WHEN status = 'returned' THEN 1 ELSE 0 END) AS returned_order_count,
		SUM(CASE WHEN status = 'cancelled' THEN 0 ELSE final_amount END) AS sale_value_after_cancellation_and_returns,
		(SELECT COUNT(*) FROM users) AS users_registered
		FROM orders`,
	).Scan(&dashBoardData).Error
	if err != nil {
		return nil, nil, e.DBQueryError_500(&err)
	}

	err = repo.DB.Raw(`
    SELECT TO_CHAR(order_date_and_time, 'YYYY-MM-DD') AS date, SUM(final_amount) AS sale
    FROM orders
    GROUP BY TO_CHAR(order_date_and_time, 'YYYY-MM-DD')`,
	).Scan(&salePerDay).Error
	if err != nil {
		return nil, nil, e.DBQueryError_500(&err)
	}

	return &dashBoardData, &salePerDay, nil
}
