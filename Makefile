.PHONY: build test

build:
	go get github.com/tools/godep
	go get github.com/onsi/ginkgo/ginkgo
	godep save
	GOOS=linux GOARCH=amd64 go build .

docker_build: build
	docker build . -t redis_proxy

docker_run: docker_build
	docker run -e REDIS_PASSWORD='' \
		-p 3000:3000 \
		-p 7777:7777 \
		redis_proxy

test: build
	ginkgo -r .
