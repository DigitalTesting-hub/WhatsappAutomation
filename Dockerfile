FROM golang:1.24-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ffmpeg

# Set working directory
WORKDIR /app

# Copy go mod files from src directory
COPY src/go.mod src/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY src/ ./

# Build the application
RUN go build -o whatsapp .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ffmpeg ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/whatsapp .

# Copy necessary files
COPY --from=builder /app/views ./views
COPY --from=builder /app/statics ./statics

# Create storage directory
RUN mkdir -p /app/storages

# Expose port
EXPOSE 3000

# Run the application in REST mode
CMD ["./whatsapp", "rest"]
