# Based on https://github.com/timberio/docker-client/blob/master/Dockerfile
ARG GOLANG_VERSION="1.19"
ARG DOCKER_VERSION="24.0"

FROM golang:${GOLANG_VERSION}-alpine AS golang

WORKDIR /usr/src/app
RUN apk add --no-cache make build-base
COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

COPY . ./
RUN go build -ldflags "-s -w" -v -o /usr/local/bin/stevedore ./cmd/stevedore.go

FROM docker:${DOCKER_VERSION}-dind

COPY --from=golang /usr/local/go /usr/local/go
COPY --from=golang /usr/local/bin/stevedore /usr/local/bin/stevedore

ENV GOPATH /go
ENV PATH /usr/local/go/bin:$PATH
ENV PATH $GOPATH/bin:$PATH 

WORKDIR $GOPATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

COPY test/stack/client/entrypoint.sh /usr/local/bin/entrypoint.sh
COPY test/stack/client/wait-for-dockerd.sh /usr/local/bin/wait-for-dockerd.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["go","version"]
HEALTHCHECK --interval=15s --timeout=2s --start-period=5s --retries=3 CMD nc -vz localhost 2376 || exit 1
