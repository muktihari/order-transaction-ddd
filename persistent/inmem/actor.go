package inmem

import (
	"context"
	"sync"

	"github.com/muktihari/order-transaction-ddd/transaction"
	"github.com/sirupsen/logrus"
)

type customerRepository struct {
	mu        sync.RWMutex
	customers map[string]*transaction.Customer
}

// NewCustomerRepository creates new customer repository in memory
func NewCustomerRepository() transaction.CustomerRepository {
	return &customerRepository{
		customers: map[string]*transaction.Customer{
			"CUSTOMER1": {ID: "CUSTOMER1", Name: "Hari", PhoneNumber: "+62-12345", Email: "example@email.com", Address: "No, Street, City, Indonesia"},
		},
	}
}

func (r *customerRepository) FindByID(ctx context.Context, id string) (*transaction.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	logrus.Warnf("id: %s", id)
	if val, ok := r.customers[id]; ok {
		return val, nil
	}
	return nil, transaction.ErrCustomerNotFound
}

type adminRepository struct {
	mu     sync.RWMutex
	admins map[string]*transaction.Admin
}

// NewAdminRepository creates new admin repository in memory
func NewAdminRepository() transaction.CustomerRepository {
	return &customerRepository{
		customers: map[string]*transaction.Customer{
			"ADMIN1": {ID: "ADMIN1", Name: "Mukti", PhoneNumber: "+62-67890", Email: "example@email.com", Address: "No, Street, City, Indonesia"},
		},
	}
}

func (r *adminRepository) FindByID(ctx context.Context, id string) (*transaction.Admin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if val, ok := r.admins[id]; ok {
		return val, nil
	}
	return nil, transaction.ErrCustomerNotFound
}
