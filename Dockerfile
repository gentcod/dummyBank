# Build stage
FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.18 AS alpine
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=alpine /app/main .
COPY app.env .
COPY sql/schemas /app/sql/schemas
ADD start.sh .
RUN chmod +x /app/sql/schemas/start.sh

EXPOSE 5000
ENTRYPOINT [ "/app/start.sh" ]
CMD ["/app/main"]