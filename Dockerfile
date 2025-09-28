# 1. Base image Go untuk build
FROM golang:1.25.1-alpine AS builder

# Install dependency build essentials
RUN apk add --no-cache gcc g++ libc-dev libwebp-dev

# Set working directory
WORKDIR /app

# Copy go mod files dan download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy codebase
COPY . .

# Build static binary
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o turubot cmd/bot/main.go

# 2. Minimal runtime image
FROM alpine:latest

# Install ffmpeg & libwebp
RUN apk add --no-cache ffmpeg libwebp

# Copy binary from builder
COPY --from=builder /app/turubot /usr/local/bin/turubot

# Set working directory
WORKDIR /app

# program entrypoint
ENTRYPOINT ["turubot"]
