package v1

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Chandra5468/cfp-Products-Service/internal/middleware"
	"github.com/Chandra5468/cfp-Products-Service/internal/services/external/products"
	"github.com/Chandra5468/cfp-Products-Service/internal/types"
	"github.com/Chandra5468/cfp-Products-Service/internal/utils/responses"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	store    types.OrdersStore
	mdbStore types.MongoComplaintStore
}

func NewHandler(store types.OrdersStore, mdbStore types.MongoComplaintStore) *Handler {
	return &Handler{
		store:    store,
		mdbStore: mdbStore,
	}
}

func (h *Handler) RegisterRoutes(router *chi.Mux) {
	// Create orders for a customer
	router.Post("/v1/api/orders/create", h.createAOrder)

	// Retreive details of specific order of a customer
	router.Get("/v1/api/orders/{id}", h.getAOrder)

	// Submits a complaint of order from customer into mongo. Validate using middleware (jwt valid, purchased the order or not status='delivered')
	router.With(func(next http.Handler) http.Handler {
		return middleware.OrderOwnerShipMiddleware(h.store, next)
	}).Post("/v1/api/orders/register/complaint", h.registerComplaint)

	// ADMIN APIs (Restaurent, KITECHEN, etc... who updates order status)
	router.Patch("/v1/api/admin/orders/update/status", h.updateOrderStatus)

}

func (h *Handler) createAOrder(w http.ResponseWriter, r *http.Request) {
	// body will be like this
	/*
	   	{
	     "user_id": "uuid",
	     "location_id": "uuid", // For now taking it as varchar
	     "items": [
	       {
	         "product_id": "uuid",
	         "quantity": 2
	       },
	       {
	         "product_id": "uuid",
	         "quantity": 1
	       }
	     ]
	   }

	*/

	orderDetails := &types.CreateAOrder{}

	err := json.NewDecoder(r.Body).Decode(orderDetails)
	if err != nil {
		responses.WriteJson(w, http.StatusBadRequest, "unable to deserialize payload")
		return
	}
	if orderDetails.Location == "" || len(orderDetails.Items) == 0 {
		responses.WriteJson(w, http.StatusBadRequest, "empty fields before inserting order")
		return
	}

	wg := &sync.WaitGroup{}

	allProducts := []types.PurchasedProduct{}
	var mu sync.Mutex
	for _, item := range orderDetails.Items {
		if item.ProductId == uuid.Nil || item.Quantity == 0 {
			responses.WriteJson(w, http.StatusBadRequest, "product id or quanity is missing")
			return
		} else {
			wg.Add(1)
			go func(wg *sync.WaitGroup, item *types.Item) {
				defer wg.Done()

				// calling products microservice
				purchasedProducts := products.ValidateCart(item)
				if purchasedProducts != nil {
					mu.Lock()
					allProducts = append(allProducts, *purchasedProducts)
					mu.Unlock()
				}
			}(wg, item)
		}
	}

	wg.Wait()

	/*

				Logic :
						Validate user_id (call User Service). // Will be done in api gateway

		Validate location_id, quantity of products available, fetch unit price of product (call and update the Products Service).

		Calculate total_amount (sum of quantity * unit_price).
		Store in orders and order_items tables.

		Trigger notification (call Notification Service).
	*/

	var totalBill float32

	userInvoice := &types.UserInvoice{}
	for _, eachPurchase := range allProducts {
		if eachPurchase.Status == 1 || eachPurchase.Status == 2 {
			userInvoice.ProductsList = append(userInvoice.ProductsList, &eachPurchase)
			totalBill += eachPurchase.TotalAmount
		}
	}
	userInvoice.TotalAmount = totalBill
	userInvoice.Message = "If there are less quantity of products. We are sorry as we added as per product quantity available"

	orderId := h.store.CreateAOrder(orderDetails, totalBill, "processing")
	h.store.OrderDetailsUpload(orderId, allProducts)

	responses.WriteJson(w, http.StatusCreated, userInvoice)
}
func (h *Handler) getAOrder(w http.ResponseWriter, r *http.Request) {
	// Gets a specific order detail

	/*
				Output response should be

					{
		  "id": "uuid",
		  "user_id": "uuid",
		  "location_id": "uuid",
		  "total_amount": 54.98,
		  "status": "pending",
		  "items": [
		    {
		      "product_id": "uuid",
		      "quantity": 2,
		      "unit_price": 29.99
		    }
		  ],
		  "created_at": "2025-06-24T15:13:00Z",
		  "updated_at": "2025-06-24T15:13:00Z"
		}

	*/

	/*

			Logic:
		Fetch from orders and order_items tables.

		Allow only the order’s user_id (Slef check only allow.) //  Admin can do (Using RBAC) but use a different api
	*/
}

func (h *Handler) updateOrderStatus(w http.ResponseWriter, r *http.Request) {
	//Purpose: Updates the order status (e.g., from pending to shipped). Access only to admin (Kitchen rbac or common admin)
	// Update status in orders table
	// Trigger notification service

}
func (h *Handler) registerComplaint(w http.ResponseWriter, r *http.Request) {
	//Submits a complaint for an order (stored in MongoDB).
	/*
					Request Body (multipart/form-data):
			text: Complaint description (string).
			images: Optional image files.

			Validate order_id and user_id (call User Service).
		Store complaint in MongoDB (complaints_collection).
		Trigger notification for admins (call Notification Service).

	*/

	// h.mdbStore
}
