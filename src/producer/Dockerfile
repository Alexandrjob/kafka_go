FROM golang:1.24.2 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /producer && \
    chmod +x /producer

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -u 1000 appuser
COPY --from=builder --chown=appuser:appuser /producer /producer
USER appuser
CMD ["/producer"]