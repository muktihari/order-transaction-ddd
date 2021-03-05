package ordering

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

func (s *loggingService) MakeOrder(ctx context.Context, customerID string) (order *transaction.Order, err error) {
	defer func(begin time.Time) {
		var orderID string
		if order != nil {
			orderID = order.ID
		}
		s.log.WithFields(log.Fields{
			"method":      "make_order",
			"customer_id": customerID,
			"order_id":    orderID,
			"took":        time.Since(begin),
			"err":         err,
		}).Println()
	}(time.Now())
	return s.Service.MakeOrder(ctx, customerID)
}

func (s *loggingService) AddProduct(ctx context.Context, orderID, productID string, quantity int64) (err error) {
	defer func(begin time.Time) {
		s.log.WithFields(log.Fields{
			"method":     "add_product",
			"order_id":   orderID,
			"product_id": productID,
			"quantity":   quantity,
			"took":       time.Since(begin),
			"err":        err,
		}).Println()
	}(time.Now())
	return s.Service.AddProduct(ctx, orderID, productID, quantity)
}

func (s *loggingService) ApplyCoupon(ctx context.Context, orderID, couponCode string) (err error) {
	defer func(begin time.Time) {
		s.log.WithFields(log.Fields{
			"method":      "apply_coupon",
			"order_id":    orderID,
			"coupon_code": couponCode,
			"took":        time.Since(begin),
			"err":         err,
		}).Println()
	}(time.Now())
	return s.Service.ApplyCoupon(ctx, orderID, couponCode)
}

func (s *loggingService) SubmitOrder(ctx context.Context, orderID string) (err error) {
	defer func(begin time.Time) {
		s.log.WithFields(log.Fields{
			"method":   "submit_order",
			"order_id": orderID,
			"took":     time.Since(begin),
			"err":      err,
		}).Println()
	}(time.Now())
	return s.Service.SubmitOrder(ctx, orderID)
}

func (s *loggingService) MakePayment(ctx context.Context, orderID string, ps transaction.PaymentSpecification) (err error) {
	defer func(begin time.Time) {
		s.log.WithFields(log.Fields{
			"method":                "make_payment",
			"order_id":              orderID,
			"payment_specification": ps,
			"took":                  time.Since(begin),
			"err":                   err,
		}).Println()
	}(time.Now())
	return s.Service.MakePayment(ctx, orderID, ps)
}

func (s *loggingService) CheckOrderStatus(ctx context.Context, orderID string) (status transaction.OrderStatus, err error) {
	defer func(begin time.Time) {
		s.log.WithFields(log.Fields{
			"method":   "check_order_status",
			"order_id": orderID,
			"took":     time.Since(begin),
			"status":   status,
			"err":      err,
		}).Println()
	}(time.Now())
	return s.Service.CheckOrderStatus(ctx, orderID)
}

func (s *loggingService) CheckShipmentStatus(ctx context.Context, shippingID transaction.ShippingID) (status transaction.ShipmentStatus, err error) {
	defer func(begin time.Time) {
		s.log.WithFields(log.Fields{
			"method":      "check_shipment_status",
			"shipping_id": shippingID,
			"took":        time.Since(begin),
			"status":      status,
			"err":         err,
		}).Println()
	}(time.Now())
	return s.Service.CheckShipmentStatus(ctx, shippingID)
}
