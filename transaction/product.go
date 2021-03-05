package transaction

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"
)

var (
	// ErrQuantityExceedProductStock occurs when customer trying to add product more than product available
	ErrQuantityExceedProductStock = errors.New("error quantity exceed product's stock")
	// ErrProductNotFound tells that product can not be found
	ErrProductNotFound = errors.New("product not found")
)

// Product represent item to sell
type Product struct {
	ID       string          `bson:"_id" json:"id"`
	Name     string          `bson:"name" json:"name"`
	Price    decimal.Decimal `bson:"price" json:"price"`
	Quantity int64           `bson:"quantity" json:"quantity"`
}

// TryReserveQuantity checks whether product's quantity can be reserved
func (p *Product) TryReserveQuantity(quantity int64) error {
	if p.Quantity < quantity {
		return ErrQuantityExceedProductStock
	}
	return nil
}

// ReserveQuantity subtracts quantity from product
func (p *Product) ReserveQuantity(quantity int64) {
	p.Quantity -= quantity
}

// RollbackQuantity adds quantity from order to product
func (p *Product) RollbackQuantity(quantity int64) {
	p.Quantity += quantity
}

// ProductRepository represent
type ProductRepository interface {
	FindByID(ctx context.Context, id string) (*Product, error)
	FindAll(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, product *Product) error
}
