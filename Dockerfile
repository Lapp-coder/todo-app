FROM golang:1.16.4

COPY ./ ./

ENV GOPATH=/

RUN go mod download
RUN go build -o build/bin/start cmd/main.go

CMD ["./build/bin/start"]
