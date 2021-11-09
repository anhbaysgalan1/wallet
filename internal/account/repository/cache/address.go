package cache

import (
	"context"
	"github.com/leaf-rain/wallet/internal/account/entity"
	"github.com/leaf-rain/wallet/internal/common"
	"github.com/leaf-rain/wallet/pkg/hcode"
	"github.com/leaf-rain/wallet/pkg/log"
	"go.uber.org/zap"
	"time"
)

func getKeyAddress(currency string) string {
	return common.Address + currency
}

// AddressInset 插入地址池
func (c *cache) AddressInset(ctx context.Context, addrS []*entity.EntityAddressPrivate) (err error) {
	if len(addrS) == 0 {
		return nil
	}
	if len(addrS[0].Currency) == 0 {
		log.GetLogger().Error("[AddressInset] currency failed", zap.Any("currency", addrS[0].Currency), zap.Error(err))
		return nil
	}
	var key = getKeyAddress(addrS[0].Currency)
	var valueMap = make(map[string]interface{}, len(addrS))
	var valueList = make([]interface{}, len(addrS))
	for i := range addrS {
		valueMap[addrS[i].Address] = addrS[i].Id.Hex()
		valueList[i] = addrS[i].Address
	}
	if err = c.redis.HMSet(common.AllAddress, valueMap).Err(); err != nil {
		log.GetLogger().Error("[AddressInset] failed", zap.Any("addrS", addrS), zap.Error(err))
	}
	if err = c.redis.SAdd(key, valueList...).Err(); err != nil {
		log.GetLogger().Error("[AddressInset] failed", zap.Any("addrS", addrS), zap.Error(err))
	} else {
		return err
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
func (c *cache) AddressGet(ctx context.Context, currency string) (address, id string, err error) {
	for i := 0; i <= 10; i++ { // 重试10次
		address, id, err = c.addressGet(ctx, currency)
		if err == nil && len(address) >= 0 && len(id) >= 0 {
			return
		}
		time.Sleep(time.Second / 2)
	}
	return "", "", hcode.ErrAddressGet
}

func (c *cache) addressGet(ctx context.Context, currency string) (address, id string, err error) {
	var key = getKeyAddress(currency)
	address, err = c.redis.SPop(key).Result()
	if err != nil {
		log.GetLogger().Error("[AddressGet] redis.SPop failed",
			zap.Any("currency", currency),
			zap.Any("key", key),
			zap.Error(err))
	}
	id, err = c.redis.HGet(common.AllAddress, address).Result()
	if err != nil {
		log.GetLogger().Error("[AddressGet] redis.HGet failed",
			zap.Any("currency", currency),
			zap.Any("key", common.AllAddress),
			zap.Any("address", address),
			zap.Error(err))
	}
	return address, id, err
}

// AddressIsItOurs 是否监听地址
func (c *cache) AddressIsItOurs(ctx context.Context, addr string) (isItOurs bool, err error) {
	isItOurs, err = c.redis.SIsMember(common.AllAddress, addr).Result()
	if err != nil {
		log.GetLogger().Error("[AddressIsItOurs] redis SIsMember failed", zap.Any("addr", addr), zap.Error(err))
	}
	return isItOurs, err
}
