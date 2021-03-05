package postgresql

import (
	"context"
	"database/sql"

	"github.com/muktihari/order-transaction-ddd/transaction"
)

type customerRepository struct {
	db *sql.DB
}

// NewCustomerRepository creates new customer repository
func NewCustomerRepository(db *sql.DB) transaction.CustomerRepository {
	return &customerRepository{db}
}

func (r *customerRepository) FindByID(ctx context.Context, id string) (*transaction.Customer, error) {
	stmt, err := r.db.PrepareContext(ctx, "select * from customers where id = $1")
	if err != nil {
		return nil, err
	}
	sqlRows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	var customer transaction.Customer
	if sqlRows.Next() == false {
		return nil, transaction.ErrCustomerNotFound
	}
	if err := Map(sqlRows, &customer); err != nil {
		return nil, err
	}
	sqlRows.Close()

	return &customer, nil
}

type adminRepository struct {
	db *sql.DB
}

func (r *adminRepository) FindByID(ctx context.Context, id string) (*transaction.Admin, error) {
	stmt, err := r.db.PrepareContext(ctx, "select * from admins where id = $1")
	if err != nil {
		return nil, err
	}
	sqlRows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	var admin transaction.Admin
	if sqlRows.Next() == false {
		return nil, transaction.ErrAdminNotFound
	}
	if err := Map(sqlRows, &admin); err != nil {
		return nil, err
	}
	sqlRows.Close()

	return &admin, nil
}
