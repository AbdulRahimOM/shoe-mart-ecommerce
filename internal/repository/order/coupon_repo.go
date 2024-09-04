package orderrepo

import (
	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
)

var (
	errCouponDoesNotExist_400 = e.Error{StatusCode: 400, Status: "Failed", Msg: "Coupon doesn't exist", Err: nil}
)

func (repo *OrderRepo) IsCouponCodeTaken(code string) (bool, *e.Error) {
	var temp entities.Coupon
	query := repo.DB.Raw(`
		SELECT *
		FROM coupons
		WHERE code = ?`,
		code).Scan(&temp)

	if query.Error != nil {
		return false, e.DBQueryError_500(&query.Error)
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (repo *OrderRepo) CreateNewCoupon(coupon *entities.Coupon) *e.Error {
	result := repo.DB.Create(&coupon)
	if result.Error != nil {
		return e.DBQueryError_500(&result.Error)
	}

	return nil
}

func (repo *OrderRepo) BlockCoupon(couponID uint) *e.Error {

	result := repo.DB.Model(&entities.Coupon{}).Where("id = ?", couponID).Update("blocked", true)
	if result.Error != nil {
		return e.DBQueryError_500(&result.Error)
	}
	return nil
}

func (repo *OrderRepo) UnblockCoupon(couponID uint) *e.Error {

	result := repo.DB.Model(&entities.Coupon{}).Where("id = ?", couponID).Update("blocked", false)
	if result.Error != nil {
		return e.DBQueryError_500(&result.Error)
	}
	return nil
}

// GetAllCoupons
func (repo *OrderRepo) GetAllCoupons() (*[]entities.Coupon, *e.Error) {
	var coupons []entities.Coupon
	result := repo.DB.Find(&coupons)
	if result.Error != nil {
		return nil, e.DBQueryError_500(&result.Error)
	}

	return &coupons, nil
}

// GetExpiredCoupons
func (repo *OrderRepo) GetExpiredCoupons() (*[]entities.Coupon, *e.Error) {
	var coupons []entities.Coupon
	result := repo.DB.Where("end_date < now()").Find(&coupons)
	if result.Error != nil {
		return nil, e.DBQueryError_500(&result.Error)
	}

	return &coupons, nil
}

// GetActiveCoupons
func (repo *OrderRepo) GetActiveCoupons() (*[]entities.Coupon, *e.Error) {
	var coupons []entities.Coupon
	result := repo.DB.Where("start_date < now() AND end_date > now() AND blocked=?", "false").Find(&coupons)
	if result.Error != nil {
		return nil, e.DBQueryError_500(&result.Error)
	}

	return &coupons, nil
}

// GetUpcomingCoupons
func (repo *OrderRepo) GetUpcomingCoupons() (*[]entities.Coupon, *e.Error) {
	var coupons []entities.Coupon
	result := repo.DB.Where("start_date > now()").Find(&coupons)
	if result.Error != nil {
		return nil, e.DBQueryError_500(&result.Error)
	}

	return &coupons, nil
}

// GetCouponByID
func (repo *OrderRepo) GetCouponByID(couponID uint) (*entities.Coupon, *e.Error) {
	var coupon entities.Coupon
	result := repo.DB.Where("id = ?", couponID).Find(&coupon)
	if result.Error != nil {
		return nil, e.DBQueryError_500(&result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, &errCouponDoesNotExist_400

	}

	return &coupon, nil
}

// GetCouponUsageCount implements repo.IOrderRepo.
func (repo *OrderRepo) GetCouponUsageCount(userID uint, couponID uint) (uint, *e.Error) {
	var count uint
	result := repo.DB.Raw(`
		SELECT COUNT(*)
		FROM orders
		WHERE user_id = ? AND coupon_id = ?`,
		userID, couponID).Scan(&count)
	if result.Error != nil {
		return 0, e.DBQueryError_500(&result.Error)
	}

	return count, nil
}
