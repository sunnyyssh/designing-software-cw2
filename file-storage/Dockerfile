FROM golang:1.24.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG GOARCH

RUN GOOS=linux GOARCH=${GOARCH} go build -o main cmd/gateway/main.go

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 80

CMD ["./main"]
