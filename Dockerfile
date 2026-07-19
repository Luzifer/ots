FROM docker.io/library/golang:1.26.5-alpine@sha256:0178a641fbb4858c5f1b48e34bdaabe0350a330a1b1149aabd498d0699ff5fb2 AS builder

COPY --from=ghcr.io/luzifer-docker/pnpm:v11.13.1@sha256:0230302ad1be2918b4690f236730150ca023b3e30604519d03ebae555393f81d . /

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
    org.opencontainers.image.version='1.21.8' \
    org.opencontainers.image.url='https://github.com/Luzifer/ots/pkgs/container/ots' \
    org.opencontainers.image.documentation='https://github.com/Luzifer/ots/wiki' \
    org.opencontainers.image.source='https://github.com/Luzifer/ots' \
    org.opencontainers.image.licenses='Apache-2.0'

COPY --from=builder /src/ots/ots /usr/local/bin/ots

EXPOSE 3000

USER 1000:1000

ENTRYPOINT ["/usr/local/bin/ots"]
