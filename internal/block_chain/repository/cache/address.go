package cache

import (
	"context"
	"github.com/leaf-rain/wallet/common"
	"github.com/leaf-rain/wallet/internal/block_chain/model"
	"github.com/leaf-rain/wallet/pkg/log"
	"go.uber.org/zap"
)

func getKeyAddress(currency string) string {
	return common.Address + currency
}

// AddressInset 插入地址池
func (c *cache) AddressInset(ctx context.Context, addrS []model.AddressPrivate) (err error) {
	if len(addrS) == 0 {
		return nil
	}
	if len(addrS[0].Currency) == 0 {
		log.GetLogger().Error("[AddressInset] currency failed", zap.Any("currency", addrS[0].Currency), zap.Error(err))
		return nil
	}
	var key = getKeyAddress(addrS[0].Currency)
	var value = make([]interface{}, len(addrS))
	for i := range addrS {
		value[i] = addrS[i].Address
	}
	if err = c.redis.SAdd(common.AllAddress, value...).Err(); err != nil {
		log.GetLogger().Error("[AddressInset] failed", zap.Any("addrS", addrS), zap.Error(err))
	} else {
		return err
	}
	if err = c.redis.SAdd(key, value...).Err(); err != nil {
		log.GetLogger().Error("[AddressInset] failed", zap.Any("addrS", addrS), zap.Error(err))
	}
	return err
}

// AddressGetTotal 获取地址池数量
func (c *cache) AddressGetTotal(ctx context.Context, currency string) (total int64, err error) {
	var key = getKeyAddress(currency)
	total, err = c.redis.SCard(key).Result()
	if err != nil {
		log.GetLogger().Error("[AddressGetTotal] failed", zap.Any("currency", currency), zap.Error(err))
	}
	return total, err
}

// AddressGet 获取地址
func (c *cache) AddressGet(ctx context.Context, currency string) (address string, err error) {
	var key = getKeyAddress(currency)
	address, err = c.redis.SPop(key).Result()
	if err != nil {
		log.GetLogger().Error("[AddressGet] failed", zap.Any("currency", currency), zap.Error(err))
	}
	return address, err
}

// AddressIsItOurs 是否监听地址
func (c *cache) AddressIsItOurs(ctx context.Context, addr string) (isItOurs bool, err error) {
	isItOurs, err = c.redis.SIsMember(common.AllAddress, addr).Result()
	if err != nil {
		log.GetLogger().Error("[AddressIsItOurs] redis SIsMember failed", zap.Any("addr", addr), zap.Error(err))
	}
	return isItOurs, err
}
