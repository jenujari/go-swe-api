FROM golang:1.25.6-alpine3.23 AS builder

RUN apk add --no-cache build-base

COPY --from=docker.io/jhon5456/sweph-build-base:v1 /usr/local/lib/libswe.so /usr/local/lib/libswe.so
COPY --from=docker.io/jhon5456/sweph-build-base:v1 /usr/local/lib/ephe /usr/local/lib/ephe

ENV SWISSEPH_PATH=/usr/local/lib/ephe
ENV CGO_ENABLED=1
ENV LD_LIBRARY_PATH=/usr/local/lib

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o sweAPI main.go

FROM alpine:3.23

COPY --from=docker.io/jhon5456/sweph-build-base:v1 /usr/local/lib/libswe.so /usr/local/lib/libswe.so
COPY --from=docker.io/jhon5456/sweph-build-base:v1 /usr/local/lib/ephe /usr/local/lib/ephe

ENV SWISSEPH_PATH=/usr/local/lib/ephe
ENV LD_LIBRARY_PATH=/usr/local/lib

WORKDIR /app
COPY --from=builder /app/sweAPI ./
COPY config/conf.yml config/conf.yml

EXPOSE 5678
ENTRYPOINT ["./sweAPI"]
