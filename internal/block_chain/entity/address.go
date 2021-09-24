package entity

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"github.com/leaf-rain/wallet/internal/block_chain/model"
	"github.com/leaf-rain/wallet/internal/block_chain/repository/cache"
	"github.com/leaf-rain/wallet/internal/block_chain/repository/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AddressPrivate struct {
	Value  model.AddressPrivate
	ValueS []model.AddressPrivate
	db     db.WalletDb
	cache  cache.WalletCache
}

func NewAddressPrivate(Value model.AddressPrivate, ValueS []model.AddressPrivate, db db.WalletDb, cache cache.WalletCache) AddressPrivate {
	return AddressPrivate{
		Value:  Value,
		ValueS: ValueS,
		db:     db,
		cache:  cache,
	}
}

// AddressCreate create address and private
func (addr AddressPrivate) AddressCreate(ctx context.Context) (id primitive.ObjectID, err error) {
	err = addr.db.AddressCreate(ctx, &addr.Value)
	if err != nil {
		return [12]byte{}, err
	}
	return addr.Value.Id, nil
}

// AddressSet 创建单个地址私钥对象
func (addr *AddressPrivate) AddressSet(ctx context.Context) (err error) {
	return addr.db.AddressSet(ctx, &addr.Value)
}

// AddressSCreate 创建多个地址私钥对象
func (addr AddressPrivate) AddressSCreate(ctx context.Context) (err error) {
	return addr.db.AddressSCreate(ctx, addr.ValueS)
}

// 获取地址
func (addr *AddressPrivate) GetToUsed(ctx context.Context) (address string, err error) {
	address, err = addr.cache.AddressGet(ctx, addr.Value.Currency)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return address, err
		}
		// 缓存没有，从数据库中读取
		err = addr.addressGetUnUsed(ctx, 1000)
		if err != nil {
			return "", err
		}
		if len(addr.ValueS) > 0 {
			go addr.cache.AddressInset(ctx, addr.ValueS)
			address = addr.ValueS[0].Address
		}
	}
	if len(address) == 0 {
		return "", err
	}
	return address, nil
}

// addressGetUnUsed 获取多个未被使用的地址
func (addr *AddressPrivate) addressGetUnUsed(ctx context.Context, total int64) (err error) {
	var opt = &options.FindOptions{}
	opt.SetLimit(total)
	opt.SetProjection(bson.M{"address": true, "currency": true})
	var filter = bson.D{{"currency", addr.Value.Currency}}
	addr.ValueS, err = addr.db.AddressGetByFilter(ctx, filter, opt)
	if err != nil {
		return err
	}
	return nil
}
