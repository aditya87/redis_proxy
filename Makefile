.PHONY: build run test

build:
	@echo "Building Docker images..."
	docker build -q -f Dockerfile . -t redis_proxy

# Defaults
capacity := "20" #20 keys
expiry := "30" #30 sec
port := "3000"
redis_host := "localhost"
redis_port := "7000"
redis_pass := ""

run: build
	docker run -e REDIS_PASSWORD=$(redis_pass) \
		-e REDIS_HOST=$(redis_host) \
		-e REDIS_PORT=$(redis_port) \
		-e PORT=$(port) \
		-e CACHE_CAPACITY=$(capacity) \
		-e EXPIRATION_TIME=$(expiry) \
		-p $(port):$(port) \
		redis_proxy \
		/app/redis_proxy

# Run unit tests with ginkgo
unit_test:
	docker build -q -f Dockerfile.test . -t redis_proxy_test
	@echo ''
	@echo "#########Running unit test suite....##########"
	docker run redis_proxy_test
	@echo "DONE UNIT TESTS"

# Run integration tests
test: build
	@echo ''
	@echo "#########Running integration test suite....#############"
	docker run -e REDIS_PASSWORD='' \
		-it redis_proxy \
		/app/integration
