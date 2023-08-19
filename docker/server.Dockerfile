FROM golang:1.21.0-alpine3.18 AS builder
WORKDIR /app
COPY ../. .
RUN go build -o main ./server/main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]