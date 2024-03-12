package orderrepo

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	"errors"
	"fmt"
)

func (repo *OrderRepo) DoCouponExistByCode(code string) (bool, *e.Error) {
	var temp entities.Coupon
	query := repo.DB.Raw(`
		SELECT *
		FROM coupons
		WHERE code = ?`,
		code).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-coupon is existing or not. query.Error= ", query.Error, "\n----")
		return false, &e.Error{Err: query.Error, StatusCode: 500}
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (repo *OrderRepo) CreateNewCoupon(coupon *entities.Coupon) *e.Error {
	fmt.Println("coupon.StartDate= ", coupon.StartDate)
	fmt.Println("coupon.EndDate= ", coupon.EndDate)
	// starr
	result := repo.DB.Create(&coupon)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't create new coupon. query.Error= ", result.Error, "\n----")
		return &e.Error{Err: result.Error, StatusCode: 500}
	}

	return nil
}

func (repo *OrderRepo) BlockCoupon(couponID uint) *e.Error {

	result := repo.DB.Model(&entities.Coupon{}).Where("id = ?", couponID).Update("blocked", true)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't block coupon. query.Error= ", result.Error, "\n----")
		return &e.Error{Err: result.Error, StatusCode: 500}
	}
	return nil
}

func (repo *OrderRepo) UnblockCoupon(couponID uint) *e.Error {

	result := repo.DB.Model(&entities.Coupon{}).Where("id = ?", couponID).Update("blocked", false)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't unblock coupon. query.Error= ", result.Error, "\n----")
		return &e.Error{Err: result.Error, StatusCode: 500}
	}
	return nil
}

// GetAllCoupons
func (repo *OrderRepo) GetAllCoupons() (*[]entities.Coupon, *e.Error) {
	var coupons []entities.Coupon
	result := repo.DB.Find(&coupons)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get all coupons. query.Error= ", result.Error, "\n----")
		return nil, &e.Error{Err: result.Error, StatusCode: 500}
	}

	return &coupons, nil
}

// GetExpiredCoupons
func (repo *OrderRepo) GetExpiredCoupons() (*[]entities.Coupon, *e.Error) {
	var coupons []entities.Coupon
	result := repo.DB.Where("end_date < now()").Find(&coupons)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get expired coupons. query.Error= ", result.Error, "\n----")
		return nil, &e.Error{Err: result.Error, StatusCode: 500}
	}

	return &coupons, nil
}

// GetActiveCoupons
func (repo *OrderRepo) GetActiveCoupons() (*[]entities.Coupon, *e.Error) {
	var coupons []entities.Coupon
	result := repo.DB.Where("start_date < now() AND end_date > now() AND blocked=?", "false").Find(&coupons)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get active coupons. query.Error= ", result.Error, "\n----")
		return nil, &e.Error{Err: result.Error, StatusCode: 500}
	}

	return &coupons, nil
}

// GetUpcomingCoupons
func (repo *OrderRepo) GetUpcomingCoupons() (*[]entities.Coupon, *e.Error) {
	var coupons []entities.Coupon
	result := repo.DB.Where("start_date > now()").Find(&coupons)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get upcoming coupons. query.Error= ", result.Error, "\n----")
		return nil, &e.Error{Err: result.Error, StatusCode: 500}
	}

	return &coupons, nil
}

// GetCouponByID
func (repo *OrderRepo) GetCouponByID(couponID uint) (*entities.Coupon, *e.Error) {
	var coupon entities.Coupon
	result := repo.DB.Where("id = ?", couponID).Find(&coupon)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get coupon by id. query.Error= ", result.Error, "\n----")
		return nil, &e.Error{Err: result.Error, StatusCode: 500}
	}
	if result.RowsAffected == 0 {
		return nil, &e.Error{Err: errors.New("coupon doesn't exist"), StatusCode: 400}
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
		fmt.Println("-------\nquery error happened. couldn't get coupon usage count. query.Error= ", result.Error, "\n----")
		return 0, &e.Error{Err: result.Error, StatusCode: 500}
	}

	return count, nil
}
