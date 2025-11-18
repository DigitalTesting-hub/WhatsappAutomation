FROM golang:1.24-alpine AS builder

# Install build dependencies including gcc for CGO
RUN apk add --no-cache git ffmpeg gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files from src directory
COPY src/go.mod src/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY src/ ./

# Build the application with CGO enabled
ENV CGO_ENABLED=1
RUN go build -o whatsapp .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ffmpeg ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/whatsapp .

# Copy necessary files (if they exist)
COPY --from=builder /app/views ./views 2>/dev/null || true
COPY --from=builder /app/statics ./statics 2>/dev/null || true

# Create storage directory
RUN mkdir -p /app/storages

# Expose port
EXPOSE 3000

# Run the application in REST mode
CMD ["./whatsapp", "rest"]
