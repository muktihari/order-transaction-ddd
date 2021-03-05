package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/muktihari/order-transaction-ddd/transaction"
)

type couponRepository struct {
	db    *sql.DB
	table string
}

// NewCouponRepository creates new coupon repository
func NewCouponRepository(db *sql.DB) transaction.CouponRepository {
	return &couponRepository{db, "coupons"}
}

func (r *couponRepository) FindByCode(ctx context.Context, code string) (*transaction.Coupon, error) {
	stmt, err := r.db.PrepareContext(ctx, "select * from "+r.table+" where code = $1")
	if err != nil {
		return nil, err
	}
	sqlRows, err := stmt.Query(code)
	if err != nil {
		return nil, err
	}

	var coupon transaction.Coupon
	if sqlRows.Next() == false {
		return nil, transaction.ErrCouponNotFound
	}
	if err := Map(sqlRows, &coupon); err != nil {
		return nil, err
	}
	sqlRows.Close()

	return &coupon, nil
}

func (r *couponRepository) Update(ctx context.Context, coupon *transaction.Coupon) error {
	c, err := r.FindByCode(ctx, coupon.Code)
	if err != nil {
		return err
	}

	diff, err := KeyValsDiff(c, coupon)
	if err != nil {
		return err
	}

	var (
		columnValues []string
		args         []interface{}
	)
	var i int
	for key, val := range diff {
		columnValues = append(columnValues, fmt.Sprintf("%v = $%d", key, i+1))
		args = append(args, val)
		i++
	}

	query := fmt.Sprintf("update %s set %s where code = %s", r.table, strings.Join(columnValues, ", "), c.Code)
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(ctx, args...); err != nil {
		return err
	}
	return nil

}
