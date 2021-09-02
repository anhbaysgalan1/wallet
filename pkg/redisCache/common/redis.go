package common

import (
	"github.com/go-redis/redis"
)

type Rediser interface {
	redis.Cmdable
}
