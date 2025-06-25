package orders

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Store struct {
	mdb *mongo.Client
}

func NewStore(mdb *mongo.Client) *Store {
	return &Store{
		mdb: mdb,
	}
}

func (m *Store) RegisterComplaint(userId uuid.UUID, orderId uuid.UUID, itemId uuid.UUID) {

}
