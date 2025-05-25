FROM golang:1.24.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main cmd/analysis/main.go

FROM ubuntu:latest

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
