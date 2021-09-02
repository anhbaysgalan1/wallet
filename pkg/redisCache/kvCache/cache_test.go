package kvCache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	rd := NewRedis()

	opt := &Options{
		Redis:        rd,
		LocalCache:   nil,
		StatsEnabled: false,
		Marshal:      nil,
		Unmarshal:    nil,
	}

	cache := New(opt)
	var ctx = context.Background()
	//success, b, err := cache.TryLock("lock:test", 1, time.Second*100, ctx)
	//fmt.Println("success:", success)
	//fmt.Println("b:", b)
	//fmt.Println("err:", err)
	item := &Item{
		Ctx:            ctx,
		Key:            "test_1",
		Value:          "sb",
		TTL:            time.Second * 100,
		Do:             nil,
		SetXX:          false,
		SetNX:          false,
		SkipLocalCache: true,
	}
	err := cache.Set(item)
	fmt.Println("err:", err)
	var s string
	err = cache.Get(ctx, "test_1", &s)
	fmt.Println("err:", err)
	fmt.Println("err:", s)
	result, err := cache.GetS(ctx, []string{"test_1"}, false)
	fmt.Println("result:", result)
	fmt.Println("err:", err)

}

func NewRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "172.12.12.165:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
