package orderManageUseCase

import (
	"MyShoo/internal/domain/entities"
	msg "MyShoo/internal/domain/messages"
	requestModels "MyShoo/internal/models/requestModels"
	requestValidation "MyShoo/pkg/validation"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/copier"
)

// CreateNewCoupon
func (uc *OrderUseCase) CreateNewCoupon(req *requestModels.NewCouponReq) (string, error) {
	// fmt.Println("req1= ", req)
	//check if coupon already exists

	//logical validations
	if req.Type == entities.Fixed && req.Percentage != 0 {
		return msg.Forbidden, errors.New("percentage should be 0 for fixed coupon type")
	}

	var coupon entities.Coupon
	if err := copier.Copy(&coupon, &req); err != nil {
		fmt.Println("Error occured while copying req to coupon")
		return "Some error occured.", err
	}

	startDate3, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		fmt.Println("Error occured while parsing start date")
		return "Some error occured.", err
	}
	endDate3, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		fmt.Println("Error occured while parsing end date")
		return "Some error occured.", err
	}
	startDate3 = startDate3.UTC().Add(-5*time.Hour - 30*time.Minute)
	endDate3 = endDate3.UTC().Add(-5*time.Hour - 30*time.Minute)

	startDate3 = startDate3.Local()
	endDate3 = endDate3.Local()

	//validate and set start and end time
	if startTime, err := requestValidation.ValidateAndParseDate(req.StartDate); err != nil {
		return msg.InvalidRequest, errors.New("invalid start time")
	} else {
		coupon.StartDate = startTime
	}
	if endTime, err := requestValidation.ValidateAndParseDate(req.EndDate); err != nil {
		return msg.InvalidRequest, errors.New("invalid end time")
	} else {
		endTime = endTime.AddDate(0, 0, 1) //to include the end day (upto 23:59:59)
		// endTime=endTime.Date()
		coupon.EndDate = endTime
	}

	coupon.StartDate = startDate3
	coupon.EndDate = endDate3
	// fmt.Println("coupon.EndDate= ", coupon.EndDate)
	// fmt.Println("coupon.StartDate= ", coupon.StartDate)
	couponExists, err := uc.orderRepo.DoCouponExistByCode(req.Code)
	if err != nil {
		fmt.Println("Error occured while checking if coupon exists")
		return "Some error occured.", err
	}
	if couponExists {
		return "Coupon already exists", errors.New("coupon already exists")
	}

	//initialise coupon
	coupon.Blocked = false

	// fmt.Println("coupon from uc= ", coupon)
	//create coupon
	err = uc.orderRepo.CreateNewCoupon(&coupon)
	if err != nil {
		fmt.Println("Error occured while creating coupon")
		return "Some error occured.", err
	}

	return "Coupon created successfully", nil
}

// BlockCoupon
func (uc *OrderUseCase) BlockCoupon(req *requestModels.BlockCouponReq) (string, error) {
	
	err := uc.orderRepo.BlockCoupon(req.ID)
	if err != nil {
		fmt.Println("Error occured while blocking coupon")
		return "Some error occured.", err
	}

	return "Coupon blocked successfully", nil
}

// UnblockCoupon
func (uc *OrderUseCase) UnblockCoupon(req *requestModels.UnblockCouponReq) (string, error) {

	err := uc.orderRepo.UnblockCoupon(req.ID)
	if err != nil {
		fmt.Println("Error occured while unblocking coupon")
		return "Some error occured.", err
	}

	return "Coupon unblocked successfully", nil
}

// GetAllCoupons
func (uc *OrderUseCase) GetAllCoupons() (*[]entities.Coupon, string, error) {
	coupons, message,err := uc.orderRepo.GetAllCoupons()
	if err != nil {
		fmt.Println("Error occured while getting all coupons")
		return nil, message, err
	}

	return coupons, "", nil
}

// GetExpiredCoupons
func (uc *OrderUseCase) GetExpiredCoupons() (*[]entities.Coupon, string, error) {
	coupons, message,err := uc.orderRepo.GetExpiredCoupons()
	if err != nil {
		fmt.Println("Error occured while getting expired coupons")
		return nil, message, err
	}

	return coupons, "", nil
}

// GetActiveCoupons
func (uc *OrderUseCase) GetActiveCoupons() (*[]entities.Coupon, string, error) {
	coupons, message,err := uc.orderRepo.GetActiveCoupons()
	if err != nil {
		fmt.Println("Error occured while getting active coupons")
		return nil, msg.RepoError, err
	}

	return coupons, message, nil
}

// GetUpcomingCoupons
func (uc *OrderUseCase) GetUpcomingCoupons() (*[]entities.Coupon, string, error) {
	coupons, message,err := uc.orderRepo.GetUpcomingCoupons()
	if err != nil {
		fmt.Println("Error occured while getting upcoming coupons")
		return nil, msg.RepoError, err
	}

	return coupons, message, nil
}