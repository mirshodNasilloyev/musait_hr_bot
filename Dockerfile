
FROM golang:1.23 AS builder

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o myapp .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/myapp .
COPY --from=builder /app/configs /app/configs

EXPOSE 80

CMD ["./myapp"]
