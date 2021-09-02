package common

import (
	"context"
	"github.com/go-redis/redis"
	"gluttonous/pkg/app"
	"gluttonous/pkg/tool"
	"strconv"
)

type CacheIndex interface {
}

type redisIndex struct {
	redis Rediser
}

func NewCacheIndex(redis Rediser) CacheIndex {
	return &redisLock{
		redis: redis,
	}
}

func (r *redisIndex) CacheIndexAdd(ctx context.Context, key string, value []redis.Z) (result int64, err error) {
	return r.redis.ZAdd(key, value...).Result()
}

func (r *redisIndex) CacheIndexGetByOrder(ctx context.Context, key string, page *app.Pagination) (result []int64, err error) {
	if r.redis.Exists(key).Val() != 1 {
		return nil, redis.Nil
	}
	var resp []string
	resp, err = r.redis.ZRevRange(key, int64(page.Offset), int64(page.Offset+page.Limit-1)).Result()
	for i := range resp {
		num, err := strconv.ParseInt(resp[i], 10, 64)
		if err == nil {
			result = append(result, num)
		}
	}
	return result, err
}

func (r *redisIndex) CacheIndexGetByScore(ctx context.Context, key string, page *app.Pagination, start, end int64) (result []int64, err error) {
	if r.redis.Exists(key).Val() != 1 {
		return nil, redis.Nil
	}
	var resp []string
	resp, err = r.redis.ZRangeByScore(key, redis.ZRangeBy{
		Min:    tool.Int64ToStr(start),
		Max:    tool.Int64ToStr(end),
		Offset: int64(page.Offset),
		Count:  int64(page.Limit),
	}).Result()
	for i := range resp {
		num, err := strconv.ParseInt(resp[i], 10, 64)
		if err == nil {
			result = append(result, num)
		}
	}
	return result, err
}
