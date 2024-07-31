package mongostore

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	db *mongo.Database
}

func New(db *mongo.Database) *Store {
	return &Store{db: db}
}
