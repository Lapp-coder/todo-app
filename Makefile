.PHONY: build run migrate test swag
build:
	docker-compose build todo-app

run:
	docker-compose up todo-app

migrate:
	migrate -path ./schema/ -database "postgres://postgres:${POSTGRES_PASSWORD}@localhost:5436/postgres?sslmode=disable" up

test:
	go test -v -race -cover ./...

swag:
	${HOME}/go/bin/swag init -g cmd/main.go -o swagger/docs/
