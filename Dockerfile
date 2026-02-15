FROM golang:1.25.6-alpine3.23 AS builder

RUN apk add --no-cache build-base protobuf protobuf-dev

# Install go plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Copy libs from prebuilt swe-builder image (build this image first)
COPY --from=swe-builder:latest /usr/local/lib/libswe.so /usr/local/lib/libswe.so
COPY --from=swe-builder:latest /usr/local/lib/ephe /usr/local/lib/ephe

ENV SWISSEPH_PATH=/usr/local/lib/ephe
ENV CGO_ENABLED=1
ENV LD_LIBRARY_PATH=/usr/local/lib

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate proto code before building
RUN protoc -Iproto --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative proto/swe.proto

RUN go build -o sweAPI main.go

FROM alpine:3.23
RUN apk add --no-cache libc6-compat

# Need to copy shared libs to the final image too
COPY --from=swe-builder:latest /usr/local/lib/libswe.so /usr/local/lib/libswe.so
COPY --from=swe-builder:latest /usr/local/lib/ephe /usr/local/lib/ephe

ENV SWISSEPH_PATH=/usr/local/lib/ephe
ENV LD_LIBRARY_PATH=/usr/local/lib

WORKDIR /app
COPY --from=builder /app/sweAPI .
COPY config/conf.yml config/conf.yml

EXPOSE 5678
ENTRYPOINT ["./sweAPI"]
