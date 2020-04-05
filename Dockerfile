FROM golang:alpine AS build
WORKDIR /app
COPY . /app
RUN apk add --update-cache yarn make gcc libc-dev && \
    rm -rf /var/cache/apk/* && \
    make dist

FROM alpine:latest
ENV DOIT_DATABASE="/data/doit.db" \
    DOIT_HOST="0.0.0.0" \
    DOIT_PORT="8443" \
    DOIT_MAILFROM="doit@example.com" \
    DOIT_MAILHOST="localhost:25" \
    DOIT_TLS="false" \
    DOIT_TLSCERT="/data/server.crt" \
    DOIT_TLSKEY="/data/server.key"

WORKDIR /
COPY --from=build /app/doit.tar.gz /
COPY entrypoint.sh /
RUN apk add --update-cache openssl && \
    rm -rf /var/cache/apk/* && \
    addgroup -g 1000 doit && \
    adduser -Ss /bin/false -u 1000 -G doit -h /doit doit && \
    tar xvf /doit.tar.gz && \
    rm /doit.tar.gz && \
    mkdir /data && \
    chown 1000:1000 /doit && \
    chmod 755 /entrypoint.sh

EXPOSE $DOIT_PORT

ENTRYPOINT [ "/entrypoint.sh" ]


