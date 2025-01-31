package cache_test

import (
	"time"

	"github.com/aditya87/redis_proxy/cache"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cache", func() {
	var subject *cache.Cache
	capacity := 3
	expirationTime := 2 * time.Second

	BeforeEach(func() {
		subject = cache.NewCache(capacity, expirationTime)
	})

	It("can be written to and read", func() {
		subject.Set("key1", "value1")
		subject.Set("key2", "value2")

		v1, err := subject.Get("key1")
		Expect(err).NotTo(HaveOccurred())
		Expect(v1).To(Equal("value1"))

		v2, err := subject.Get("key2")
		Expect(err).NotTo(HaveOccurred())
		Expect(v2).To(Equal("value2"))

		_, err = subject.Get("key3")
		Expect(err).To(MatchError(ContainSubstring("key key3 not found")))
	})

	It("can list its keys", func() {
		subject.Set("key1", "value1")
		subject.Set("key2", "value2")
		subject.Set("key3", "value3")

		Expect(subject.Keys()).To(ConsistOf("key1", "key2", "key3"))
	})

	It("uses LRU replacement when it hits its capacity", func() {
		subject.Set("key1", "value1")
		subject.Set("key2", "value2")
		subject.Set("key3", "value3")

		_, _ = subject.Get("key2")
		_, _ = subject.Get("key1")

		subject.Set("key4", "value4")

		Expect(subject.Keys()).To(ConsistOf("key1", "key2", "key4"))
		Expect(len(subject.Keys())).To(Equal(3))

		_, _ = subject.Get("key2")
		_, _ = subject.Get("key4")

		subject.Set("key5", "value5")

		Expect(subject.Keys()).To(ConsistOf("key2", "key4", "key5"))
		Expect(len(subject.Keys())).To(Equal(3))
	})

	It("expires keys after the expiration time is elapsed", func() {
		subject.Set("key1", "value1")
		subject.Set("key2", "value2")

		time.Sleep(3 * time.Second)
		Expect(subject.Keys()).To(BeEmpty())
	})

	It("can delete keys", func() {
		subject.Set("key1", "value1")
		subject.Set("key2", "value2")
		Expect(subject.Keys()).To(ConsistOf("key1", "key2"))

		subject.Remove("key1")
		subject.Remove("key3")
		Expect(subject.Keys()).To(ConsistOf("key2"))
	})
})
