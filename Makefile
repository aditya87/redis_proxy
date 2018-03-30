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
		-it redis_proxy \
		/app/integration
