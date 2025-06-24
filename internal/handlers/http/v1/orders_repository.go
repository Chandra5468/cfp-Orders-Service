package v1

import (
	"github.com/Chandra5468/cfp-Products-Service/internal/types"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store types.OrdersStore
}

func NewHandler(store types.OrdersStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *chi.Mux) {

}
