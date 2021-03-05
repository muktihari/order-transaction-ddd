package handling

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

// NewInstrumentingService create new logging service
func NewInstrumentingService(
	request *prometheus.CounterVec,
	latency *prometheus.SummaryVec,
	s Service,
) Service {
	prometheus.MustRegister(request, latency)
	return &instrumentingService{request, latency, s}
}

func (s *instrumentingService) ViewOrder(ctx context.Context, orderID string) (order *transaction.Order, err error) {
	defer func(begin time.Time) {
		s.request.WithLabelValues("view_order", fmt.Sprintf("%t", err != nil)).Inc()
		s.latency.WithLabelValues("view_order", fmt.Sprintf("%t", err != nil)).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.ViewOrder(ctx, orderID)
}

func (s *instrumentingService) CancelOrder(ctx context.Context, orderID string) (err error) {
	defer func(begin time.Time) {
		s.request.WithLabelValues("cancel_order", fmt.Sprintf("%t", err != nil)).Inc()
		s.latency.WithLabelValues("cancel_order", fmt.Sprintf("%t", err != nil)).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.CancelOrder(ctx, orderID)
}

func (s *instrumentingService) ShipOrderToLogisticsPartner(ctx context.Context, orderID string) (shippingID transaction.ShippingID, err error) {
	defer func(begin time.Time) {
		s.request.WithLabelValues("ship_order_to_logistics_partner", fmt.Sprintf("%t", err != nil)).Inc()
		s.latency.WithLabelValues("ship_order_to_logistics_partner", fmt.Sprintf("%t", err != nil)).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.ShipOrderToLogisticsPartner(ctx, orderID)
}
