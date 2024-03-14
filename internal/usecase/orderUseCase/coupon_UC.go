package orderusecase

import (
	e "MyShoo/internal/domain/customErrors"
	"MyShoo/internal/domain/entities"
	request "MyShoo/internal/models/requestModels"
	requestValidation "MyShoo/pkg/validation"
	"time"

	"github.com/jinzhu/copier"
)

// CreateNewCoupon
func (uc *OrderUseCase) CreateNewCoupon(req *request.NewCouponReq) *e.Error {
	var err *e.Error
	//logical validations
	if req.Type == entities.Fixed && req.Percentage != 0 {
		return e.TextError("percentage should be 0 for fixed coupon type", 400)
	}

	var coupon entities.Coupon
	if errr := copier.Copy(&coupon, &req); errr != nil {
		return e.TextCumError("error occured while copying req to coupon", errr, 500)
	}

	startDate, errr := time.Parse("2006-01-02", req.StartDate)
	if errr != nil {
		return e.TextCumError("invalid start date. error on parsing:", errr, 400)
	}
	endDate, errr := time.Parse("2006-01-02", req.EndDate)
	if errr != nil {
		return e.TextCumError("invalid end date. error on parsing:", errr, 400)
	}
	startDate = startDate.UTC().Add(-5*time.Hour - 30*time.Minute)
	endDate = endDate.UTC().Add(-5*time.Hour - 30*time.Minute)

	startDate = startDate.Local()
	endDate = endDate.Local()

	//validate and set start and end time
	if startTime, errr := requestValidation.ValidateAndParseDate(req.StartDate); errr != nil {
		return e.TextError("invalid start time", 400)
	} else {
		coupon.StartDate = startTime
	}
	if endTime, err := requestValidation.ValidateAndParseDate(req.EndDate); err != nil {
		return e.TextError("invalid end time", 400)
	} else {
		endTime = endTime.AddDate(0, 0, 1) //to include the end day (upto 23:59:59)
		coupon.EndDate = endTime
	}

	coupon.StartDate = startDate
	coupon.EndDate = endDate
	codeAlreadyUsed, err := uc.orderRepo.IsCouponCodeTaken(req.Code)
	if err != nil {
		return err
	}
	if codeAlreadyUsed {
		return e.TextError("coupon code already exists", 400)
	}

	//initialise coupon
	coupon.Blocked = false

	//create coupon
	return uc.orderRepo.CreateNewCoupon(&coupon)
}

// BlockCoupon
func (uc *OrderUseCase) BlockCoupon(req *request.BlockCouponReq) *e.Error {
	return uc.orderRepo.BlockCoupon(req.ID)
}

// UnblockCoupon
func (uc *OrderUseCase) UnblockCoupon(req *request.UnblockCouponReq) *e.Error {
	return uc.orderRepo.UnblockCoupon(req.ID)
}

// GetAllCoupons
func (uc *OrderUseCase) GetAllCoupons() (*[]entities.Coupon, *e.Error) {
	return uc.orderRepo.GetAllCoupons()
}

// GetExpiredCoupons
func (uc *OrderUseCase) GetExpiredCoupons() (*[]entities.Coupon, *e.Error) {
	return uc.orderRepo.GetExpiredCoupons()
}

// GetActiveCoupons
func (uc *OrderUseCase) GetActiveCoupons() (*[]entities.Coupon, *e.Error) {
	return uc.orderRepo.GetActiveCoupons()
}

// GetUpcomingCoupons
func (uc *OrderUseCase) GetUpcomingCoupons() (*[]entities.Coupon, *e.Error) {
	return uc.orderRepo.GetUpcomingCoupons()
}
