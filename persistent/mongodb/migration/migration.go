package migration

import (
	"context"
	"time"

	"github.com/muktihari/order-transaction-ddd/transaction"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MigratePredefinedData drop transaction-order database and create new one with predefined data such as actor, product, coupon
func MigratePredefinedData(ctx context.Context, client *mongo.Client) (err error) {
	if err := client.Database("transaction-order").Drop(ctx); err != nil {
		return err
	}

	db := client.Database("transaction-order")
	_, err = db.Collection("admins").InsertMany(ctx, []interface{}{
		transaction.Admin{ID: primitive.NewObjectID().Hex(), Name: "Mukti"},
	})
	if err != nil {
		return err
	}

	_, err = db.Collection("customers").InsertMany(ctx, []interface{}{
		transaction.Customer{ID: primitive.NewObjectID().Hex(), Name: "Hari", PhoneNumber: "+62-12345", Email: "example@email.com", Address: "No, Street, City, Indonesia"},
	})
	if err != nil {
		return err
	}

	_, err = db.Collection("products").InsertMany(ctx, []interface{}{
		transaction.Product{ID: primitive.NewObjectID().Hex(), Name: "Sony Xperia 10", Price: decimal.NewFromInt(500), Quantity: 100},
		transaction.Product{ID: primitive.NewObjectID().Hex(), Name: "Ultramilk 1 Liter", Price: decimal.NewFromInt(5), Quantity: 1000},
	})
	if err != nil {
		return err
	}

	_, err = db.Collection("coupons").InsertMany(ctx, []interface{}{
		transaction.Coupon{Code: "DISCOUNT_5$", Type: transaction.CouponTypeNominal, Amount: decimal.NewFromInt(5), Quantity: 100, Begin: time.Now(), End: time.Now().Add(10 * 24 * time.Hour)},
		transaction.Coupon{Code: "DISCOUNT_20%", Type: transaction.CouponTypeNominal, Amount: decimal.NewFromFloat(0.2), Quantity: 100, Begin: time.Now(), End: time.Now().Add(10 * 24 * time.Hour)},
	})

	return nil
}
