package common

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"runtime"
	"time"
)

type redisLock struct {
	redis Rediser
}

type LockType int

const (
	Writer LockType = 1
	Read   LockType = 2
)

type RedisLock interface {
	// 加锁, 返回加锁是否成功
	TryLock(key string, val interface{}, ttl time.Duration, ctx context.Context) (bool, []byte, error)
	// 解锁
	UnLock(ctx context.Context, key string, val interface{}) error
	// 等待释放锁
	WaitUnLock(ctx context.Context, key string) (err error)
	// 续约
	RenewLock(ctx context.Context, key string, value interface{}, ttl time.Duration) (err error)
	// 加锁, 如果有锁，等待锁释放再添加锁
	WaitLock(key string, val interface{}, ttl time.Duration, ctx context.Context) (bool, error)
}

func NewRedisLock(redis Rediser) RedisLock {
	return &redisLock{
		redis: redis,
	}
}

// 加锁, 返回加锁是否成功
func (c redisLock) TryLock(key string, val interface{}, ttl time.Duration, ctx context.Context) (bool, []byte, error) {
	var lockKey = NewLockKey(key)
	var result []byte
	success, err := c.redis.SetNX(lockKey, val, ttl).Result()
	if err != nil {
		return false, nil, err
	}
	if !success {
		result, err = c.redis.Get(lockKey).Bytes()
		if errors.Is(err, redis.Nil) {
			return c.TryLock(key, val, ttl, ctx)
		}
		return success, result, err
	}
	return success, result, err
}

// 加锁, 如果有锁，等待锁释放再添加锁
func (c redisLock) WaitLock(key string, val interface{}, ttl time.Duration, ctx context.Context) (bool, error) {
	var lockKey = NewLockKey(key)
	err := c.WaitUnLock(ctx, lockKey)
	if err != nil {
		return false, err
	}
	lua := redis.NewScript("if redis.call('exists', KEYS[1]) == 0 " +
		"then return redis.call('set', KEYS[1], ARGV[1], 'EX', ARGV[2]) " +
		"else return 'NO' end")
	cmd := lua.Run(c.redis, []string{lockKey}, val, int(ttl/time.Second))
	if cmd.Err() != nil {
		return false, cmd.Err()
	} else {
		result, _ := cmd.String()
		if result == "OK" {
			return true, cmd.Err()
		} else {
			return false, cmd.Err()
		}
	}
}

// 解锁
func (c redisLock) UnLock(ctx context.Context, key string, val interface{}) error {
	var lockKey = NewLockKey(key)
	luaDel := redis.NewScript("if redis.call('get',KEYS[1]) == ARGV[1] then " +
		"return redis.call('del',KEYS[1]) else return 0 end")
	return luaDel.Run(c.redis, []string{lockKey}, val).Err()
}

// 等待释放锁
func (c redisLock) WaitUnLock(ctx context.Context, key string) (err error) {
	var lockKey = NewLockKey(key)
	ttl, err := c.redis.TTL(lockKey).Result()
	if err != nil {
		return err
	}
	if ttl > 0 {
		time.Sleep(time.Second * 1)
		runtime.Gosched()
		return c.WaitUnLock(ctx, key)
	}
	return nil
}

// 续约
func (c redisLock) RenewLock(ctx context.Context, key string, value interface{}, ttl time.Duration) (err error) {
	return c.redis.Set(NewLockKey(key), value, ttl).Err()
}

// 拼接锁key
func NewLockKey(key string) (lockKey string) {
	return key + ":lock"
}
