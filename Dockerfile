FROM alpine:3.18

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY bin/ogbuser-linux-amd64 /app/ogbuser
COPY user-config.yaml /app/user-config.yaml

RUN chmod +x /app/ogbuser && \
    chown -R appuser:appgroup /app

USER appuser

EXPOSE 12121
EXPOSE 12122

CMD ["/app/ogbuser", "serve", "--config", "/app/user-config.yaml"]
