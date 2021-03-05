package ordering

import (
	"context"

	"github.com/muktihari/order-transaction-ddd/transaction"
)

// Service is the interface that provides ordering methods.
type Service interface {
	// MakeOrder creates new open order for the customer
	MakeOrder(ctx context.Context, customerID string) (*transaction.Order, error)
	// AddProduct adds product with set quantity to the order
	AddProduct(ctx context.Context, orderID, productID string, quantity int64) error
	// ApplyCoupon applies coupon to the order
	ApplyCoupon(ctx context.Context, orderID, couponCode string) error
	// SubmitOrder reserves added products and its quantity and finalize order
	SubmitOrder(ctx context.Context, orderID string) error
	// MakePayment makes payment for submitted order
	MakePayment(ctx context.Context, orderID string, ps transaction.PaymentSpecification) error
	// CheckOrderStatus checks status order
	CheckOrderStatus(ctx context.Context, orderID string) (transaction.OrderStatus, error)
	// CheckShipmentStatus checks shipment status
	CheckShipmentStatus(ctx context.Context, shippingID transaction.ShippingID) (transaction.ShipmentStatus, error)
}

type service struct {
	orders    transaction.OrderRepository
	customers transaction.CustomerRepository
	products  transaction.ProductRepository
	coupons   transaction.CouponRepository
	logistics transaction.LogisticsPartner
}

// NewService creates a ordering service with necessary dependencies
func NewService(
	orders transaction.OrderRepository,
	customers transaction.CustomerRepository,
	products transaction.ProductRepository,
	coupons transaction.CouponRepository,
	logistics transaction.LogisticsPartner,
) Service {
	return &service{
		orders:    orders,
		customers: customers,
		products:  products,
		coupons:   coupons,
		logistics: logistics,
	}
}

func (s *service) MakeOrder(ctx context.Context, customerID string) (*transaction.Order, error) {
	c, err := s.customers.FindByID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	o := transaction.NewOrder(c)

	if err := s.orders.Store(ctx, o); err != nil {
		return nil, err
	}

	return o, nil
}

func (s *service) AddProduct(ctx context.Context, orderID string, productID string, quantity int64) error {
	o, err := s.orders.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	p, err := s.products.FindByID(ctx, productID)
	if err != nil {
		return err
	}

	if err := p.TryReserveQuantity(quantity); err != nil {
		return err
	}

	if err := o.AddProduct(p, quantity); err != nil {
		return err
	}

	if err := s.orders.Update(ctx, o); err != nil {
		return err
	}

	return nil
}

func (s *service) ApplyCoupon(ctx context.Context, orderID, couponCode string) error {
	o, err := s.orders.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	c, err := s.coupons.FindByCode(ctx, couponCode)
	if err != nil {
		return err
	}

	if err := o.ApplyCoupon(*c); err != nil {
		return err
	}

	if err := s.orders.Update(ctx, o); err != nil {
		return err
	}

	return nil
}

func (s *service) SubmitOrder(ctx context.Context, orderID string) error {
	o, err := s.orders.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if err := o.ChangeStatusTo(transaction.OrderStatusSubmitted); err != nil {
		return err
	}

	if err := s.orders.FinalizeAndReserveProducts(ctx, o); err != nil {
		return err
	}

	return nil
}

func (s *service) MakePayment(ctx context.Context, orderID string, ps transaction.PaymentSpecification) error {
	if err := ps.Validate(); err != nil {
		return err
	}

	o, err := s.orders.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	o.SpecifyNewPayment(ps)
	if err := o.ChangeStatusTo(transaction.OrderStatusPaid); err != nil {
		return err
	}

	if err := s.orders.Update(ctx, o); err != nil {
		return err
	}

	return nil
}

func (s *service) CheckOrderStatus(ctx context.Context, orderID string) (transaction.OrderStatus, error) {
	o, err := s.orders.FindByID(ctx, orderID)
	if err != nil {
		return 0, err
	}
	return o.Status, nil
}

func (s *service) CheckShipmentStatus(ctx context.Context, shippingID transaction.ShippingID) (transaction.ShipmentStatus, error) {
	return s.logistics.CheckShipmentStatus(ctx, shippingID)
}
