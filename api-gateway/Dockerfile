FROM golang:1.25-alpine3.22 AS dependencies
WORKDIR /app
COPY go.mod go.sum ./

RUN apk add --no-cache wget

FROM dependencies AS builder
WORKDIR /app
COPY . .
RUN go build -o api ./cmd

FROM alpine:3.23.3
WORKDIR /app
COPY --from=builder /app /app
EXPOSE 8080
ENTRYPOINT [ "/app/api" ]