FROM golang:1.24.5-alpine@sha256:48ee313931980110b5a91bbe04abdf640b9a67ca5dea3a620f01bacf50593396 AS builder

ENV CGO_ENABLED=0 \
    GOPATH=/go \
    NODE_ENV=production

COPY . /go/src/github.com/Luzifer/ots
WORKDIR /go/src/github.com/Luzifer/ots

RUN set -ex \
 && apk --no-cache add \
      curl \
      git \
      make \
      nodejs-current \
      npm \
      tar \
      unzip \
 && make frontend_prod generate-apidocs \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly


FROM alpine:3.22@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1

LABEL org.opencontainers.image.authors='Knut Ahlers <knut@ahlers.me>' \
    org.opencontainers.image.version='1.17.2' \
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
