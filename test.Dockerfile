FROM dhi.io/golang:1.27.0-alpine3.24-dev

RUN apk add --no-cache build-base git

COPY --from=docker.io/jhon5456/sweph-build-base:v1 /usr/local/lib/libswe.so /usr/local/lib/libswe.so
COPY --from=docker.io/jhon5456/sweph-build-base:v1 /usr/local/lib/ephe /usr/local/lib/ephe

ENV SWISSEPH_PATH=/usr/local/lib/ephe
ENV CGO_ENABLED=1
ENV LD_LIBRARY_PATH=/usr/local/lib

WORKDIR /test_sweisseph

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENTRYPOINT ["go", "test", "./...", "-v"]
