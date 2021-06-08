.PHONY: build run migrate swag
build:
	go build cmd/todo-app/main.go

run:
	./main

migrate:
	migrate -path ./schema/ -database "postgres://postgres:${POSTGRES_PASSWORD}@localhost:5432/todo-db?sslmode=disable" up

swag:
	swag init -g cmd/todo-app/main.go -o swagger/docs/
