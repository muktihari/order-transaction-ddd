package transaction

import (
	"context"
	"errors"
)

var (
	// ErrLogisticsRegister occurs when trying to register order to logistics partner
	ErrLogisticsRegister = errors.New("error register logistics")
	// ErrLogisticsCheckShipment occurs when trying to check shipment status from logistics partner
	ErrLogisticsCheckShipment = errors.New("error check shipment logistics")
)

// ShipmentStatus type shipment status
type ShipmentStatus int

const (
	// ShipmentStatusShipped tells an order is already shipped
	ShipmentStatusShipped ShipmentStatus = iota + 1
	// ShipmentStatusOnDelivery tells an order is on the way to customer's address
	ShipmentStatusOnDelivery
	// ShipmentStatusDelived tells an order has been received by customer
	ShipmentStatusDelived
)

func (s ShipmentStatus) String() string {
	switch s {
	case ShipmentStatusShipped:
		return "Status Shipped"
	case ShipmentStatusOnDelivery:
		return "Status OnDelivery"
	case ShipmentStatusDelived:
		return "Status Delived"
	}
	return ""
}

// ShippingID type of shipping id
type ShippingID string

// LogisticsPartner provides access to logistics partner API
type LogisticsPartner interface {
	RegisterShipment(ctx context.Context, orderID string) (ShippingID, error)
	CheckShipmentStatus(ctx context.Context, shippingID ShippingID) (ShipmentStatus, error)
}
