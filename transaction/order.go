// Package transaction contains the heart of the domain model.
package transaction

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"
)

var (
	// ErrOrderIsAlreadyFinalized tells that an order can not be changed since it's already finalized.
	ErrOrderIsAlreadyFinalized = errors.New("error order is already finalized")
	// ErrOrderIsAlreadyCompleted tells that an order can not be changed since it's already completed.
	ErrOrderIsAlreadyCompleted = errors.New("error order is already completed")
	// ErrOrderIsAlreadyCanceled tells that an order can not be changed since it's already canceled.
	ErrOrderIsAlreadyCanceled = errors.New("error order is already completed")
	// ErrOrderIsAlreadyShipped tells that an order can not be changed since it's already shipped.
	ErrOrderIsAlreadyShipped = errors.New("error order is already shipped")
	// ErrOrderNotFound tells that order can not be found
	ErrOrderNotFound = errors.New("order not found")
)

// Order is the central class in the domain model
type Order struct {
	ID                   string               `bson:"_id" json:"id"`
	Coupon               Coupon               `bson:"coupon" json:"coupon"`
	Cart                 []CartItem           `bson:"cart" json:"cart"`
	Status               OrderStatus          `bson:"status" json:"status"`
	Price                decimal.Decimal      `bson:"price" json:"price"`
	PriceAfterReduction  decimal.Decimal      `bson:"price_after_reduction" json:"price_after_reduction"`
	Customer             Customer             `bson:"customer" json:"customer"`
	PaymentSpecification PaymentSpecification `bson:"payment_specification" json:"payment_specification"`
	ShippingID           ShippingID           `bson:"shipping_id" json:"shipping_id"`
}

// CartItem represents list of potential bought product with its quantity
type CartItem struct {
	Product  *Product
	Quantity int64
}

// OrderStatus type of status order
type OrderStatus int

const (
	// OrderStatusOpen tells an order is newly created
	OrderStatusOpen OrderStatus = iota + 1
	// OrderStatusSubmitted tells an order is finalized and can not be changed
	OrderStatusSubmitted
	// OrderStatusPaid tells an order is paid
	OrderStatusPaid
	// OrderStatusShipped tells an order is shipped via logistic partner
	OrderStatusShipped
	// OrderStatusCompleted tells an order has received payment proof and completed
	OrderStatusCompleted
	// OrderStatusCancelled tells an order has been canceled, all reserved product quantity are returned
	OrderStatusCancelled
)

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusOpen:
		return "Status Open"
	case OrderStatusSubmitted:
		return "Status Submitted"
	case OrderStatusPaid:
		return "Status Paid"
	case OrderStatusShipped:
		return "Status Shipped"
	case OrderStatusCompleted:
		return "Status Completed"
	case OrderStatusCancelled:
		return "Status Cancelled"
	}
	return ""
}

// NewOrder makes an order
func NewOrder(customer *Customer) *Order {
	return &Order{
		Customer: *customer,
		Cart:     []CartItem{},
		Status:   OrderStatusOpen,
	}
}

// AddProduct add product to the order's ChartItems
func (o *Order) AddProduct(p *Product, quantity int64) error {
	if o.Status != OrderStatusOpen {
		return ErrOrderIsAlreadyFinalized
	}
	for i := range o.Cart {
		if o.Cart[i].Product.ID == p.ID {
			o.Cart[i].Quantity = quantity
			return nil
		}
	}
	o.Cart = append(o.Cart, CartItem{Product: p, Quantity: quantity})
	o.CalculateTotalPrice()
	return nil
}

// ApplyCoupon applies coupon to the Order
func (o *Order) ApplyCoupon(coupon Coupon) error {
	if o.Status != OrderStatusOpen {
		return ErrOrderIsAlreadyFinalized
	}
	if err := coupon.Validate(); err != nil {
		return err
	}
	o.Coupon = coupon
	o.CalculateTotalPrice()
	return nil
}

// ChangeStatusTo changes the status order
func (o *Order) ChangeStatusTo(status OrderStatus) error {
	if o.Status == OrderStatusCompleted {
		return ErrOrderIsAlreadyCompleted
	}
	if o.Status != OrderStatusOpen && status == OrderStatusSubmitted {
		return ErrOrderIsAlreadyFinalized
	}
	if o.Status == status && status == OrderStatusCancelled {
		return ErrOrderIsAlreadyCanceled
	}
	if o.Status == status && status == OrderStatusShipped {
		return ErrOrderIsAlreadyShipped
	}
	o.Status = status
	return nil
}

// CalculateTotalPrice calculates total price of added products and price after reduction if any coupon is applied
func (o *Order) CalculateTotalPrice() {
	for _, cartItem := range o.Cart {
		o.Price = o.Price.Add(cartItem.Product.Price.Mul(decimal.NewFromInt(cartItem.Quantity)))
	}
	if o.Coupon != (Coupon{}) {
		o.PriceAfterReduction = o.Coupon.GetPriceAfterReduction(o.Price)
	}
}

// AllowMakePayment is a policy to an order is allowed payment to be made
func (o *Order) AllowMakePayment() bool {
	return o.Status == OrderStatusSubmitted
}

// SpecifyNewPayment specifies new payment for purchasing the order
func (o *Order) SpecifyNewPayment(ps PaymentSpecification) {
	o.PaymentSpecification = ps
}

// SpecifyShippingID specifies shippingID retrieved from logistics partner
func (o *Order) SpecifyShippingID(shippingID ShippingID) {
	o.ShippingID = shippingID
}

// OrderRepository provides access to orders
type OrderRepository interface {
	FindByID(ctx context.Context, id string) (*Order, error)
	Store(ctx context.Context, order *Order) error
	Update(ctx context.Context, order *Order) error
	FinalizeAndReserveProducts(ctx context.Context, order *Order) error
	CancelAndReleaseProducts(ctx context.Context, order *Order) error
}
