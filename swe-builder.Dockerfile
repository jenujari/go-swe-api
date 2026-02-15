FROM alpine:3.23

RUN apk add --no-cache build-base git

RUN mkdir -p /usr/local/lib
ENV LD_LIBRARY_PATH=/usr/local/lib

WORKDIR /swisseph
RUN git clone https://github.com/aloistr/swisseph.git .
RUN make libswe.so

RUN cp libswe.so /usr/local/lib/libswe.so && cp -r ephe /usr/local/lib/ephe

CMD ["/bin/sh"]
