FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o pr_chat_bot ./cmd/pr_chat_bot

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/pr_chat_bot .

CMD ["./pr_chat_bot"]
