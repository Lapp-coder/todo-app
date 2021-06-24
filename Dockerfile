FROM golang:1.16.4-alpine3.13 AS builder

COPY ./ /github.com/Lapp-coder/todo-app/
WORKDIR /github.com/Lapp-coder/todo-app/

RUN go mod download
RUN go build -o ./build/bin/start ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/Lapp-coder/todo-app/build/bin/start .
COPY --from=0 /github.com/Lapp-coder/todo-app/configs configs/

EXPOSE 8080

CMD ["./start"]
