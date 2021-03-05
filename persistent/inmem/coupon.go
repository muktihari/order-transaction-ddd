package inmem

import (
	"context"
	"sync"
	"time"

	"github.com/muktihari/order-transaction-ddd/transaction"
	"github.com/shopspring/decimal"
)

type couponRepository struct {
	mu      sync.RWMutex
	coupons map[string]*transaction.Coupon
}

// NewCouponRepository creates new coupon repository in memory
func NewCouponRepository() transaction.CouponRepository {
	return &couponRepository{
		coupons: map[string]*transaction.Coupon{
			"DISCOUNT_$5": {
				Code:     "DISCOUNT_$5",
				Quantity: 100,
				Amount:   decimal.NewFromInt(5),
				Type:     transaction.CouponTypeNominal,
				Begin:    time.Now(),
				End:      time.Now().Add(10 * 24 * time.Hour),
			},
			"DISCOUNT_20%": {
				Code:     "DISCOUNT_20%",
				Quantity: 100,
				Amount:   decimal.NewFromFloat(0.2),
				Type:     transaction.CouponTypePercentage,
				Begin:    time.Now(),
				End:      time.Now().Add(10 * 24 * time.Hour),
			},
		},
	}
}

func (r *couponRepository) FindByCode(ctx context.Context, code string) (*transaction.Coupon, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if val, ok := r.coupons[code]; ok {
		return val, nil
	}
	return nil, transaction.ErrCouponNotFound
}

func (r *couponRepository) Update(ctx context.Context, coupon *transaction.Coupon) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.coupons[coupon.Code] = coupon
	return nil
}
