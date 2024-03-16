# Build stage
FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.18 AS alpine
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY sql/migrations /app/sql/migrations
COPY start.sh .
RUN chmod +x /app/start.sh

EXPOSE 5000
CMD ["/app/main"]