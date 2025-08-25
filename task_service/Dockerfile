FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/. .

RUN CGO_ENABLED=0 go build -o /task-service ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /task-service /app/task-service

COPY src/app.env .
COPY src/migrations ./migrations

EXPOSE 8081

CMD ["/app/task-service"]
