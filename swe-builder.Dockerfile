FROM docker.io/library/alpine:3.24@sha256:79ff19e9084a00eece421b2523fb93e22d730e2c0e525905de047e848e56d95f

# libswe ABI must match swephgo's vendored 2.10.03 headers.
ARG SWISSEPH_REF=v2.10.03
# Ephemeris data files are not in the v2.10.03 tag; pin the commit that ships them.
ARG SWISSEPH_EPHE_REF=3fd0f956d73898b91cc4f67cf18b21af656d1342
ENV LD_LIBRARY_PATH=/usr/local/lib

WORKDIR /swisseph

RUN apk add --no-cache build-base git \
    && git clone --depth 1 --branch "${SWISSEPH_REF}" https://github.com/aloistr/swisseph.git . \
    && make libswe.so \
    && gcc -shared -Wl,-soname,libswe.so -o libswe.so \
        swedate.o swehouse.o swejpl.o swemmoon.o swemplan.o sweph.o swephlib.o swecl.o swehel.o \
        -lm -ldl \
    && install -Dm755 libswe.so /usr/local/lib/libswe.so \
    && git init /tmp/swe-ephe \
    && git -C /tmp/swe-ephe remote add origin https://github.com/aloistr/swisseph.git \
    && git -C /tmp/swe-ephe fetch --depth 1 origin "${SWISSEPH_EPHE_REF}" \
    && git -C /tmp/swe-ephe checkout FETCH_HEAD \
    && mkdir -p /usr/local/lib/ephe \
    && cp -a /tmp/swe-ephe/ephe/. /usr/local/lib/ephe/ \
    && cp -f seleapsec.txt /usr/local/lib/ephe/ \
    && rm -rf /tmp/swe-ephe
