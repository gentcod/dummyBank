# Build stage
FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Mid stage
FROM alpine:3.18 AS alpine
WORKDIR /app
COPY --from=builder /app/main .

# Database migration stage
FROM gomicro/goose as migration
WORKDIR /app
COPY --from=alpine /app/main .
COPY app.env .
COPY sql/schemas /app/sql/schemas
ADD start.sh /app/sql/schemas
RUN chmod +x /app/sql/schemas/start.sh

EXPOSE 5000
ENTRYPOINT [ "/app/sql/schemas/start.sh" ]
CMD ["/app/main"]