package ordering

import (
	"context"
	"fmt"
	"time"

	"github.com/muktihari/order-transaction-ddd/transaction"
	"github.com/prometheus/client_golang/prometheus"
)

type instrumentingService struct {
	request *prometheus.CounterVec
	latency *prometheus.SummaryVec
	Service
}

// NewInstrumentinService creates new instrumenting service
func NewInstrumentinService(
	request *prometheus.CounterVec,
	latency *prometheus.SummaryVec,
	s Service,
) Service {
	prometheus.MustRegister(request, latency)
	return &instrumentingService{request, latency, s}
}

func (s *instrumentingService) MakeOrder(ctx context.Context, customerID string) (order *transaction.Order, err error) {
	defer func(begin time.Time) {
		s.request.WithLabelValues("make_order", fmt.Sprintf("%t", err != nil)).Inc()
		s.latency.WithLabelValues("make_order", fmt.Sprintf("%t", err != nil)).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.MakeOrder(ctx, customerID)
}

func (s *instrumentingService) AddProduct(ctx context.Context, orderID, productID string, quantity int64) (err error) {
	defer func(begin time.Time) {
		s.request.WithLabelValues("add_product", fmt.Sprintf("%t", err != nil)).Inc()
		s.latency.WithLabelValues("add_product", fmt.Sprintf("%t", err != nil)).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.AddProduct(ctx, orderID, productID, quantity)
}

func (s *instrumentingService) ApplyCoupon(ctx context.Context, orderID, couponCode string) (err error) {
	defer func(begin time.Time) {
		s.request.WithLabelValues("apply_coupon", fmt.Sprintf("%t", err != nil)).Inc()
		s.latency.WithLabelValues("apply_coupon", fmt.Sprintf("%t", err != nil)).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.ApplyCoupon(ctx, orderID, couponCode)
}

func (s *instrumentingService) SubmitOrder(ctx context.Context, orderID string) (err error) {
	defer func(begin time.Time) {
		s.request.WithLabelValues("submit_order", fmt.Sprintf("%t", err != nil)).Inc()
		s.latency.WithLabelValues("submit_order", fmt.Sprintf("%t", err != nil)).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.SubmitOrder(ctx, orderID)
}

func (s *instrumentingService) MakePayment(ctx context.Context, orderID string, ps transaction.PaymentSpecification) (err error) {
	defer func(begin time.Time) {
		s.request.WithLabelValues("make_payment", fmt.Sprintf("%t", err != nil)).Inc()
		s.latency.WithLabelValues("make_payment", fmt.Sprintf("%t", err != nil)).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.MakePayment(ctx, orderID, ps)
}

func (s *instrumentingService) CheckOrderStatus(ctx context.Context, orderID string) (status transaction.OrderStatus, err error) {
	defer func(begin time.Time) {
		s.request.WithLabelValues("check_order_status", fmt.Sprintf("%t", err != nil)).Inc()
		s.latency.WithLabelValues("check_order_status", fmt.Sprintf("%t", err != nil)).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.CheckOrderStatus(ctx, orderID)
}

func (s *instrumentingService) CheckShipmentStatus(ctx context.Context, shippingID transaction.ShippingID) (status transaction.ShipmentStatus, err error) {
	defer func(begin time.Time) {
		s.request.WithLabelValues("check_shipment_status", fmt.Sprintf("%t", err != nil)).Inc()
		s.latency.WithLabelValues("check_shipment_status", fmt.Sprintf("%t", err != nil)).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.CheckShipmentStatus(ctx, shippingID)
}
