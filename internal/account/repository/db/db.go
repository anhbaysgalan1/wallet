package db

import (
	"context"
	"github.com/leaf-rain/wallet/internal/account/entity"
	"github.com/leaf-rain/wallet/internal/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AccountDb interface {
	AddressCreate(ctx context.Context, addr *entity.EntityAddressPrivate) (err error)
	AddressSet(ctx context.Context, addr *entity.EntityAddressPrivate) (err error)
	AddressSCreate(ctx context.Context, addrS []*entity.EntityAddressPrivate) (err error)
	AddressGetByFilter(ctx context.Context, filter bson.D, opt *options.FindOptions) (addrS []*entity.EntityAddressPrivate, err error)
}

type accountDb struct {
	mongo *mongo.Client
}

func NewWalletDb(ctx context.Context, mg *mongo.Client) AccountDb {
	return &accountDb{mongo: mg}
}

const (
	collectionName = "address_private"
)

func addressPrivateGetCollection(mg *mongo.Client) *mongo.Collection {
	return mg.Database(common.DbNameForWallet).Collection(collectionName)
}
