package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

func (rp RedisProxy) ServePost(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusUnprocessableEntity)
		return
	}

	var kvPair map[string]interface{}
	err = json.Unmarshal(reqBody, &kvPair)
	if err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusUnprocessableEntity)
		return
	}

	var value interface{}
	for k, v := range kvPair {
		_, err = rp.RClient.Set(k, v, 0).Result()
		if err != nil {
			http.Error(w, fmt.Sprintf("ERROR: Failed to set value %s for key %s: %s", v, k, err.Error()), http.StatusInternalServerError)
			return
		}
		value = v
		rp.LocalCache.Remove(k)
	}

	io.WriteString(w, fmt.Sprintf("%v", value))
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

	expTime, err := strconv.Atoi(os.Getenv("EXPIRATION_TIME"))
	if err != nil {
		log.Fatalf("Could not read expiration time %s", os.Getenv("EXPIRATION_TIME"))
	}

	s := RedisProxy{
		RClient:    rClient,
		LocalCache: cache.NewCache(capacity, time.Duration(expTime)*time.Second),
	}

	fmt.Println("Starting Redis proxy")

	mux := bone.New()
	mux.GetFunc("/", s.ServeGet)
	mux.PostFunc("/", s.ServePost)

	http.ListenAndServe(":"+os.Getenv("PORT"), mux)
}
