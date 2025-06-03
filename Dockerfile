FROM golang:1.24-alpine AS builder
WORKDIR /app

ARG ARCH="amd64"
ARG OS="linux"

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
COPY . .

ENV CGO_ENABLED=0\
    GOOS=${OS}\
    GOARCH=${ARCH}

# Download dependencies and build
RUN go mod download
RUN go build -o mcp-alertmanager cmd/mcp-alertmanager/main.go

# Runtime stage
FROM alpine:latest
WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/mcp-alertmanager /usr/local/bin/mcp-alertmanager

# Default entrypoint
ENTRYPOINT ["mcp-alertmanager"]
