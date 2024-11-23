# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy backend source code
COPY . .

# Create necessary directories
RUN mkdir -p /app/data/static /app/internal/database/migrations/files

# Build for production
RUN go build -o main ./cmd/server

# Development stage with hot-reloading
FROM golang:1.23 AS development

WORKDIR /app

# Install Air
RUN go install github.com/air-verse/air@latest

# Create .air.toml configuration
RUN echo '# .air.toml\n\
root = "."\n\
tmp_dir = "tmp"\n\
[build]\n\
cmd = "go build -o ./tmp/main ./cmd/server"\n\
bin = "tmp/main"\n\
include_ext = ["go", "tpl", "tmpl", "html", "env", "sql", "json"]\n\
exclude_dir = ["assets", "tmp", "vendor"]\n\
delay = 1000\n\
stop_on_error = true\n\
send_interrupt = false\n\
kill_delay = 500\n\
[log]\n\
time = true\n\
[color]\n\
main = "magenta"\n\
watcher = "cyan"\n\
build = "yellow"\n\
runner = "green"\n\
[misc]\n\
clean_on_exit = true' > .air.toml

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

EXPOSE 8080

# Use Air for development
CMD ["air", "-c", ".air.toml"]

# Production stage
FROM golang:1.23 AS production

WORKDIR /app

# Copy binary and required files from builder
COPY --from=builder /app/main .
COPY --from=builder /app/data/static ./data/static
COPY --from=builder /app/internal/database/migrations/files ./internal/database/migrations/files

EXPOSE 8080

# Production command
CMD ["./main"]