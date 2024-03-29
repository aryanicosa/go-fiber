.PHONY: clean critic security lint test build run swag run-test-dependencies

APP_NAME=service
BUILD_DIR=$(PWD)/build
BUILDER_IMAGE=aryanicosa/go-fiber
REVISION_ID?=$(shell cat version.json | jq -r .version)

swag:
	swag init

# run go only local machine
clean:
	rm -rf ./build

critic:
	gocritic check -enableAll -disable=codegenComment ./...

security:
	gosec ./...

lint:
	golangci-lint run ./...

test: clean critic security lint
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

build: test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) .

run: swag build
	$(BUILD_DIR)/$(APP_NAME)
# end of run local development only using "make run"

####################################
# run service with make file command
docker.fiber.build:
	docker build -t $(BUILDER_IMAGE):$(REVISION_ID) .

docker.run: docker.network docker.postgres docker.redis swag docker.fiber.build docker.fiber

docker.network:
	docker network inspect dev-network >/dev/null 2>&1 || \
	docker network create -d bridge dev-network

docker.fiber:
	docker run --rm -d \
		--name fiber-rest-api-service \
		--network dev-network \
		-p 8080:8080 \
		$(BUILDER_IMAGE):$(REVISION_ID)

docker.postgres:
	docker run --rm -d \
		--name fiber-rest-api-postgres \
		--network dev-network \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=password \
		-e POSTGRES_DB=postgres \
		-v ${PWD}/dev-postgres/data/:/var/lib/postgresql/data \
		-p 5432:5432 \
		postgres

docker.redis:
	docker run --rm -d \
		--name fiber-rest-api-redis \
		--network dev-network \
		-p 6379:6379 \
		redis
# end of run service with make file command
####################################

# run service: go-fiber-app, postgres, and redis with docker compose
docker.run.compose: swag run-docker-compose

docker.stop: docker.stop.fiber docker.stop.postgres docker.stop.redis

run-docker-compose:
	docker-compose -f docker-compose.yml up

docker.stop.compose:
	docker stop fiber-rest-api-service fiber-rest-api-postgres fiber-rest-api-redis
# end of run service: go-fiber-app, postgres, and redis with docker compose

# test in local machine using docker-compose-test
start-test-dependencies:
	docker-compose -f docker-compose-test.yml up

stop-test-dependencies:
	docker stop fiber-rest-api-postgres-test fiber-rest-api-redis-test

run-test:
	go test -v -cover ./...
# end of test in local machine using docker-compose
