package mongodb

import (
	"context"
	"errors"

	"github.com/muktihari/order-transaction-ddd/transaction"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type productRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewProductRepository creates new product repository in memory
func NewProductRepository(db *mongo.Database) transaction.ProductRepository {
	return &productRepository{db, db.Collection("products")}
}

func (r *productRepository) FindByID(ctx context.Context, id string) (*transaction.Product, error) {
	sr := r.collection.FindOne(ctx, bson.M{"_id": id})
	if err := sr.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, transaction.ErrCouponNotFound
		}
		return nil, err
	}

	var product transaction.Product
	if err := sr.Decode(&product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindAll(ctx context.Context) ([]transaction.Product, error) {
	cur, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, transaction.ErrCouponNotFound
		}
		return nil, err
	}
	defer cur.Close(nil)

	var products []transaction.Product
	if err := cur.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) Update(ctx context.Context, product *transaction.Product) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": product.ID}, product)
	if err != nil {
		return err
	}

	return nil
}
