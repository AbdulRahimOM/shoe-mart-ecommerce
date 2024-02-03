package entities

import (
	"database/sql/driver"
	"errors"
	"time"
)

type Coupon struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	Code          string     `json:"code" gorm:"column:code;notNull"`
	Type          CouponType `json:"type" gorm:"column:type;notNull"`
	MinOrderValue float32    `json:"minOrderValue" gorm:"column:min_order_value"`
	MaxDiscount   float32    `json:"maxDiscount" gorm:"column:max_discount"`
	Percentage    float32    `json:"percentage" gorm:"column:percentage"`
	Description   string     `json:"description" gorm:"column:description;notNull"`
	UsageLimit    uint       `json:"usageLimit" gorm:"column:usage_limit;notNull"`
	StartDate     time.Time  `json:"startDate" gorm:"column:start_date;notNull"`
	EndDate       time.Time  `json:"endDate" gorm:"column:end_date;notNull"`
	Blocked       bool       `json:"blocked" gorm:"column:blocked;notNull"`
}

type CouponType string

const (
	Percentage CouponType = "percentage"
	Fixed      CouponType = "fixed"
)

// implementing Scanner and Valuer interface for CouponType
func (c *CouponType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*c = CouponType(v)
	case string:
		*c = CouponType(v)
	default:
		return errors.New("unexpected type for CouponType")
	}
	return nil
}

func (c CouponType) Value() (driver.Value, error) {
	return string(c), nil
}
