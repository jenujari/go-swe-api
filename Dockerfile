FROM dhi.io/golang:1.27.0-alpine3.24-dev AS builder

RUN apk add --no-cache build-base

COPY --from=docker.io/jhon5456/sweph-build-base:v2.10.03 /usr/local/lib/libswe.so /usr/local/lib/libswe.so
COPY --from=docker.io/jhon5456/sweph-build-base:v2.10.03 /usr/local/lib/ephe /usr/local/lib/ephe

ENV SWISSEPH_PATH=/usr/local/lib/ephe
ENV CGO_ENABLED=1
ENV LD_LIBRARY_PATH=/usr/local/lib

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o sweAPI main.go

FROM dhi.io/alpine-base:3.24

COPY --from=docker.io/jhon5456/sweph-build-base:v2.10.03 /usr/local/lib/libswe.so /usr/local/lib/libswe.so
COPY --from=docker.io/jhon5456/sweph-build-base:v2.10.03 /usr/local/lib/ephe /usr/local/lib/ephe

ENV SWISSEPH_PATH=/usr/local/lib/ephe
ENV LD_LIBRARY_PATH=/usr/local/lib

WORKDIR /app
COPY --from=builder /app/sweAPI ./
COPY config/conf.yml config/conf.yml

EXPOSE 5678
ENTRYPOINT ["./sweAPI"]
