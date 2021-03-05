package handling

import (
	"context"
	"time"

	"github.com/muktihari/order-transaction-ddd/transaction"
	log "github.com/sirupsen/logrus"
)

type loggingService struct {
	log *log.Logger
	Service
}

// NewLoggingService create new logging service
func NewLoggingService(log *log.Logger, s Service) Service {
	return &loggingService{log, s}
}

func (s *loggingService) ViewOrder(ctx context.Context, orderID string) (order *transaction.Order, err error) {
	defer func(begin time.Time) {
		s.log.WithFields(log.Fields{
			"method":   "view_order",
			"order_id": orderID,
			"took":     time.Since(begin),
			"err":      err,
		}).Println()
	}(time.Now())
	return s.Service.ViewOrder(ctx, orderID)
}

func (s *loggingService) CancelOrder(ctx context.Context, orderID string) (err error) {
	defer func(begin time.Time) {
		s.log.WithFields(log.Fields{
			"method":   "cancel_order",
			"order_id": orderID,
			"took":     time.Since(begin),
			"err":      err,
		}).Println()
	}(time.Now())
	return s.Service.CancelOrder(ctx, orderID)
}

func (s *loggingService) ShipOrderToLogisticsPartner(ctx context.Context, orderID string) (shippingID transaction.ShippingID, err error) {
	defer func(begin time.Time) {
		s.log.WithFields(log.Fields{
			"method":      "ship_order_to_logistics_partner",
			"order_id":    orderID,
			"shipping_id": shippingID,
			"took":        time.Since(begin),
			"err":         err,
		}).Println()
	}(time.Now())
	return s.Service.ShipOrderToLogisticsPartner(ctx, orderID)
}
