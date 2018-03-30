.PHONY: build run test

build:
	@echo "Building Docker images..."
	docker build -q -f Dockerfile . -t redis_proxy

run: build
	docker run -e REDIS_PASSWORD='' \
		-p 3000:3000 \
		-p 7777:7777 \
		redis_proxy

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
