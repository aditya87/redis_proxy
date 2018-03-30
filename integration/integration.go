package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/go-redis/redis"
)

var rClient *redis.Client

func main() {
	fmt.Println("Starting test suite...")

	fmt.Println("Starting redis server...")
	cmd := exec.Command("redis-server", "--port", "7777")
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Could not start redis-server: %v\n", err)
	}
	time.Sleep(3 * time.Second)

	fmt.Println("Setting up environment...")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "7777")
	os.Setenv("PORT", "3000")
	os.Setenv("REDIS_PASSWORD", "")
	os.Setenv("CACHE_CAPACITY", "5")
	os.Setenv("EXPIRATION_TIME", "10")

	fmt.Println("Starting redis proxy with cache size 5 and expiration time 10s...")
	cmd = exec.Command("/app/redis_proxy")
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Could not start redis_proxy: %v\n", err)
	}
	time.Sleep(3 * time.Second)

	fmt.Println("Setting up redis client")
	rClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	err = TestGet()
	if err != nil {
		log.Fatalf("FAILED: %v\n", err)
	}

	err = TestSet()
	if err != nil {
		log.Fatalf("FAILED: %v\n", err)
	}

	err = TestCacheGet()
	if err != nil {
		log.Fatalf("FAILED: %v\n", err)
	}

	err = TestCacheLRUCapacity()
	if err != nil {
		log.Fatalf("FAILED: %v\n", err)
	}

	err = TestCacheExpiry()
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

func TestSet() error {
	log.Print("Testing Sets...")
	err := postToProxy("k2", "value2")
	if err != nil {
		return err
	}

	v2, _ := rClient.Get("k2").Result()
	if v2 != "value2" {
		return fmt.Errorf("Expected value2 to be written, got %s\n", v2)
	}

	log.Println("PASS")
	return nil
}

func TestCacheGet() error {
	rClient.Set("k3", "value3", 0)

	log.Print("Testing Cache Get...")

	s1 := time.Now()
	_, err := getFromProxy("k3")
	if err != nil {
		return err
	}
	t1 := time.Since(s1)

	s2 := time.Now()
	_, err = getFromProxy("k3")
	if err != nil {
		return err
	}
	t2 := time.Since(s2)

	log.Printf("Uncached lookup took %s seconds, cached lookup took %s\n", t1.String(), t2.String())
	if t1 <= t2 {
		return fmt.Errorf("Expected second lookup (took %d time) to be faster than first (which took %d time)", t2, t1)
	}

	log.Println("PASS")
	return nil
}

func TestCacheLRUCapacity() error {
	rClient.Set("k4", "value4", 0)

	log.Print("Testing Cache capacity and LRU policy...")

	_, err := getFromProxy("k1")
	if err != nil {
		return err
	}

	_, err = getFromProxy("k2")
	if err != nil {
		return err
	}

	_, err = getFromProxy("k3")
	if err != nil {
		return err
	}

	_, err = getFromProxy("k4")
	if err != nil {
		return err
	}

	s1 := time.Now()
	_, err = getFromProxy("k1")
	if err != nil {
		return err
	}
	t1 := time.Since(s1)

	s2 := time.Now()
	_, err = getFromProxy("k4")
	if err != nil {
		return err
	}
	t2 := time.Since(s2)

	log.Printf("Uncached lookup took %s seconds, cached lookup took %s\n", t1.String(), t2.String())
	if t1 <= t2 {
		return fmt.Errorf("Expected second lookup (took %d time) to be faster than first (which took %d time)", t2, t1)
	}

	log.Println("PASS")
	return nil
}

func TestCacheExpiry() error {
	rClient.Set("k5", "value5", 0)

	log.Print("Testing Cache expiration...")

	_, err := getFromProxy("k5")
	if err != nil {
		return err
	}

	s1 := time.Now()
	_, err = getFromProxy("k5")
	if err != nil {
		return err
	}
	t1 := time.Since(s1)

	time.Sleep(10 * time.Second)
	s2 := time.Now()
	_, err = getFromProxy("k5")
	if err != nil {
		return err
	}
	t2 := time.Since(s2)

	log.Printf("Uncached lookup took %s seconds, cached lookup took %s\n", t2.String(), t1.String())
	if t2 <= t1 {
		return fmt.Errorf("Expected first lookup (took %d time) to be faster than second (which took %d time)", t1, t2)
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

func postToProxy(key string, value interface{}) error {
	m := map[string]interface{}{
		key: value,
	}

	b, _ := json.Marshal(m)

	resp, err := http.Post("http://localhost:3000", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("Got error response: %v\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Got invalid status code: %v\n", resp.StatusCode)
	}

	return nil
}
