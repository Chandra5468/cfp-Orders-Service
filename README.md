# cfp-Orders-Service
This microservice will handle orders, status, take complaints, connect with products and connect with notification service.


middleware for specific routes only

router.Route("/v1/api", func(r chi.Router) {
    r.Get("/public", publicHandler) // No middleware

    r.With(middleware.CorsHandler).Get("/product/{productId}", productHandler) // Only this route has CORS

    r.With(middleware.LoggingMiddleware).Post("/product", createProductHandler) // Only POST has logging
})


or 

router.Route("/admin", func(r chi.Router) {
    r.Use(middleware.AdminAuthMiddleware) // All /admin routes will have this middleware

    r.Get("/dashboard", dashboardHandler)
    r.Get("/settings", settingsHandler)
})


CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    location varchar(100) NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'processing', 'shipped', 'delivered', 'cancelled')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) on delete cascade, 
  # If you want order_items to be automatically deleted when their order is deleted:
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


Scope: orders is the parent entity (one order), while order_items is the child entity (multiple items per order).
Data: orders stores order-level data (e.g., total cost, status), while order_items stores item-level data (e.g., product ID, quantity).
Relationship: One order can have multiple order_items (one-to-many), linked by order_id.

Example:
A customer orders two Chocolate Cakes and one Vanilla Cake in a single order.
orders table: One record with id, user_id, location_id, total_amount (sum of all items), status (e.g., pending).
order_items table: Two records:
Record 1: order_id, product_id (Chocolate Cake), quantity (2), unit_price (e.g., 29.99).
Record 2: order_id, product_id (Vanilla Cake), quantity (1), unit_price (e.g., 24.99).

{
  "_id": ObjectId,
  "order_id": UUID,
  "user_id": UUID,
  "text": String,
  "images": [
    {
      "filename": String,
      "content_type": String,
      "data": Binary
    }
  ],
  "created_at": ISODate,
  "updated_at": ISODate
}