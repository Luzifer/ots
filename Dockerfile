FROM docker.io/library/golang:1.27.1-alpine@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414 AS builder

COPY --from=ghcr.io/luzifer-docker/pnpm:v11.28.1@sha256:438ad48e27d96671927274122bde7549d78334f13bf02a16d722d01d52c1b3f6 . /

ENV CGO_ENABLED=0 \
    GOPATH=/go \
    NODE_ENV=production

COPY . /src/ots
WORKDIR /src/ots

RUN set -ex \
 && apk --no-cache add \
      curl \
      git \
      make \
      nodejs-current \
      npm \
      tar \
      unzip \
 && make build-local


FROM scratch

LABEL org.opencontainers.image.authors='Knut Ahlers <knut@ahlers.me>' \
    org.opencontainers.image.version='1.22.0' \
    org.opencontainers.image.url='https://github.com/Luzifer/ots/pkgs/container/ots' \
    org.opencontainers.image.documentation='https://github.com/Luzifer/ots/wiki' \
    org.opencontainers.image.source='https://github.com/Luzifer/ots' \
    org.opencontainers.image.licenses='Apache-2.0'

COPY --from=builder /src/ots/ots /usr/local/bin/ots

EXPOSE 3000

USER 1000:1000

ENTRYPOINT ["/usr/local/bin/ots"]
