package orders

import (
	"database/sql"
	"log"
	"sync"

	"github.com/Chandra5468/cfp-Products-Service/internal/types"
	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateAOrder(orderDetails *types.CreateAOrder, totalBill float32, status string) *uuid.UUID {
	query := `insert into orders(user_id, location, total_amount, status) values ($1,$2,$3,$4) returning id`
	var id uuid.UUID
	err := s.db.QueryRow(query, &orderDetails.UserId, &orderDetails.Location, &totalBill, &status).Scan(&id)
	if err != nil {
		log.Printf("Error inserting order : %v", err)
		return nil
	}
	return &id
}

func (s *Store) GetSelfOrdersDetail() {

}

func (s *Store) UpdateOrderStatus() {

}

func (s *Store) ValidateUserOrder(userId uuid.UUID, orderId uuid.UUID) bool {
	return false
}

func (s *Store) OrderDetailsUpload(orderId *uuid.UUID, allProducts []types.PurchasedProduct) {
	wg := &sync.WaitGroup{}
	for _, eachPurchase := range allProducts {
		wg.Add(1)
		go func(ep *types.PurchasedProduct) {
			defer wg.Done()
			query := `insert into order_items (order_id, product_id, quantity, unit_price) values ($1,$2,$3,$4)`
			_, err := s.db.Exec(query, orderId, ep.ProductId, ep.Quantity, ep.TotalAmount)
			if err != nil {
				log.Printf("Error inserting into order_items %s", err.Error())
			} else {
				log.Println("Inserted document ")
			}
		}(&eachPurchase)
	}

	wg.Wait()
}
