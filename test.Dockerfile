FROM golang:tip-alpine3.21 AS go_build

RUN apk add build-base git

ARG GH_USER
ARG GH_TOKEN

# Copy libs from prebuilt swe-builder image (build this image first)
COPY --from=swe-builder:latest /usr/local/lib/libswe.so /usr/local/lib/libswe.so
COPY --from=swe-builder:latest /usr/local/lib/ephe /usr/local/lib/ephe

ENV SWISSEPH_PATH=/usr/local/lib/ephe
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOPRIVATE=github.com/jenujari/*
ENV GONOSUMDB=github.com/jenujari/*

WORKDIR /test_sweisseph

COPY ./go.mod  ./
COPY ./go.sum  ./

RUN git config --global url."https://$GH_USER:$GH_TOKEN@github.com/".insteadOf "https://github.com/" 

RUN go mod download

ENTRYPOINT ["go", "test", "./..."]
