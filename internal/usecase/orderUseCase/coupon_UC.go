package orderusecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	msg "MyShoo/internal/domain/messages"
	request "MyShoo/internal/models/requestModels"
	requestValidation "MyShoo/pkg/validation"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/copier"
)

// CreateNewCoupon
func (uc *OrderUseCase) CreateNewCoupon(req *request.NewCouponReq) error {
	// fmt.Println("req1= ", req)
	//check if coupon already exists

	//logical validations
	if req.Type == entities.Fixed && req.Percentage != 0 {
		return  e.Error{Err:errors.New("percentage should be 0 for fixed coupon type"),StatusCode: 400}
	}

	var coupon entities.Coupon
	if err := copier.Copy(&coupon, &req); err != nil {
		err=fmt.Errorf("Error occured while copying req to coupon,error:",err)
		return  e.Error{Err: err,StatusCode: 500}
	}

	startDate3, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		err=fmt.Errorf("Error occured while parsing start date",err)
		return e.Error{Err: err,StatusCode: 400}
	}
	endDate3, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		err=fmt.Errorf("Error occured while parsing end date",err)
		return e.Error{Err: err,StatusCode: 400}
	}
	startDate3 = startDate3.UTC().Add(-5*time.Hour - 30*time.Minute)
	endDate3 = endDate3.UTC().Add(-5*time.Hour - 30*time.Minute)

	startDate3 = startDate3.Local()
	endDate3 = endDate3.Local()

	//validate and set start and end time
	if startTime, err := requestValidation.ValidateAndParseDate(req.StartDate); err != nil {
		return  e.Error{Err:  errors.New("invalid start time"),StatusCode: 400}
	} else {
		coupon.StartDate = startTime
	}
	if endTime, err := requestValidation.ValidateAndParseDate(req.EndDate); err != nil {
		return  e.Error{Err:  errors.New("invalid end time"),StatusCode: 400}
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
func (uc *OrderUseCase) BlockCoupon(req *request.BlockCouponReq) error {

	err := uc.orderRepo.BlockCoupon(req.ID)
	if err != nil {
		return err
	}

	return nil
}

// UnblockCoupon
func (uc *OrderUseCase) UnblockCoupon(req *request.UnblockCouponReq) error {

	err := uc.orderRepo.UnblockCoupon(req.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetAllCoupons
func (uc *OrderUseCase) GetAllCoupons() (*[]entities.Coupon, error) {
	coupons, err := uc.orderRepo.GetAllCoupons()
	if err != nil {
		fmt.Println("Error occured while getting all coupons")
		return nil, err
	}

	return coupons, nil
}

// GetExpiredCoupons
func (uc *OrderUseCase) GetExpiredCoupons() (*[]entities.Coupon, error) {
	coupons, err := uc.orderRepo.GetExpiredCoupons()
	if err != nil {
		fmt.Println("Error occured while getting expired coupons")
		return nil, err
	}

	return coupons, nil
}

// GetActiveCoupons
func (uc *OrderUseCase) GetActiveCoupons() (*[]entities.Coupon, error) {
	coupons, err := uc.orderRepo.GetActiveCoupons()
	if err != nil {
		fmt.Println("Error occured while getting active coupons")
		return nil, err
	}

	return coupons, nil
}

// GetUpcomingCoupons
func (uc *OrderUseCase) GetUpcomingCoupons() (*[]entities.Coupon, error) {
	coupons, err := uc.orderRepo.GetUpcomingCoupons()
	if err != nil {
		fmt.Println("Error occured while getting upcoming coupons")
		return nil, err
	}

	return coupons, nil
}
