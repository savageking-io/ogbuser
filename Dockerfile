FROM golang:1.23.4-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.AppVersion=$(cat VERSION)" -o /app/ogbuser .

FROM alpine:3.18

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/ogbuser /app/ogbuser
COPY user-config.yaml /app/user-config.yaml

RUN chmod +x /app/ogbuser && \
    chown -R appuser:appgroup /app

USER appuser

EXPOSE 8090

CMD ["/app/ogbuser", "serve", "--config", "/app/user-config.yaml"]
