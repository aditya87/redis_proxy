package fakes

import (
	"time"

	"github.com/go-redis/redis"
)

type FakeRClient struct {
	getCalledWithKey   string
	setCalledWithKey   string
	setCalledWithValue interface{}
	users              map[string]interface{}
	err                error
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
	frc.setCalledWithKey = key
	frc.setCalledWithValue = value
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
	key := frc.getCalledWithKey
	frc.getCalledWithKey = ""
	return key
}

func (frc *FakeRClient) SetCalledWith() (string, interface{}) {
	key := frc.setCalledWithKey
	value := frc.setCalledWithValue
	frc.setCalledWithKey = ""
	frc.setCalledWithValue = nil
	return key, value
}

func (frc *FakeRClient) Keys(pattern string) *redis.StringSliceCmd {
	keys := make([]string, 0)
	for k := range frc.users {
		keys = append(keys, k)
	}
	return redis.NewStringSliceResult(keys, frc.err)
}
