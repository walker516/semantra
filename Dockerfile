# ./app/api/Dockerfile

FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o semantra-api ./cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/semantra-api .
COPY config.prod.json ./config.json

EXPOSE 8081
CMD ["./semantra-api"]
