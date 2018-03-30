package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
)

var rClient *redis.Client

func main() {
	fmt.Println("Starting test suite...")
	fmt.Println("Setting up redis client")
	rClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	err := TestGet()
	if err != nil {
		log.Fatalf("FAILED: %v\n", err)
	}
}

func TestGet() error {
	rClient.Set("k1", "value1", 0)

	log.Print("Testing Gets...")
	v1, err := getFromProxy("k1")
	if err != nil {
		return err
	}

	if v1 != "value1" {
		return fmt.Errorf("Expected value1, got %s\n", v1)
	}

	log.Println("PASS")
	return nil
}

func getFromProxy(key string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:3000?key=%s", key))
	if err != nil {
		return "", fmt.Errorf("Got error response: %v\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Got invalid status code: %v\n", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Could not read response: %v\n", err)
	}

	return string(body), nil
}
