package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRedisProxy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RedisProxy Suite")
}
