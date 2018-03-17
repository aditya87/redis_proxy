package fakes

import (
	"time"

	"github.com/go-redis/redis"
)

type FakeRClient struct {
	getCalledWithKey string
	users            map[string]interface{}
	err              error
}

func NewFakeRClient() *FakeRClient {
	return &FakeRClient{
		users: make(map[string]interface{}),
		err:   nil,
	}
}

func (frc *FakeRClient) SetError(err error) {
	frc.err = err
}

func (frc *FakeRClient) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	frc.users[key] = value
	return redis.NewStatusResult("", frc.err)
}

func (frc *FakeRClient) Get(key string) *redis.StringCmd {
	frc.getCalledWithKey = key

	s, ok := frc.users[key].(string)
	if !ok {
		return redis.NewStringResult("", frc.err)
	}

	return redis.NewStringResult(s, frc.err)
}

func (frc *FakeRClient) GetCalledWith() string {
	return frc.getCalledWithKey
}

func (frc *FakeRClient) Keys(pattern string) *redis.StringSliceCmd {
	keys := make([]string, 0)
	for k := range frc.users {
		keys = append(keys, k)
	}
	return redis.NewStringSliceResult(keys, frc.err)
}
