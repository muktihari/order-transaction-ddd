package inmem

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/muktihari/order-transaction-ddd/transaction"
)

type shipment struct {
	OrderID string
	Status  transaction.ShipmentStatus
}

type logisticsPartner struct {
	mu        sync.RWMutex
	logistics map[transaction.ShippingID]shipment
}

// NewLogisticsParner create new logistics partner in memory
func NewLogisticsParner() transaction.LogisticsPartner {
	return &logisticsPartner{
		logistics: make(map[transaction.ShippingID]shipment),
	}
}

func (r *logisticsPartner) RegisterShipment(ctx context.Context, orderID string) (transaction.ShippingID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, val := range r.logistics {
		if val.OrderID == orderID {
			return "", transaction.ErrLogisticsRegister
		}
	}

	shippingID := uuid.NewString()
	r.logistics[transaction.ShippingID(shippingID)] = shipment{OrderID: orderID, Status: transaction.ShipmentStatusShipped}
	return transaction.ShippingID(shippingID), nil
}

func (r *logisticsPartner) CheckShipmentStatus(ctx context.Context, shippingID transaction.ShippingID) (transaction.ShipmentStatus, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if val, ok := r.logistics[shippingID]; ok {
		return val.Status, nil
	}
	return 0, transaction.ErrLogisticsCheckShipment
}
