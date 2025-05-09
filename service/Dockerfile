FROM golang:1.24.3 AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

ARG GOARCH

RUN GOOS=linux GOARCH=${GOARCH} go build -o main main.go

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
