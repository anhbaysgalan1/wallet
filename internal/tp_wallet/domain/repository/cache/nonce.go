package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"tp_wallet/internal/common"
	"tp_wallet/pkg/hcode"
	"tp_wallet/pkg/log"
)

// NonceGetByAddr 获取地址nonce值
func (c *walletCache) NonceGetByAddr(ctx context.Context, addr string) (uint64, error) {
	var result uint64
	var err error
	result, err = c.cache.HGet(common.NonceAddrHash, addr).Uint64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, hcode.ErrGetAddressNonce
		}
		log.GetLogger().Error("[NonceGetByAddr] cache.HGet failed", zap.String("addr", addr), zap.String("key", common.NonceAddrHash), zap.Error(err))
		return 0, hcode.ErrInternalCache
	}
	return result, nil
}

// NonceIncr 地址nonce值添加固定值
func (c *walletCache) NonceIncr(ctx context.Context, addr string, incr int64) error {
	err := c.cache.HIncrBy(common.NonceAddrHash, addr, incr).Err()
	if err != nil {
		log.GetLogger().Error("[NonceIncr] cache.HIncrBy failed", zap.String("addr", addr), zap.String("key", common.NonceAddrHash), zap.Error(err))
		return hcode.ErrInternalCache
	}
	return nil
}
