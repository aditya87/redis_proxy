package main_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	redis_proxy "github.com/aditya87/redis_proxy"
	"github.com/aditya87/redis_proxy/cache"
	"github.com/aditya87/redis_proxy/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RedisProxy", func() {
	var subject redis_proxy.RedisProxy
	var rClient *fakes.FakeRClient
	var rr *httptest.ResponseRecorder
	var lCache *cache.Cache

	BeforeEach(func() {
		rClient = fakes.NewFakeRClient()
		lCache = cache.NewCache(3)

		subject = redis_proxy.RedisProxy{
			RClient:    rClient,
			LocalCache: lCache,
		}

		rr = httptest.NewRecorder()
	})

	It("redirects HTTP GETs to Redis gets and returns the value in the response", func() {
		rClient.Set("k", "v", 5*time.Second)
		req, _ := http.NewRequest("GET", "?key=k", nil)

		By("calling via the redis client")
		subject.ServeGet(rr, req)
		Expect(rr.Code).To(Equal(http.StatusOK))
		Expect(rClient.GetCalledWith()).To(Equal("k"))
		Expect(rr.Body.String()).To(Equal("v"))

		By("writing to the cache")
		Expect(lCache.Keys()).To(Equal([]string{"k"}))
		Expect(lCache.Get("k")).To(Equal("v"))

		By("subsequently reading from the cache")
		rr = httptest.NewRecorder()
		subject.ServeGet(rr, req)
		Expect(rr.Code).To(Equal(http.StatusOK))
		Expect(rClient.GetCalledWith()).To(Equal(""))
		Expect(rr.Body.String()).To(Equal("v"))
	})

	It("returns an error response if the Redis backend throws an error", func() {
		rClient.Set("k", "v", 5*time.Second)
		rClient.SetError(errors.New("some error"))
		req, _ := http.NewRequest("GET", "?key=k", nil)

		subject.ServeGet(rr, req)
		Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		Expect(rClient.GetCalledWith()).To(Equal("k"))
		Expect(rr.Body.String()).To(ContainSubstring("Failed to look up value for key k: some error"))
	})
})
