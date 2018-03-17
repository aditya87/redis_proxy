package main

import (
	"net/http"
	"time"

	"github.com/go-redis/redis"
)

type RedisClient interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(key string) *redis.StringCmd
	Keys(pattern string) *redis.StringSliceCmd
}

type RedisProxy struct {
	RClient RedisClient
}

func (rp RedisProxy) ServeGet(w http.ResponseWriter, r *http.Request) {
}

func main() {
}
