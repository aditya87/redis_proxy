package main

import (
	"fmt"
	"io"
	"io/ioutil"
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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error", http.StatusUnprocessableEntity)
		return
	}

	key := string(body)
	value, err := rp.RClient.Get(key).Result()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, value)
}

func main() {
}
