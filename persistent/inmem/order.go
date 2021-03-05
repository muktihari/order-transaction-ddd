package inmem

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/muktihari/order-transaction-ddd/transaction"
	"github.com/shopspring/decimal"
)

type orderRepository struct {
	mu       sync.RWMutex
	orders   map[string]*transaction.Order
	coupons  transaction.CouponRepository
	products transaction.ProductRepository
}

// NewOrderRepository creates new order repository in memory
func NewOrderRepository(coupons transaction.CouponRepository, products transaction.ProductRepository) transaction.OrderRepository {
	return &orderRepository{
		orders: map[string]*transaction.Order{
			"ORDER_OPEN": {
				ID: "ORDER_OPEN",
				Customer: transaction.Customer{
					ID: "CUSTOMER1", Name: "Hari", PhoneNumber: "+62-12345", Email: "example@email.com", Address: "No, Street, City, Indonesia",
				},
				Cart:   []transaction.CartItem{},
				Status: transaction.OrderStatusOpen,
			},
			"ORDER_WITH_PRODUCT": {
				ID: "ORDER_WITH_PRODUCT",
				Customer: transaction.Customer{
					ID: "CUSTOMER1", Name: "Hari", PhoneNumber: "+62-12345", Email: "example@email.com", Address: "No, Street, City, Indonesia",
				},
				Cart: []transaction.CartItem{
					{
						Product:  &transaction.Product{ID: "PRODUCT1", Name: "Sony Xperia 10", Price: decimal.NewFromInt(500), Quantity: 200},
						Quantity: 5,
					},
				},
				Status: transaction.OrderStatusOpen,
			},
			"ORDER_WITH_PRODUCT_AND_COUPON": {
				ID: "ORDER_WITH_PRODUCT_AND_COUPON",
				Customer: transaction.Customer{
					ID: "CUSTOMER1", Name: "Hari", PhoneNumber: "+62-12345", Email: "example@email.com", Address: "No, Street, City, Indonesia",
				},
				Coupon: transaction.Coupon{
					Code:     "DISCOUNT_20%",
					Quantity: 100,
					Amount:   decimal.NewFromFloat(0.2),
					Type:     transaction.CouponTypePercentage,
					Begin:    time.Now(),                          // will always valid
					End:      time.Now().Add(10 * 24 * time.Hour), // will always valid
				},
				Cart: []transaction.CartItem{
					{
						Product:  &transaction.Product{ID: "PRODUCT1", Name: "Sony Xperia 10", Price: decimal.NewFromInt(500), Quantity: 200},
						Quantity: 5,
					},
				},
				Price:               decimal.NewFromInt(500 * 5),
				PriceAfterReduction: decimal.NewFromInt(500 * 5).Sub(decimal.NewFromInt(500 * 5).Mul(decimal.NewFromFloat(0.2))),
				Status:              transaction.OrderStatusOpen,
			},
		},
		coupons:  coupons,
		products: products,
	}
}

func (r *orderRepository) FindByID(ctx context.Context, id string) (*transaction.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if val, ok := r.orders[id]; ok {
		return val, nil
	}
	return nil, transaction.ErrOrderNotFound
}

func (r *orderRepository) Store(ctx context.Context, order *transaction.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	order.ID = uuid.NewString()
	r.orders[order.ID] = order
	return nil
}

func (r *orderRepository) Update(ctx context.Context, order *transaction.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
	return nil
}

func (r *orderRepository) FinalizeAndReserveProducts(ctx context.Context, order *transaction.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// assume it's transactional
	coupon, err := r.coupons.FindByCode(ctx, order.Coupon.Code)
	if err != nil {
		return err
	}
	coupon.Quantity--
	if err := r.coupons.Update(ctx, coupon); err != nil {
		return err
	}

	for _, cartItem := range order.Cart {
		p, err := r.products.FindByID(ctx, cartItem.Product.ID)
		if err != nil {
			return err
		}
		if err := p.TryReserveQuantity(cartItem.Quantity); err != nil {
			return err
		}
		p.ReserveQuantity(cartItem.Quantity)
		if err := r.products.Update(ctx, p); err != nil {
			return err
		}

	}

	r.orders[order.ID] = order

	return nil
}

func (r *orderRepository) CancelAndReleaseProducts(ctx context.Context, order *transaction.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// assume it's transactional
	coupon, err := r.coupons.FindByCode(ctx, order.Coupon.Code)
	if err != nil {
		return err
	}
	coupon.Quantity++
	if err := r.coupons.Update(ctx, coupon); err != nil {
		return err
	}

	for _, cartItem := range order.Cart {
		p, err := r.products.FindByID(ctx, cartItem.Product.ID)
		if err != nil {
			return err
		}
		p.RollbackQuantity(cartItem.Quantity)
		if err := r.products.Update(ctx, p); err != nil {
			return err
		}
	}
	r.orders[order.ID] = order

	return nil
}
