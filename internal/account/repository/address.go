package repository

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"github.com/leaf-rain/wallet/internal/account/dto"
	"github.com/leaf-rain/wallet/internal/account/entity"
	"github.com/leaf-rain/wallet/internal/account/repository/cache"
	"github.com/leaf-rain/wallet/internal/account/repository/db"
	"github.com/leaf-rain/wallet/internal/common"
	"github.com/leaf-rain/wallet/pkg/hcode"
	"github.com/leaf-rain/wallet/pkg/log"
	"github.com/leaf-rain/wallet/pkg/tool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type AccountRepo struct {
	db    db.AccountDb
	cache cache.AccountCache
}

func NewAddressRepo(db db.AccountDb, cache cache.AccountCache) AccountRepo {
	return AccountRepo{
		db:    db,
		cache: cache,
	}
}

// GetToUsed 获取地址
func (a *AccountRepo) GetToUsed(ctx context.Context, currency, remarks string) (address string, err error) {
	var id string
	address, id, err = a.cache.AddressGet(ctx, currency)
	defer a.CheckInventory(ctx, currency)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return address, err
		}
		var addrs []*entity.EntityAddressPrivate
		// 缓存没有，从数据库中读取
		addrs, err = a.addressGetUnUsed(ctx, currency, common.AddressPoolTotal)
		if err != nil {
			return "", err
		}
		if len(addrs) > 0 {
			go a.cache.AddressInset(ctx, addrs[1:])
			address = addrs[0].Address
		}
	}
	if len(address) == 0 {
		return "", hcode.ErrAddressGet
	} else {
		_id, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.GetLogger().Error("[GetToUsed] primitive.ObjectIDFromHex failed",
				zap.Any("currency", currency),
				zap.String("remarks", remarks),
				zap.String("addr", address),
				zap.String("id", id),
				zap.Error(err))
			return "", hcode.ErrAddressGet
		}
		var now = tool.GetTimeUnixMilli()
		if err = a.db.AddressSet(ctx, &entity.EntityAddressPrivate{
			Id:         _id,
			Status:     dto.AccountType_AccountType_Used,
			CreateTime: now,
			UpdateTime: now,
			Remarks:    remarks,
		}); err != nil {
			return "", hcode.ErrAddressGet
		}
	}
	return address, nil
}

func (a *AccountRepo) CheckInventory(ctx context.Context, currency string) {
	var total int64
	var err error
	var addrs []*entity.EntityAddressPrivate
	total, err = a.cache.AddressGetTotal(ctx, currency)
	if err != nil && !errors.Is(err, redis.Nil) {
		return
	}
	if total <= common.AddressPoolTotal {
		addrs, err = a.addressGetUnUsed(ctx, currency, common.AddressPoolTotal)
		if err != nil {
			return
		}
		if len(addrs) > 0 {
			go a.cache.AddressInset(ctx, addrs)
		}
		if int64(len(addrs)) < common.AddressPoolTotal { // todo:地址数量不足，自动生成

		}
	}
}

// addressGetUnUsed 获取多个未被使用的地址
func (a *AccountRepo) addressGetUnUsed(ctx context.Context, currency string, total int64) (result []*entity.EntityAddressPrivate, err error) {
	var opt = &options.FindOptions{}
	opt.SetLimit(total)
	opt.SetProjection(bson.M{"address": true, "currency": true, "_id": true})
	var filter = bson.D{{"currency", currency}, {"status", dto.AccountType_AccountType_UnUsed}}
	result, err = a.db.AddressGetByFilter(ctx, filter, opt)
	if err != nil {
		return nil, err
	}
	return result, nil
}
