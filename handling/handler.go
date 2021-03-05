package handling

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/muktihari/order-transaction-ddd/transaction"
)

// MakeHandler create RestAPI handler
func MakeHandler(s Service) http.Handler {
	r := chi.NewRouter()

	r.Get("/order/{order_id}/view", func(w http.ResponseWriter, r *http.Request) {
		orderID := chi.URLParam(r, "order_id")
		o, err := s.ViewOrder(r.Context(), orderID)
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

	r.Post("/order/{order_id}/cancel", func(w http.ResponseWriter, r *http.Request) {
		orderID := chi.URLParam(r, "order_id")
		if err := s.CancelOrder(r.Context(), orderID); err != nil {
			encodeError(err, w)
			return
		}
	})

	r.Post("/order/{order_id}/ship", func(w http.ResponseWriter, r *http.Request) {
		orderID := chi.URLParam(r, "order_id")
		shippingID, err := s.ShipOrderToLogisticsPartner(r.Context(), orderID)
		if err != nil {
			encodeError(err, w)
			return
		}

		var response = map[string]interface{}{
			"shipping_id": shippingID,
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
	case transaction.ErrOrderNotFound:
		w.WriteHeader(http.StatusNotFound)
	case transaction.ErrOrderIsAlreadyFinalized:
		w.WriteHeader(http.StatusConflict)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
