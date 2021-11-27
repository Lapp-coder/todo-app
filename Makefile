.PHONY: build run migrate-up migrate-down test swag
.SILENT:
build:
	docker-compose build todo-app

run:
	docker-compose up todo-app

migrate-up:
	migrate -path ./migrations/ -database "postgres://postgres:${POSTGRES_PASSWORD}@localhost:5432/postgres?sslmode=disable" up

migrate-down:
	migrate -path ./migrations/ -database "postgres://postgres:${POSTGRES_PASSWORD}@localhost:5432/postgres?sslmode=disable" down 

test:
	go test -v -race -cover ./...

mockgen:
	go generate ./internal/service/

swag:
	swag init -g cmd/main.go -o docs/
