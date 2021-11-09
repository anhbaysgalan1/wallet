package cache

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/leaf-rain/wallet/internal/account/entity"
	"github.com/leaf-rain/wallet/pkg/redisCache/kvCache"
)

type AccountCache interface {
	// AddressGet 获取地址
	AddressGet(ctx context.Context, currency string) (address, id string, err error)
	// AddressInset 插入地址池
	AddressInset(ctx context.Context, addrS []*entity.EntityAddressPrivate) (err error)
	// AddressGetTotal 获取地址池数量
	AddressGetTotal(ctx context.Context, currency string) (total int64, err error)
	// AddressIsItOurs 是否监听地址
	AddressIsItOurs(ctx context.Context, addr string) (isItOurs bool, err error)
}

type cache struct {
	redis *redis.Client
	kv    *kvCache.Cache
}

func NewWalletCache(ctx context.Context, redis *redis.Client) AccountCache {
	c := &cache{redis: redis}
	c.kv = kvCache.New(&kvCache.Options{Redis: redis})
	return c
}
