.PHONY: build run test

build:
	docker build -f Dockerfile . -t redis_proxy
	docker build -f Dockerfile.test . -t redis_proxy_test

run: build
	docker run -e REDIS_PASSWORD='' \
		-p 3000:3000 \
		-p 7777:7777 \
		redis_proxy

test: build
	docker run redis_proxy_test
