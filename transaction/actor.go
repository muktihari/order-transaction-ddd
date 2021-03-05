package transaction

import (
	"context"
	"errors"
)

var (
	// ErrCustomerNotFound tells that product can not be found
	ErrCustomerNotFound = errors.New("customer not found")
	// ErrAdminNotFound tells that product can not be found
	ErrAdminNotFound = errors.New("admin not found")
)

// Admin represents the person in charge to administrate the shop
type Admin struct {
	ID   string `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
}

// Customer represents customer who want to buy products from the shop
type Customer struct {
	ID          string `bson:"_id" json:"id"`
	Name        string `bson:"name" json:"name"`
	PhoneNumber string `bson:"phone_number" json:"phone_number"`
	Email       string `bson:"email" json:"email"`
	Address     string `bson:"address" json:"address"`
}

// CustomerRepository provides access to customers
type CustomerRepository interface {
	FindByID(ctx context.Context, id string) (*Customer, error)
}

// AdminRepository provides access to admins
type AdminRepository interface {
	FindByID(ctx context.Context, id string) (*Admin, error)
}
