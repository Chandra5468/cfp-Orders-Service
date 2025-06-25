package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/Chandra5468/cfp-Products-Service/internal/types"
	"github.com/Chandra5468/cfp-Products-Service/internal/utils/responses"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Middleware for RBAC, Logging, checking validations etc...

func CorsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func OrderOwnerShipMiddleware(psqlStore types.OrdersStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orderIdStr := chi.URLParam(r, "id")

		if orderIdStr == "" {
			responses.WriteJson(w, http.StatusBadRequest, "Please mention the orderid")
			return
		}

		orderId, err := uuid.Parse(orderIdStr)
		if err != nil {
			responses.WriteJson(w, http.StatusBadRequest, "Unable to parse orderId from str to uuid")
			return
		}
		// Check if orderId is in postgresql orders table
		orderDetails := &types.ValidateUser{}
		err = json.NewDecoder(r.Body).Decode(orderDetails)
		if err != nil {
			responses.WriteJson(w, http.StatusBadRequest, "Unable to deserialize body for userId validation")
			return
		}
		userId := orderDetails.UserId
		valid := psqlStore.ValidateUserOrder(*userId, orderId)
		if valid {
			next.ServeHTTP(w, r)
		} else {
			responses.WriteJson(w, http.StatusInternalServerError, "userId does not get validated with orderId")
			return
		}
	})
}
