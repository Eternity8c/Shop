FROM golang:1.25-alpine3.22 AS dependencies
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && apk add --no-cache wget

FROM dependencies AS builder
WORKDIR /app
COPY . .
RUN go build -o  auth ./cmd/auth

FROM alpine:3.23.3
WORKDIR /app
COPY --from=builder /app/auth /app/auth

ENTRYPOINT [ "/app/auth" ]