package cache

import (
	"context"
	"errors"
	hcode "git.huoys.com/chain-game/rowing_proto/common/error"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"tp_wallet/internal/common"
	"tp_wallet/pkg/log"
)

func (c *walletCache) IndexForHashAndBillId(ctx context.Context, hash string, id primitive.ObjectID) (err error) {
	err = c.cache.HSet(common.KeyAllHashToBillId, hash, id.Hex()).Err()
	if err != nil {
		log.GetLogger().Error("[IndexForHashAndBillId] cache.HSet failed ",
			zap.Any("hash", hash),
			zap.Any("id", id),
			zap.Error(err))
		return hcode.ErrServer
	}
	return nil
}

func (c *walletCache) IndexGetBillIdByHash(ctx context.Context, hash string) (primitive.ObjectID, error) {
	var err error
	var id primitive.ObjectID
	var idStr string
	idStr, err = c.cache.HGet(common.KeyAllHashToBillId, hash).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.GetLogger().Error("[IndexForHashAndBillId] cache.HSet failed ",
				zap.Any("hash", hash),
				zap.Error(err))
			return id, hcode.ErrServer
		}
		return [12]byte{}, err
	}
	id, err = primitive.ObjectIDFromHex(idStr)
	if err != nil {
		log.GetLogger().Error("[IndexForHashAndBillId] primitive.ObjectIDFromHex failed ",
			zap.Any("hash", hash),
			zap.Any("idStr", idStr),
			zap.Error(err))
		return id, hcode.ErrServer
	}
	return id, nil
}
