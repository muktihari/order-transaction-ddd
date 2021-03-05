package mongodb

import (
	"context"
	"errors"

	"github.com/muktihari/order-transaction-ddd/transaction"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type customerRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewCustomerRepository creates new customer repository
func NewCustomerRepository(db *mongo.Database) transaction.CustomerRepository {
	return &customerRepository{db, db.Collection("customers")}
}

func (r *customerRepository) FindByID(ctx context.Context, id string) (*transaction.Customer, error) {
	sr := r.collection.FindOne(ctx, bson.M{"_id": id})
	if err := sr.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, transaction.ErrCustomerNotFound
		}
		return nil, err
	}

	var customer transaction.Customer
	if err := sr.Decode(&customer); err != nil {
		return nil, err
	}

	return &customer, nil
}

type adminRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewAdminRepository creates new admin repository
func NewAdminRepository(db *mongo.Database) transaction.AdminRepository {
	return &adminRepository{db, db.Collection("customers")}
}

func (r *adminRepository) FindByID(ctx context.Context, id string) (*transaction.Admin, error) {
	sr := r.collection.FindOne(ctx, bson.M{"_id": id})
	if err := sr.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, transaction.ErrAdminNotFound
		}
		return nil, err
	}

	var admin transaction.Admin
	if err := sr.Decode(&admin); err != nil {
		return nil, err
	}

	return &admin, nil
}
