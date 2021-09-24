package db

import (
	"context"
	"github.com/leaf-rain/wallet/common"
	"github.com/leaf-rain/wallet/internal/block_chain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WalletDb interface {
	AddressCreate(ctx context.Context, addr *model.AddressPrivate) (err error)
	AddressSet(ctx context.Context, addr *model.AddressPrivate) (err error)
	AddressSCreate(ctx context.Context, addrS []model.AddressPrivate) (err error)
	AddressGetByFilter(ctx context.Context, filter bson.D, opt *options.FindOptions) (addrS []model.AddressPrivate, err error)
}

type repository struct {
	mongo *mongo.Client
}

func NewWalletDb(ctx context.Context, mg *mongo.Client) WalletDb {
	return &repository{mongo: mg}
}

const (
	collectionName = "address_private"
)

func addressPrivateGetCollection(mg *mongo.Client) *mongo.Collection {
	return mg.Database(common.DbNameForWallet).Collection(collectionName)
}
