# 1. Base image Go untuk build
FROM golang:1.25.1-bullseye AS builder

RUN apt-get update && apt-get install -y \
    build-essential libwebp-dev ffmpeg \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o turubot cmd/bot/main.go

# 2. Runtime
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y \
    ffmpeg libwebp7 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/turubot /usr/local/bin/turubot

WORKDIR /app
ENTRYPOINT ["turubot"]
