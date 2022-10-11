# Based on https://github.com/timberio/docker-client/blob/master/Dockerfile
ARG GOLANG_VERSION="1.19"
ARG DOCKER_VERSION="20.10"

FROM golang:${GOLANG_VERSION}-alpine AS golang

WORKDIR /usr/src/app

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

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["go","version"]
