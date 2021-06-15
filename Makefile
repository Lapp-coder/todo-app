.PHONY: build run migrate test swag
build:
	go build -o build/bin/start cmd/main.go

run:
	./build/bin/start

migrate:
	migrate -path ./schema/ -database "postgres://postgres:${POSTGRES_PASSWORD}@localhost:5432/todo-db?sslmode=disable" up

test:
	go test -v -race -cover ./...

swag:
	${HOME}/go/bin/swag init -g cmd/main.go -o swagger/docs/
