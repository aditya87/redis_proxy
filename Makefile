.PHONY: build test

build:
	go get github.com/tools/godep
	go get github.com/onsi/ginkgo/ginkgo
	godep save
	go build .

test: build
	ginkgo -r .
