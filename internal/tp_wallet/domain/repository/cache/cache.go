package cache

import (
	"context"
	"encoding/json"
	"tp_wallet/internal/wallet/entity"
	"tp_wallet/pkg/database/redis"
	"tp_wallet/pkg/redisCache/common"
	"tp_wallet/pkg/redisCache/kvCache"
)

type WalletCache interface {
	BalanceSave(ctx context.Context, addr []*entity.Balance) error
	BalanceDel(ctx context.Context, uid uint64) error
	BalanceGetByUid(ctx context.Context, uid uint64) (map[string]entity.Amount, error)
	BalanceSysSet(ctx context.Context, sysUid uint64, amount, currency string, IsAdd bool) (entity.Amount, entity.Amount, error)

	AddressSave(ctx context.Context, addr *entity.Address) error
	AddressGetByUid(ctx context.Context, uid uint64) (map[string]struct{}, error)
	AddressSaveByUid(ctx context.Context, uid uint64, addrMap map[string]struct{}) error
	UidGetByAddr(ctx context.Context, addr string) (uint64, error)

	NonceGetByAddr(ctx context.Context, addr string) (uint64, error)
	NonceIncr(ctx context.Context, addr string, incr int64) error
}

type walletCache struct {
	cache *redis.Client
	kv    *kvCache.Cache
	lock  common.RedisLock
}

func NewCache(d *redis.Client) WalletCache {
	wc := &walletCache{cache: d}
	wc.kv = kvCache.New(&kvCache.Options{
		Redis: d,
	})
	wc.lock = common.NewRedisLock(d)
	return wc
}

func (c *walletCache) Close() {
	if c.cache != nil {
		_ = c.cache.Close()
	}
}

func (c *walletCache) CacheMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c *walletCache) CacheUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
