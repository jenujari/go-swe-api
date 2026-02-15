FROM alpine:3.21
RUN apk add --no-cache build-base git

WORKDIR /swisseph
RUN git clone https://github.com/aloistr/swisseph.git .
RUN make libswe.so

RUN mkdir -p /usr/local/lib \
  && cp libswe.so /usr/local/lib/libswe.so \
  && cp -r ephe /usr/local/lib/ephe

CMD ["/bin/sh"]
