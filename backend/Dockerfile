# Build stage
FROM golang:1.21 AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy backend source code
COPY . .

# Create necessary directories
RUN mkdir -p /app/data/static /app/internal/database/migrations/files

# Build for production
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Production stage
FROM alpine:latest

WORKDIR /app

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy binary and required files from builder
COPY --from=builder /app/main .
COPY --from=builder /app/data/static ./data/static
COPY --from=builder /app/internal/database/migrations/files ./internal/database/migrations/files

EXPOSE 8080

# Production command
CMD ["./main"]