FROM golang:tip-alpine3.21 AS go_build

RUN apk add build-base

ARG timestamp=unknown

# Copy libs from prebuilt swe-builder image (build this image first)
COPY --from=swe-builder:latest /usr/local/lib/libswe.so /usr/local/lib/libswe.so
COPY --from=swe-builder:latest /usr/local/lib/ephe /usr/local/lib/ephe

ENV SWISSEPH_PATH=/usr/local/lib/ephe
ENV CGO_ENABLED=1
ENV GOOS=linux

WORKDIR /test_sweisseph

COPY ./sweAPI/go.mod  ./
COPY ./sweAPI/go.sum  ./

RUN go mod download

COPY ./sweAPI ./

RUN go build -o sweAPI main.go

ENTRYPOINT ["./sweAPI"]
