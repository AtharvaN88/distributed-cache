FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o cache-server main.go

FROM debian:bullseye-slim

WORKDIR /root/

COPY --from=builder /app/cache-server .

EXPOSE 8081 8082 8083

CMD ["./cache-server"]