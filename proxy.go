package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-zoo/bone"
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
	params := r.URL.Query()
	key := params.Get("key")

	value, err := rp.RClient.Get(key).Result()
	if err != nil {
		http.Error(w, fmt.Sprintf("ERROR: Failed to look up value for key %s: %s", key, err.Error()), http.StatusInternalServerError)
		return
	}

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

	s := RedisProxy{
		RClient: rClient,
	}

	fmt.Println("Starting Redis proxy")

	mux := bone.New()
	mux.GetFunc("/", s.ServeGet)

	http.ListenAndServe(":"+os.Getenv("PORT"), mux)
}
