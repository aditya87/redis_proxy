.PHONY: build run test

build:
	docker build . -t redis_proxy

run: build
	docker run -e REDIS_PASSWORD='' \
		-p 3000:3000 \
		-p 7777:7777 \
		redis_proxy

test:
	ginkgo -r .
