FROM alpine:3.23

ARG SWISSEPH_REF=master
ENV LD_LIBRARY_PATH=/usr/local/lib

WORKDIR /swisseph

RUN apk add --no-cache build-base git \
    && git clone --depth 1 --branch "${SWISSEPH_REF}" https://github.com/aloistr/swisseph.git . \
    && make libswe.so \
    && install -Dm755 libswe.so /usr/local/lib/libswe.so \
    && cp -r ephe /usr/local/lib/ephe
