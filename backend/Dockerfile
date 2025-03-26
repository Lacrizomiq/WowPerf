# Build stage
FROM golang:1.23.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build optimized for production
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main ./cmd/server

# Development stage
FROM golang:1.23.3-alpine AS development

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

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]

# Production stage
FROM alpine:3.18 AS production

WORKDIR /app

RUN apk add --no-cache tzdata ca-certificates

COPY --from=builder /app/main .
COPY --from=builder /app/data/static ./data/static
COPY --from=builder /app/internal/database/migrations/files ./internal/database/migrations/files

RUN adduser -D -H -h /app appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["./main"]