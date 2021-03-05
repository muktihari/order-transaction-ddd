package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/muktihari/order-transaction-ddd/transaction"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type orderRepository struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

// NewOrderRepository creates new order repository
func NewOrderRepository(client *mongo.Client, db *mongo.Database) transaction.OrderRepository {
	return &orderRepository{client, db, db.Collection("orders")}
}

func (r *orderRepository) FindByID(ctx context.Context, id string) (*transaction.Order, error) {
	sr := r.collection.FindOne(ctx, bson.M{"_id": id})
	if err := sr.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, transaction.ErrOrderNotFound
		}
		return nil, err
	}

	var order transaction.Order
	if err := sr.Decode(&order); err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) Store(ctx context.Context, order *transaction.Order) error {
	order.ID = primitive.NewObjectID().Hex()
	ir, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return err
	}
	order.ID = fmt.Sprintf("%v", ir.InsertedID)

	return nil
}

func (r *orderRepository) Update(ctx context.Context, order *transaction.Order) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": order.ID}, order)
	if err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) FinalizeAndReserveProducts(ctx context.Context, order *transaction.Order) error {
	sess, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer sess.EndSession(nil)

	_, err = sess.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		for _, cartItem := range order.Cart {
			_, err := r.db.Collection("products").UpdateOne(sessCtx,
				bson.M{"_id": cartItem.Product.ID},
				bson.M{"$inc": bson.M{"quantity": cartItem.Quantity * -1}})
			if err != nil {
				return nil, err
			}
			_, err = r.db.Collection("coupons").UpdateOne(sessCtx,
				bson.M{"code": order.Coupon.Code},
				bson.M{"$inc": bson.M{"quantity": -1}})
			if err != nil {
				return nil, err
			}
		}
		_, err = r.collection.ReplaceOne(sessCtx, bson.M{"_id": order.ID}, order)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	if err := sess.CommitTransaction(ctx); err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) CancelAndReleaseProducts(ctx context.Context, order *transaction.Order) error {
	sess, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer sess.EndSession(nil)

	_, err = sess.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		for _, cartItem := range order.Cart {
			_, err := r.db.Collection("products").UpdateOne(sessCtx,
				bson.M{"_id": cartItem.Product.ID},
				bson.M{"$inc": bson.M{"quantity": cartItem.Quantity}})
			if err != nil {
				return nil, err
			}
		}
		_, err = r.collection.ReplaceOne(sessCtx, bson.M{"_id": order.ID}, order)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	if err := sess.CommitTransaction(ctx); err != nil {
		return err
	}

	return nil
}
