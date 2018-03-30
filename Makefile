.PHONY: build run test

build:
	@echo "Building Docker images..."
	docker build -q -f Dockerfile . -t redis_proxy
	docker build -q -f Dockerfile.test . -t redis_proxy_test

run: build
	docker run -e REDIS_PASSWORD='' \
		-p 3000:3000 \
		-p 7777:7777 \
		redis_proxy

test: build
	@echo ''
	@echo "#########Running unit test suite....##########"
	docker run redis_proxy_test
	@echo "DONE UNIT TESTS"
	@echo ''
	@echo "#########Running integration test suite....#############"
	docker run -e REDIS_PASSWORD='' \
		-e PORT=3000 \
	  -e REDIS_HOST=localhost \
		-e EXPIRATION_TIME=10 \
		-e CACHE_CAPACITY=5 \
		-e REDIS_PORT=7777 \
		-p 3000:3000 \
		-p 7777:7777 \
		-it redis_proxy \
		/app/integration
