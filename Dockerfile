# syntax=docker/dockerfile:1
FROM golang:1.22 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o webot cmd/main.go

FROM alpine:latest AS runner

WORKDIR /app
COPY --from=builder /app/webot ./webot

ENV TZ=Asia/Shanghai
RUN apk add --no-cache tzdata

ENTRYPOINT ["/app/webot"]