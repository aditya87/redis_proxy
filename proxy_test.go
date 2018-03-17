package main_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	redis_proxy "github.com/aditya87/redis_proxy"
	"github.com/aditya87/redis_proxy/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RedisProxy", func() {
	var subject redis_proxy.RedisProxy
	var rClient *fakes.FakeRClient
	var rr *httptest.ResponseRecorder

	BeforeEach(func() {
		rClient = fakes.NewFakeRClient()
		subject = redis_proxy.RedisProxy{
			RClient: rClient,
		}

		rr = httptest.NewRecorder()
	})

	It("redirects HTTP GETs to Redis gets", func() {
		req, _ := http.NewRequest("GET", "/", bytes.NewBuffer([]byte(`key`)))
		subject.ServeGet(rr, req)
		Expect(rr.Code).To(Equal(http.StatusOK))
	})
})
