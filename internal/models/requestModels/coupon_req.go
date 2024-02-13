package request

import "MyShoo/internal/domain/entities"

type NewCouponReq struct {
	Code          string              `json:"code" validate:"required,min=3,max=15"`
	Type          entities.CouponType `json:"type" validate:"required,oneof=percentage fixed"`
	MinOrderValue float32             `json:"minOrderValue" validate:"required,number,min=0"`
	MaxDiscount   float32             `json:"maxDiscount" validate:"required,number,min=0"`
	Percentage    float32             `json:"percentage" validate:"number,min=0,max=100"`
	Description   string              `json:"description" validate:"min=3,max=100"`
	UsageLimit    uint                `json:"usageLimit" validate:"required,number,min=1"`
	StartDate     string              `json:"startDate" validate:"required"`
	EndDate       string              `json:"endDate" validate:"required"`
}

type BlockCouponReq struct {
	ID uint `json:"id" validate:"required,number"`
}

type UnblockCouponReq struct {
	ID uint `json:"id" validate:"required,number"`
}
