FROM golang:alpine as builder
RUN apk update && apk add ca-certificates tzdata && update-ca-certificates
COPY . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build \
-a -ldflags '-extldflags "-static"' -o main .
FROM scratch
COPY --from=builder /build/main /app/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app/data
ENTRYPOINT ["/app/main", "-dir", "/app/data"]
