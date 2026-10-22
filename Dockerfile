FROM golang:1.27.3-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/nginx-prometheus .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && adduser -S app -G app

COPY --from=builder /out/nginx-prometheus /usr/bin/nginx-prometheus

EXPOSE 1514/udp 1514/tcp 11080/tcp

USER app
ENTRYPOINT ["/usr/bin/nginx-prometheus"]
