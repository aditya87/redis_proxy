package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aditya87/redis_proxy/cache"
	"github.com/go-redis/redis"
	"github.com/go-zoo/bone"
)

type RedisClient interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(key string) *redis.StringCmd
	Keys(pattern string) *redis.StringSliceCmd
}

type RedisProxy struct {
	RClient    RedisClient
	LocalCache *cache.Cache
}

func (rp RedisProxy) ServeGet(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get("key")

	value, err := rp.LocalCache.Get(key)
	if err == nil {
		io.WriteString(w, value)
		return
	}

	value, err = rp.RClient.Get(key).Result()
	if err != nil {
		http.Error(w, fmt.Sprintf("ERROR: Failed to look up value for key %s: %s", key, err.Error()), http.StatusInternalServerError)
		return
	}

	rp.LocalCache.Set(key, value)

	io.WriteString(w, value)
}

func main() {
	rClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	capacity, err := strconv.Atoi(os.Getenv("CACHE_CAPACITY"))
	if err != nil {
		log.Fatalf("Could not read cache capacity %s", os.Getenv("CACHE_CAPACITY"))
	}

	s := RedisProxy{
		RClient:    rClient,
		LocalCache: cache.NewCache(capacity),
	}

	fmt.Println("Starting Redis proxy")

	mux := bone.New()
	mux.GetFunc("/", s.ServeGet)

	http.ListenAndServe(":"+os.Getenv("PORT"), mux)
}
