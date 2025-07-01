FROM golang:1.24 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /matching-api ./cmd/main.go

FROM scratch
COPY --from=builder /matching-api /matching-api

ENV APP_PORT=<port> \
    JWT_SECRET=<super-secret> \
    DRIVER_API_BASE_URL=<Driver Location API Base URL> \
    DRIVER_API_TOKEN=<auth token> \
    BREAKER_TIMEOUT=5s \

EXPOSE 9000
HEALTHCHECK --interval=30s --timeout=3s CMD \
  wget -qO- http://localhost:${APP_PORT}/v1/healthz || exit 1

ENTRYPOINT ["/matching-api"]
