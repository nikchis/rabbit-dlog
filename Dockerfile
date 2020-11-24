FROM golang:1.14.12-buster as builder
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY . /build/
WORKDIR /build
ENV CGO_ENABLED=0
ARG TARGETOS
ARG TARGETARCH
RUN go get -v -t -d ./...
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
-a -ldflags '-extldflags "-static"' -o main .
FROM scratch
COPY --from=builder /build/main /app/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
WORKDIR /app/data
ENTRYPOINT ["/app/main", "-dir", "/app/data"]
