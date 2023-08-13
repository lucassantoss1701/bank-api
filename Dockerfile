FROM golang:1.19-alpine as builder

RUN apk add --no-cache git
WORKDIR /build
COPY . .
RUN go build -o server ./cmd/server/main.go

FROM alpine:3.14

RUN apk add --no-cache wget netcat-openbsd bash

WORKDIR /

COPY --from=builder /build/wait-for-services.sh ./
RUN chmod +x wait-for-services.sh

COPY --from=builder /build/internal/infra/database/migrations ./internal/infra/database/migrations

COPY --from=builder /build/server ./server
