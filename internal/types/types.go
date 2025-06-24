package types

import (
	"time"

	"github.com/google/uuid"
)

type OrdersStore interface {
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
	Name        string     `json:"name"`
	Description string     `json:"description"` // You can also do Description *string `json:"description"`. So, when psql gives null value we will get nil. Else we will get zero pointed value in this case ""
	Price       float32    `json:"price"`
	Quantity    int16      `json:"quantity"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
