package redistool

import (
	"github.com/garyburd/redigo/redis"
)

var RedisPool *redis.Pool

func init() {
	RedisPool = GetRedis("redis://localhost:6379")
}
