# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and static files
COPY . .
COPY data/static ./static

# Build for production
RUN go build -o main ./cmd/server

# Development stage with hot-reloading
FROM golang:1.23 AS development

WORKDIR /app

# Install Air avec le nouveau chemin
RUN go install github.com/air-verse/air@latest

# Create .air.toml configuration
RUN echo '# .air.toml\n\
root = "."\n\
tmp_dir = "tmp"\n\
[build]\n\
cmd = "go build -o ./tmp/main ./cmd/server"\n\
bin = "tmp/main"\n\
include_ext = ["go", "tpl", "tmpl", "html", "env"]\n\
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

# Copy source code and static files
COPY . .
COPY data/static ./static

EXPOSE 8080

# Use Air for development
CMD ["air", "-c", ".air.toml"]

# Production stage
FROM golang:1.23 AS production

WORKDIR /app

# Copy binary and static files from builder
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static

EXPOSE 8080

# Production command
CMD ["./main"]