package ordering

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/muktihari/order-transaction-ddd/transaction"
)

var (
	// ErrInvalidArgument occurs when payload argument is invalid
	ErrInvalidArgument = errors.New("invalid argument")
)

// MakeHandler create RestAPI handler
func MakeHandler(s Service) http.Handler {
	r := chi.NewRouter()

	r.Post("/order/make", func(w http.ResponseWriter, r *http.Request) {
		payload := struct {
			CustomerID string `json:"customer_id"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			encodeError(err, w)
			return
		}

		o, err := s.MakeOrder(r.Context(), payload.CustomerID)
		if err != nil {
			encodeError(err, w)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(o); err != nil {
			encodeError(err, w)
			return
		}
	})

	r.Put("/order/{order_id}/addproduct", func(w http.ResponseWriter, r *http.Request) {
		orderID := chi.URLParam(r, "order_id")
		payload := struct {
			ProductID string `json:"product_id"`
			Quantity  int64  `json:"quantity"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			encodeError(err, w)
			return
		}

		if err := s.AddProduct(r.Context(), orderID, payload.ProductID, payload.Quantity); err != nil {
			encodeError(err, w)
			return
		}
	})

	r.Put("/order/{order_id}/applycoupon", func(w http.ResponseWriter, r *http.Request) {
		orderID := chi.URLParam(r, "order_id")
		payload := struct {
			CouponCode string `json:"coupon_code"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			encodeError(err, w)
			return
		}

		if err := s.ApplyCoupon(r.Context(), orderID, payload.CouponCode); err != nil {
			encodeError(err, w)
			return
		}
	})

	r.Post("/order/{order_id}/submit", func(w http.ResponseWriter, r *http.Request) {
		orderID := chi.URLParam(r, "order_id")

		if err := s.SubmitOrder(r.Context(), orderID); err != nil {
			encodeError(err, w)
			return
		}
	})

	r.Post("/order/{order_id}/makepayment", func(w http.ResponseWriter, r *http.Request) {
		orderID := chi.URLParam(r, "order_id")
		var ps transaction.PaymentSpecification
		if err := json.NewDecoder(r.Body).Decode(&ps); err != nil {
			encodeError(err, w)
			return
		}

		if err := s.MakePayment(r.Context(), orderID, ps); err != nil {
			encodeError(err, w)
			return
		}
	})

	r.Get("/order/{order_id}/status", func(w http.ResponseWriter, r *http.Request) {
		orderID := chi.URLParam(r, "order_id")

		status, err := s.CheckOrderStatus(r.Context(), orderID)
		if err != nil {
			encodeError(err, w)
			return
		}
		var response = map[string]interface{}{
			"status": status,
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			encodeError(err, w)
			return
		}
	})

	r.Get("/shipment/{shipping_id}", func(w http.ResponseWriter, r *http.Request) {
		shippingID := transaction.ShippingID(chi.URLParam(r, "shipping_id"))

		status, err := s.CheckShipmentStatus(r.Context(), shippingID)
		if err != nil {
			encodeError(err, w)
			return
		}
		var response = map[string]interface{}{
			"status": status,
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			encodeError(err, w)
			return
		}
	})

	return r
}

func encodeError(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch err {
	case transaction.ErrCustomerNotFound:
		fallthrough
	case transaction.ErrOrderNotFound:
		fallthrough
	case transaction.ErrCouponNotFound:
		fallthrough
	case transaction.ErrProductNotFound:
		w.WriteHeader(http.StatusNotFound)
	case transaction.ErrInvalidCoupon:
		w.WriteHeader(http.StatusConflict)
	case transaction.ErrQuantityExceedProductStock:
		w.WriteHeader(http.StatusConflict)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
