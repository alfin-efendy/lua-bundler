# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags "-X main.Version=${VERSION} -s -w" \
    -o lua-bundler .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/lua-bundler .

# Create a non-root user
RUN adduser -D -s /bin/sh appuser

# Create directories for input/output
RUN mkdir -p /app/input /app/output && \
    chown -R appuser:appuser /app

USER appuser
WORKDIR /app

EXPOSE 8080

# Default command
ENTRYPOINT ["/root/lua-bundler"]
CMD ["-help"]