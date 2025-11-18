FROM golang:1.24-alpine AS builder

# Install ALL required build dependencies
RUN apk add --no-cache \
    git \
    ffmpeg \
    gcc \
    musl-dev \
    sqlite-dev

# Set working directory
WORKDIR /app

# Copy entire src directory first
COPY src/ ./

# Verify go.mod exists and download dependencies
RUN if [ -f go.mod ]; then \
        echo "go.mod found, downloading dependencies..."; \
        CGO_ENABLED=1 go mod download; \
    else \
        echo "ERROR: go.mod not found!"; \
        ls -la; \
        exit 1; \
    fi

# Build with CGO enabled
ENV CGO_ENABLED=1
ENV GOOS=linux
RUN go build -v -ldflags="-w -s" -o whatsapp .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ffmpeg \
    ca-certificates \
    sqlite-libs \
    tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/whatsapp .

# Create storage directory
RUN mkdir -p /app/storages && chmod 777 /app/storages

# Expose port
EXPOSE 3000

# Run the application in REST mode
CMD ["./whatsapp", "rest", "--port=3000"]
