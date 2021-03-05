package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

var (
	// ErrInvalidCoupon tells that the coupon is not valid, whether has not been started, expired, or limit exceeded.
	ErrInvalidCoupon = errors.New("error invalid coupon")
	// ErrCouponNotFound tells that order can not be found
	ErrCouponNotFound = errors.New("coupon not found")
)

// Coupon is price reduction scheme that can be applied to an order
type Coupon struct {
	Code     string          `bson:"code" json:"code"`
	Quantity int             `bson:"quantity" json:"quantity"`
	Amount   decimal.Decimal `bson:"amount" json:"amount"`
	Begin    time.Time       `bson:"begin" json:"begin"`
	End      time.Time       `bson:"end" json:"end"`
	Type     CouponType      `bson:"type" json:"type"`
}

// CouponType type of the coupon
type CouponType int

const (
	// CouponTypePercentage represents the Amount of the coupon is percentage
	CouponTypePercentage CouponType = iota + 1
	// CouponTypeNominal represents the Amount of the coupon is a nominal of money
	CouponTypeNominal
)

func (s CouponType) String() string {
	switch s {
	case CouponTypePercentage:
		return "Percentage"
	case CouponTypeNominal:
		return "Nominal"
	}
	return ""
}

// Validate checks whether coupon is valid or not
func (c *Coupon) Validate() error {
	if c.Quantity == 0 {
		return ErrInvalidCoupon
	}
	now := time.Now()
	if now.Before(c.Begin) || now.After(c.End) {
		return ErrInvalidCoupon
	}
	return nil
}

// GetPriceAfterReduction applies coupon reduction to the price
func (c *Coupon) GetPriceAfterReduction(price decimal.Decimal) decimal.Decimal {
	switch c.Type {
	case CouponTypeNominal:
		return price.Sub(c.Amount)
	case CouponTypePercentage:
		return price.Sub(price.Mul(c.Amount))
	}
	return price
}

// CouponRepository provides access to coupons
type CouponRepository interface {
	FindByCode(ctx context.Context, code string) (*Coupon, error)
	Update(ctx context.Context, coupon *Coupon) error
}
