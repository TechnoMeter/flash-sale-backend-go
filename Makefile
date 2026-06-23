.PHONY: build run test docker-build docker-up migrate

build:
	go build -o flash-sale ./cmd/api

run:
	go run ./cmd/api

test:
	go test -race -v ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate:
	docker exec -i $$(docker ps -qf "name=postgres") psql -U postgres -d flash_sale < migrations/001_init.up.sql