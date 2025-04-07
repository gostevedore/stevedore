# Based on https://github.com/timberio/docker-client/blob/master/Dockerfile
ARG GOLANG_VERSION="1.23"
ARG DOCKER_VERSION="27.1"

FROM golang:${GOLANG_VERSION}-alpine AS golang

WORKDIR /usr/src/app
RUN apk add --no-cache make build-base

COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -ldflags "-s -w" -v -o /usr/local/bin/stevedore ./cmd/stevedore.go

FROM docker:${DOCKER_VERSION}-dind

COPY --from=golang /usr/local/bin/stevedore /usr/local/bin/stevedore

COPY examples/02-wildcard-version-example/stack/stevedore/entrypoint.sh /usr/local/bin/entrypoint.sh
COPY examples/02-wildcard-version-example/stack/stevedore/wait-for-dockerd.sh /usr/local/bin/wait-for-dockerd.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["stevedore","version"]
HEALTHCHECK --interval=15s --timeout=2s --start-period=5s --retries=3 CMD nc -vz localhost 2376 || exit 1
