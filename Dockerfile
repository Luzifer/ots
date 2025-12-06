FROM golang:1.25.4-alpine@sha256:d3f0cf7723f3429e3f9ed846243970b20a2de7bae6a5b66fc5914e228d831bbb AS builder

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


FROM alpine:3.23@sha256:51183f2cfa6320055da30872f211093f9ff1d3cf06f39a0bdb212314c5dc7375

LABEL org.opencontainers.image.authors='Knut Ahlers <knut@ahlers.me>' \
    org.opencontainers.image.version='1.20.0' \
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
