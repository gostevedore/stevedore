# Based on https://github.com/timberio/docker-client/blob/master/Dockerfile
ARG GOLANG_VERSION="1.24"
ARG DOCKER_VERSION="27.1"

FROM golang:${GOLANG_VERSION}-alpine AS golang

FROM docker:${DOCKER_VERSION}-dind

ARG GO_TOOL_ERRCHECK_VERSION="latest"
ARG GO_TOOL_GOLANGCI_LINT_VERSION="latest"
ARG GO_TOOL_GOLINT_VERSION="latest"
ARG GO_TOOL_GOSEC_VERSION="latest"
ARG GO_TOOL_STATICCHECK_VERSION="latest"

COPY --from=golang /usr/local/go /usr/local/go

WORKDIR /usr/src/app
RUN apk add --update --no-cache \
    build-base \
    docker-cli-compose \
    git

ENV GOPATH /go
ENV PATH /usr/local/go/bin:$PATH
ENV PATH $GOPATH/bin:$PATH 
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# Install errcheck
RUN go install github.com/kisielk/errcheck@${GO_TOOL_ERRCHECK_VERSION}

# Install golangci-lint
# Upgrade to golang 1.21 before using it
# RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GO_TOOL_GOLANGCI_LINT_VERSION}

# Install golint
RUN go install golang.org/x/lint/golint@${GO_TOOL_GOLINT_VERSION}

# Install gosec
RUN go install github.com/securego/gosec/cmd/gosec@${GO_TOOL_GOSEC_VERSION}

# Install staticcheck
RUN go install honnef.co/go/tools/cmd/staticcheck@${GO_TOOL_STATICCHECK_VERSION}

COPY build/entrypoint.sh /usr/local/bin/entrypoint.sh
COPY build/wait-for-dockerd.sh /usr/local/bin/wait-for-dockerd.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
HEALTHCHECK --interval=15s --timeout=2s --start-period=5s --retries=3 CMD docker info || exit 1
