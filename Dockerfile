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

# Copy go mod files from src directory
COPY src/go.mod src/go.sum ./

# Download dependencies with CGO enabled
ENV CGO_ENABLED=1
ENV GOOS=linux
RUN go mod download

# Copy source code
COPY src/ ./

# Verify CGO is enabled and build
RUN echo "Building with CGO_ENABLED=1..." && \
    go env && \
    go build -v -ldflags="-w -s" -o whatsapp .

# Final stage
FROM alpine:latest

# Install runtime dependencies including sqlite libs
RUN apk add --no-cache \
    ffmpeg \
    ca-certificates \
    sqlite-libs \
    tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/whatsapp .

# Copy necessary files (gracefully handle if they don't exist)
COPY --from=builder /app/views ./views 2>/dev/null || true
COPY --from=builder /app/statics ./statics 2>/dev/null || true

# Create storage directory with proper permissions
RUN mkdir -p /app/storages && chmod 777 /app/storages

# Expose port (Render uses PORT environment variable)
EXPOSE 3000

# Run the application in REST mode
CMD ["./whatsapp", "rest", "--port=3000"]
