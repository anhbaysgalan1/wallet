package repository

import (
	"gluttonous/pkg/database/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	mgo *mongo.Client
	rds *redis.Client
}

type Repository struct {
	Wallet     *Wallet
	Bill       *Bill
	MgoSession *MgoSession
}

func NewRepository(mgo *mongo.Client, rds *redis.Client) *Repository {
	var a = repository{mgo, rds}
	return &Repository{
		Wallet:     &Wallet{a},
		Bill:       &Bill{a},
		MgoSession: &MgoSession{a},
	}
}
