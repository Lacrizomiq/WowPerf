FROM golang:1.22.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY data/static ./static

RUN go build -o main ./cmd/server

EXPOSE 8080

CMD ["./main"]
