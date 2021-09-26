FROM luzifer/archlinux as builder

ENV CGO_ENABLED=0 \
    GOPATH=/go

COPY . /go/src/github.com/Luzifer/ots
WORKDIR /go/src/github.com/Luzifer/ots

RUN set -ex \
 && pacman --noconfirm -Syy \
      curl \
      git \
      go \
      make \
      nodejs-lts-fermium \
      npm \
      tar \
      unzip \
 && make -C src -f ../Makefile generate-inner \
 && make download_libs generate-apidocs \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly


FROM alpine:latest

LABEL maintainer "Knut Ahlers <knut@ahlers.me>"

RUN set -ex \
 && apk --no-cache add \
      ca-certificates

COPY --from=builder /go/bin/ots /usr/local/bin/ots

EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/ots"]
CMD ["--"]

# vim: set ft=Dockerfile:
