package ordermanagementrepo

import (
	"MyShoo/internal/domain/entities"
	msg "MyShoo/internal/domain/messages"
	"errors"
	"fmt"
)

func (repo *OrderRepo) DoCouponExistByCode(code string) (bool, error) {
	var temp entities.Coupon
	query := repo.DB.Raw(`
		SELECT *
		FROM coupons
		WHERE code = ?`,
		code).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-coupon is existing or not. query.Error= ", query.Error, "\n----")
		return false, query.Error
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (repo *OrderRepo) CreateNewCoupon(coupon *entities.Coupon) error {
	fmt.Println("coupon.StartDate= ", coupon.StartDate)
	fmt.Println("coupon.EndDate= ", coupon.EndDate)
	// starr
	result := repo.DB.Create(&coupon)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't create new coupon. query.Error= ", result.Error, "\n----")
		return result.Error
	}

	return nil
}

func (repo *OrderRepo) BlockCoupon(couponID uint) error {

	result := repo.DB.Model(&entities.Coupon{}).Where("id = ?", couponID).Update("blocked", true)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't block coupon. query.Error= ", result.Error, "\n----")
		return result.Error
	}
	return nil
}

func (repo *OrderRepo) UnblockCoupon(couponID uint) error {

	result := repo.DB.Model(&entities.Coupon{}).Where("id = ?", couponID).Update("blocked", false)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't unblock coupon. query.Error= ", result.Error, "\n----")
		return result.Error
	}
	return nil
}

// GetAllCoupons
func (repo *OrderRepo) GetAllCoupons() (*[]entities.Coupon, string, error) {
	var coupons []entities.Coupon
	result := repo.DB.Find(&coupons)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get all coupons. query.Error= ", result.Error, "\n----")
		return nil, "Some error occured", result.Error
	}

	return &coupons, "", nil
}

// GetExpiredCoupons
func (repo *OrderRepo) GetExpiredCoupons() (*[]entities.Coupon, string, error) {
	var coupons []entities.Coupon
	result := repo.DB.Where("end_date < now()").Find(&coupons)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get expired coupons. query.Error= ", result.Error, "\n----")
		return nil, "Some error occured", result.Error
	}

	return &coupons, "", nil
}

// GetActiveCoupons
func (repo *OrderRepo) GetActiveCoupons() (*[]entities.Coupon, string, error) {
	var coupons []entities.Coupon
	result := repo.DB.Where("start_date < now() AND end_date > now() AND blocked=?", "false").Find(&coupons)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get active coupons. query.Error= ", result.Error, "\n----")
		return nil, "Some error occured", result.Error
	}

	return &coupons, "", nil
}

// GetUpcomingCoupons
func (repo *OrderRepo) GetUpcomingCoupons() (*[]entities.Coupon, string, error) {
	var coupons []entities.Coupon
	result := repo.DB.Where("start_date > now()").Find(&coupons)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get upcoming coupons. query.Error= ", result.Error, "\n----")
		return nil, "Some error occured", result.Error
	}

	return &coupons, "", nil
}

// GetCouponByID
func (repo *OrderRepo) GetCouponByID(couponID uint) (*entities.Coupon, string, error) {
	var coupon entities.Coupon
	result := repo.DB.Where("id = ?", couponID).Find(&coupon)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get coupon by id. query.Error= ", result.Error, "\n----")
		return nil, msg.RepoError, result.Error
	}else if result.RowsAffected == 0 {
		return nil, msg.InvalidRequest, errors.New("coupon doesn't exist")
	}

	return &coupon, "", nil
}

// GetCouponUsageCount implements repository_interface.IOrderRepo.
func (repo *OrderRepo) GetCouponUsageCount(userID uint, couponID uint) (uint, string, error) {
	var count uint
	result := repo.DB.Raw(`
		SELECT COUNT(*)
		FROM orders
		WHERE user_id = ? AND coupon_id = ?`,
		userID, couponID).Scan(&count)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't get coupon usage count. query.Error= ", result.Error, "\n----")
		return 0, msg.RepoError, result.Error
	}

	return count, "", nil
}