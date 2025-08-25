# build
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/. .

RUN CGO_ENABLED=0 go build -o /notification-service ./cmd/server

# final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /notification-service /app/notification-service

COPY src/app.env .

CMD ["/app/notification-service"]
