package inmem

import (
	"context"
	"sync"

	"github.com/muktihari/order-transaction-ddd/transaction"
	"github.com/shopspring/decimal"
)

type productRepository struct {
	mu       sync.RWMutex
	products map[string]*transaction.Product
}

// NewProductRepository creates new product repository in memory
func NewProductRepository() transaction.ProductRepository {
	return &productRepository{
		products: map[string]*transaction.Product{
			"PRODUCT1": {ID: "PRODUCT1", Name: "Sony Xperia 10", Price: decimal.NewFromInt(500), Quantity: 200},
			"PRODUCT2": {ID: "PRODUCT2", Name: "Ultramilk 1L", Price: decimal.NewFromInt(5), Quantity: 2000},
		},
	}
}

func (r *productRepository) FindByID(ctx context.Context, id string) (*transaction.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if val, ok := r.products[id]; ok {
		return val, nil
	}
	return nil, transaction.ErrProductNotFound
}

func (r *productRepository) FindAll(ctx context.Context) ([]transaction.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	products := make([]transaction.Product, 0)
	for _, p := range r.products {
		products = append(products, *p)
	}
	return products, nil
}

func (r *productRepository) Update(ctx context.Context, product *transaction.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.products[product.ID] = product
	return nil
}
