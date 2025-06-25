package types

import (
	"time"

	"github.com/google/uuid"
)

type OrdersStore interface {
	CreateAOrder(orderDetails *CreateAOrder, totalBill float32, status string) *uuid.UUID
	GetSelfOrdersDetail()
	UpdateOrderStatus()
	ValidateUserOrder(userId uuid.UUID, orderId uuid.UUID) bool
	OrderDetailsUpload(orderId *uuid.UUID, allProducts []PurchasedProduct)
}
type MongoComplaintStore interface {
	RegisterComplaint(userId uuid.UUID, orderId uuid.UUID, itemId uuid.UUID)
}

// Declaring structs here because :
/*
	Model or Domain Layer: The struct representing your database entity should reside in a separate package,
	typically named types or models,
	which is independent of the business logic (services).
	This allows easy reusability and separation of concerns.


	Why put the struct in a separate types package?
Reusability: The struct can be used by different parts of your application (e.g., database layer, API layer, etc.).

Separation of concerns: Keeping data types separate from the business logic makes your code more modular and easier to test.

*/
type Orders struct {
	Id          *uuid.UUID `json:"id,omitempty"` // using pointer instead of uuid.UUID because pointer will give nil if there is null(record not found) in sql. And null can be omitted not zero pointed vales
	UserId      *uuid.UUID `json:"user_id"`
	Location    string     `json:"location"` // You can also do Location *string `json:"location"`. So, when psql gives null value we will get nil. Else we will get zero pointed value in this case ""
	TotalAmount float32    `json:"total_amount"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Item struct {
	ProductId uuid.UUID `json:"product_id"`
	Quantity  int16     `json:"quantity"`
}

type CreateAOrder struct {
	UserId   uuid.UUID `json:"user_id"`
	Location string    `json:"location"`
	Items    []*Item   `json:"items,omitempty"`
}

type ValidateUser struct {
	Id     *uuid.UUID `json:"order_id,omitempty"` // This will be orderId
	UserId *uuid.UUID `json:"user_id"`
}

type PurchasedProduct struct {
	ProductId   uuid.UUID `json:"product_id"`
	Quantity    int16     `json:"quantity"`
	TotalAmount float32   `json:"total_amount"`
	Status      int8      `json:"status"` // 0: Products are not available. Out of stock 1: all products purchased, 2: only few products are available
}

//	type PurchasedProducts struct {
//		AllPurchases []PurchasedProduct `json:"final_settlement"`
//	}
type UserInvoice struct {
	ProductsList []*PurchasedProduct `json:"purchases"`
	TotalAmount  float32             `json:"total_amount"`
	Message      string              `json:"message"`
}
