package redis

import (
	redisio "github.com/hoisie/redis"
	"os"
)

var (
	RedisClient = NewRedisClient(13)
	RedisMaster = NewRedisMaster(13)
)

func NewClient(poolsize int) (c *redisio.Client){
	_c := new(redisio.Client)
	_c.MaxPoolSize = poolsize

	return _c
}

func NewRedisClient(db int) *redisio.Client {
	c := NewClient(50)

	c.Addr = "127.0.0.1:6379"
	if addr := os.Getenv("REDIS_CLIENT_ADDRESS"); addr != "" {
		c.Addr = addr
	}

	c.Db = db
	c.Password = os.Getenv("REDIS_PASSWORD")

	return c
}

func NewRedisMaster(db int) *redisio.Client {
	c := NewClient(50)

	c.Addr = "127.0.0.1:6379"
	if addr := os.Getenv("REDIS_MASTER_ADDRESS"); addr != "" {
		c.Addr = addr
	}

	c.Db = db
	c.Password = os.Getenv("REDIS_PASSWORD")

	return c
}
