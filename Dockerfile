FROM docker.io/jhon5456/sweph-build-base:v2.10.03 AS swe

FROM dhi.io/golang:1.27.0-alpine3.24-dev@sha256:a078c40b4db93bdddad7b24137402fad8caf5b8944776c2f2d90756fa2b9596b AS builder

RUN apk add --no-cache build-base

COPY --from=swe /usr/local/lib/libswe.so /usr/local/lib/libswe.so
COPY --from=swe /usr/local/lib/ephe /usr/local/lib/ephe

ENV SWISSEPH_PATH=/usr/local/lib/ephe
ENV CGO_ENABLED=1
ENV LD_LIBRARY_PATH=/usr/local/lib
ENV CGO_LDFLAGS="-L/usr/local/lib -Wl,-rpath,/usr/local/lib"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -trimpath -ldflags="-s -w" -o /app/sweAPI . \
    && /lib/ld-musl-$(uname -m).so.1 --list /app/sweAPI

FROM dhi.io/alpine-base:3.24@sha256:aa2aa13f40cfc9e17296ba19e64d9439f040fc5f01e4b29d81661ef7063e4e46

COPY --from=swe /usr/local/lib/libswe.so /usr/local/lib/libswe.so
COPY --from=swe /usr/local/lib/ephe /usr/local/lib/ephe

ENV SWISSEPH_PATH=/usr/local/lib/ephe
ENV LD_LIBRARY_PATH=/usr/local/lib

WORKDIR /app
COPY --from=builder --chown=65532:65532 /app/sweAPI /app/sweAPI
COPY --chown=65532:65532 config/conf.yml /app/config/conf.yml

EXPOSE 5678
USER 65532
ENTRYPOINT ["/app/sweAPI"]
