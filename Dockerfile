# The certs stage is used to obtain a current set of CA certificates.
FROM docker.io/library/alpine:3.19 AS deps

# hadolint ignore=DL3018
RUN apk add --no-cache \
    ca-certificates \
    docker-cli

# The builder build stage compiles the Go code into a static binary.
FROM golang:1.21-alpine as build

WORKDIR /go/src/github.com/joshdk/actions-docker-shim

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -o /bin/actions-docker-shim \
    -buildvcs=false \
    -ldflags "-buildid= -s -w" \
    -trimpath \
    .

# The final build stage copies in the compiled binary.
FROM scratch

ARG CREATED
ARG REVISION
ARG VERSION

# hadolint ignore=DL4000
MAINTAINER Josh Komoroske <github.com/joshdk>

# Standard OCI image labels.
# See: https://github.com/opencontainers/image-spec/blob/v1.0.1/annotations.md#pre-defined-annotation-keys
LABEL org.opencontainers.image.created="$CREATED"
LABEL org.opencontainers.image.authors="Josh Komoroske <github.com/joshdk>"
LABEL org.opencontainers.image.url="https://github.com/joshdk/actions-docker-shim"
LABEL org.opencontainers.image.documentation="https://github.com/joshdk/actions-docker-shim/blob/master/README.md"
LABEL org.opencontainers.image.source="https://github.com/joshdk/actions-docker-shim"
LABEL org.opencontainers.image.version="$VERSION"
LABEL org.opencontainers.image.revision="$REVISION"
LABEL org.opencontainers.image.vendor="Josh Komoroske <github.com/joshdk>"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.ref.name="ghcr.io/joshdk/actions-docker-shim:$VERSION"
LABEL org.opencontainers.image.title="actions-docker-shim"
LABEL org.opencontainers.image.description="Shim that enables using private ghcr.io images in GitHub Actions"

COPY LICENSE.txt /
COPY --from=deps  /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=deps  /usr/bin/docker /usr/bin/docker
COPY --from=deps  /lib/ld-musl-x86_64.so.1 /lib/ld-musl-x86_64.so.1
COPY README.md /
COPY --from=build /bin/actions-docker-shim /bin/actions-docker-shim

ENTRYPOINT ["/bin/actions-docker-shim"]
