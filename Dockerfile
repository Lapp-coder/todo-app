FROM golang:1.16.4-alpine3.13 AS builder

COPY ./ /github.com/Lapp-coder/todo-app/
WORKDIR /github.com/Lapp-coder/todo-app/

RUN chmod +x wait-for-postgres.sh

RUN go mod download
RUN go build -o ./build/bin/todo-app ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

RUN apk --update --no-cache add postgresql-client 

COPY --from=builder /github.com/Lapp-coder/todo-app/build/bin/todo-app .
COPY --from=builder /github.com/Lapp-coder/todo-app/configs configs/
COPY --from=builder /github.com/Lapp-coder/todo-app/wait-for-postgres.sh .

EXPOSE 8080

CMD ["./todo-app"]
