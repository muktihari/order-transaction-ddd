package mongodb

import (
	"context"
	"errors"

	"github.com/muktihari/order-transaction-ddd/transaction"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type couponRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewCouponRepository creates new coupon repository
func NewCouponRepository(db *mongo.Database) transaction.CouponRepository {
	return &couponRepository{db, db.Collection("coupons")}
}

func (r *couponRepository) FindByCode(ctx context.Context, code string) (*transaction.Coupon, error) {
	sr := r.collection.FindOne(ctx, bson.M{"code": code})
	if err := sr.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, transaction.ErrCouponNotFound
		}
		return nil, err
	}

	var coupon transaction.Coupon
	if err := sr.Decode(&coupon); err != nil {
		return nil, err
	}

	return &coupon, nil
}

func (r *couponRepository) Update(ctx context.Context, coupon *transaction.Coupon) error {
	sr := r.collection.FindOneAndReplace(ctx, bson.M{"code": coupon.Code}, coupon)
	if err := sr.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return transaction.ErrCouponNotFound
		}
		return err
	}

	return nil
}
