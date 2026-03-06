FROM golang:1.25.6-alpine3.23 AS builder

RUN apk add --no-cache build-base git protobuf protobuf-dev

# Install go plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Copy libs from prebuilt swe-builder image (build this image first)
COPY --from=swe-builder:latest /usr/local/lib/libswe.so /usr/local/lib/libswe.so
COPY --from=swe-builder:latest /usr/local/lib/ephe /usr/local/lib/ephe

ENV SWISSEPH_PATH=/usr/local/lib/ephe
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV LD_LIBRARY_PATH=/usr/local/lib

WORKDIR /test_sweisseph

COPY ./go.mod  ./
COPY ./go.sum  ./

RUN go mod download

COPY . .

# Generate proto code before testing
RUN protoc -Iproto --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative proto/swe.proto

ENTRYPOINT ["go", "test", "./..." , "-v"]
