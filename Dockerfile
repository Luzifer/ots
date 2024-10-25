FROM golang:1-alpine AS builder

ENV CGO_ENABLED=0 \
    GOPATH=/go \
    NODE_ENV=production

COPY . /go/src/github.com/Luzifer/ots
WORKDIR /go/src/github.com/Luzifer/ots

RUN set -ex \
 && apk update && apk add \
      curl \
      git \
      make \
      nodejs-lts \
      npm \
      tar \
      unzip \
 && make download_libs generate-inner generate-apidocs \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly


FROM alpine:latest

LABEL org.opencontainers.image.authors='Knut Ahlers <knut@ahlers.me>' \
    org.opencontainers.image.version='1.13.0' \
    org.opencontainers.image.url='https://github.com/Luzifer/ots/pkgs/container/ots' \
    org.opencontainers.image.documentation='https://github.com/Luzifer/ots/wiki' \
    org.opencontainers.image.source='https://github.com/Luzifer/ots' \
    org.opencontainers.image.licenses='Apache-2.0'

RUN set -ex \
 && apk --no-cache add \
      ca-certificates

COPY --from=builder /go/bin/ots /usr/local/bin/ots

EXPOSE 3000

USER 1000:1000

ENTRYPOINT ["/usr/local/bin/ots"]
CMD ["--"]

# vim: set ft=Dockerfile:
