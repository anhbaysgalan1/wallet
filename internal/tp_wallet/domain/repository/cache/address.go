package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"tp_wallet/internal/common"
	"tp_wallet/internal/wallet/entity"
	"tp_wallet/pkg/log"
)

func (c *walletCache) AddressSave(ctx context.Context, addr *entity.Address) error {
	var err error
	err = c.cache.SAdd(common.KeyAddressUidToAddr(addr.Uid), addr.Address).Err()
	if err != nil {
		log.GetLogger().Error("[AddressSave] cache.HSet failed",
			zap.Any("req", addr),
			zap.Error(err))
		return err
	}
	err = c.cache.HSet(common.AddressAddrToUid, addr.Address, addr.Uid).Err()
	if err != nil {
		log.GetLogger().Error("[AddressSave] cache.HSet failed",
			zap.Any("req", addr),
			zap.Error(err))
		return err
	}
	return nil
}

func (c *walletCache) AddressGetByUid(ctx context.Context, uid uint64) (map[string]struct{}, error) {
	var addr map[string]struct{}
	var err error
	addr, err = c.cache.SMembersMap(common.KeyAddressUidToAddr(uid)).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.GetLogger().Error("[AddressGetByUid] cache.HGet failed",
				zap.Any("uid", uid),
				zap.Error(err))
		}
		return nil, err
	}
	return addr, nil
}

func (c *walletCache) AddressSaveByUid(ctx context.Context, uid uint64, addrMap map[string]struct{}) error {
	if len(addrMap) == 0 {
		return nil
	}
	var addr = make([]interface{}, len(addrMap))
	var index = 0
	for k, _ := range addrMap {
		addr[index] = k
		index += 1
	}
	var err error
	_, err = c.cache.SAdd(common.KeyAddressUidToAddr(uid), addr...).Result()
	if err != nil {
		log.GetLogger().Error("[AddressSaveByUid] cache.HGet failed",
			zap.Any("uid", uid),
			zap.Error(err))
		return err
	}
	return nil
}

func (c *walletCache) UidGetByAddr(ctx context.Context, addr string) (uint64, error) {
	var uid uint64
	var err error
	uid, err = c.cache.HGet(common.AddressAddrToUid, addr).Uint64()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.GetLogger().Error("[UidGetByAddr] cache.HGet failed",
				zap.Any("addr", addr),
				zap.Error(err))
		}
		return 0, err
	}
	return uid, nil
}
