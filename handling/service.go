// Package handling contains process of handling order after order is being submitted
// by customer assigned to admin
package handling

import (
	"context"

	"github.com/muktihari/order-transaction-ddd/transaction"
)

// Service is the interface that provides handling method.
type Service interface {
	// ViewOrder views order details
	ViewOrder(ctx context.Context, orderID string) (*transaction.Order, error)
	// CancelOrder cancels order
	CancelOrder(ctx context.Context, orderID string) error
	// ShipOrderToLogisticsPartner ships the order to logistics partner. It will update shippingID on order
	ShipOrderToLogisticsPartner(ctx context.Context, orderID string) (transaction.ShippingID, error)
}

type service struct {
	orders    transaction.OrderRepository
	products  transaction.ProductRepository
	logistics transaction.LogisticsPartner
}

// NewService creates a handling service with necessary dependencies
func NewService(
	orders transaction.OrderRepository,
	products transaction.ProductRepository,
	logistics transaction.LogisticsPartner,
) Service {
	return &service{
		orders:    orders,
		products:  products,
		logistics: logistics,
	}
}

func (s *service) ViewOrder(ctx context.Context, orderID string) (*transaction.Order, error) {
	return s.orders.FindByID(ctx, orderID)
}

func (s *service) CancelOrder(ctx context.Context, orderID string) error {
	o, err := s.orders.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if err := o.ChangeStatusTo(transaction.OrderStatusCancelled); err != nil {
		return err
	}

	if err := s.orders.CancelAndReleaseProducts(ctx, o); err != nil {
		return err
	}

	return nil
}

func (s *service) ShipOrderToLogisticsPartner(ctx context.Context, orderID string) (transaction.ShippingID, error) {
	var shippingID transaction.ShippingID

	o, err := s.orders.FindByID(ctx, orderID)
	if err != nil {
		return shippingID, err
	}

	if err := o.ChangeStatusTo(transaction.OrderStatusShipped); err != nil {
		return shippingID, err
	}

	shippingID, err = s.logistics.RegisterShipment(ctx, orderID)
	if err != nil {
		return shippingID, err
	}

	o, err = s.orders.FindByID(ctx, orderID)
	if err != nil {
		return shippingID, err
	}

	o.SpecifyShippingID(shippingID)

	if err := s.orders.Update(ctx, o); err != nil {
		return shippingID, err
	}

	return shippingID, nil
}
