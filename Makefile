.PHONY: clean critic security lint test build run swag run-dependencies run-test-dependencies

APP_NAME = server
BUILD_DIR = $(PWD)/build

clean:
	rm -rf ./build

critic:
	gocritic check -enableAll ./...

security:
	gosec ./...

lint:
	golangci-lint run ./...

test: clean critic security lint
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

build: test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

run: swag build
	$(BUILD_DIR)/$(APP_NAME)

docker.run: docker.network docker.postgres swag docker.fiber docker.redis

docker.network:
	docker network inspect dev-network >/dev/null 2>&1 || \
	docker network create -d bridge dev-network

docker.fiber.build:
	docker build -t fiber .

docker.fiber: docker.fiber.build
	docker run --rm -d \
		--name go-fiber \
		--network dev-network \
		-p 8080:8080 \
		fiber

docker.postgres:
	docker run --rm -d \
		--name go-fiber-postgres \
		--network dev-network \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=password \
		-e POSTGRES_DB=postgres \
		-v ${PWD}/dev-postgres/data/:/var/lib/postgresql/data \
		-p 5432:5432 \
		postgres

docker.redis:
	docker run --rm -d \
		--name go-fiber-redis \
		--network dev-network \
		-p 6379:6379 \
		redis

docker.stop: docker.stop.fiber docker.stop.postgres docker.stop.redis

docker.stop.fiber:
	docker stop go-fiber

docker.stop.postgres:
	docker stop go-fiber-postgres

docker.stop.redis:
	docker stop go-fiber-redis

swag:
	swag init

# run in local machine using docker-compose
run-dependencies:
	docker-compose -f docker-compose-dependencies.yml up

stop-dependencies:
	docker stop fiber-rest-api-postgres fiber-rest-api-redis

run-local:
	go run main.go

# test in local machine using docker-compose
run-test-dependencies:
	docker-compose -f docker-compose-test.yml up

stop-test-dependencies:
	docker stop fiber-rest-api-postgres-test fiber-rest-api-redis-test

run-test:
	go test -v -cover ./...
