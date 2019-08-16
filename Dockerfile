FROM alpine:3.7

## copy compiled binary
COPY bin/docker-namejoked /namejoked
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

## keep web rest port open
EXPOSE 8082
ENTRYPOINT ["/namejoked", "--loglevel", "debug"]
# CMD "/namejoked --loglevel debug"