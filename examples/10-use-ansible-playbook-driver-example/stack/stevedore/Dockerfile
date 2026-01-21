# Based on https://github.com/timberio/docker-client/blob/master/Dockerfile
ARG GOLANG_VERSION="1.24"
ARG DOCKER_VERSION="28.5"

## docker dind image is used to copy its entrypoint
FROM docker:${DOCKER_VERSION}-dind AS dind

# Compiling Stevedore binary
FROM golang:${GOLANG_VERSION}-alpine AS golang

WORKDIR /usr/src/app
RUN apk add --no-cache make build-base

COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -ldflags "-s -w -X 'github.com/gostevedore/stevedore/internal/core/domain/release.BuildDate=$(date +%c)' -X 'github.com/gostevedore/stevedore/internal/core/domain/release.Version=example' -X 'github.com/gostevedore/stevedore/internal/core/domain/release.Commit=unknown'" -v -o /usr/local/bin/stevedore ./cmd/stevedore.go

# Define an image with Stevedore and Ansible
FROM debian:bookworm

ARG DEBIAN_FRONTEND=noninteractive
ARG ansible_version=5.9.0
ARG ansible_stdout_callback=yaml
ARG user_id=10001
ARG working_dir=/src
ARG example_name

## Installing packages
RUN apt-get update \
  && apt-get install --no-install-recommends --yes \
    build-essential \
    curl \
    git \
    libffi-dev \
    libssl-dev \
    openssh-client \
    python3 \
    python3-dev \
    python3-pip \
  && curl -sSL https://get.docker.com/ | sh \
  && apt-get clean \
  && apt-get autoremove --yes \
  && rm -rf /var/lib/apt/lists/* \
  && ln -s /usr/bin/python3 /usr/bin/python

# Install Ansible via pip.
RUN pip3 install --break-system-packages \
  ansible==${ansible_version} \
  cryptography \
  docker \
  requests

COPY --from=golang /usr/local/bin/stevedore /usr/local/bin/stevedore
COPY --from=dind /usr/local/bin/dockerd-entrypoint.sh /usr/local/bin/dockerd-entrypoint.sh
COPY examples/${example_name}/stack/stevedore/entrypoint.sh /usr/local/bin/entrypoint.sh
COPY examples/${example_name}/stack/stevedore/wait-for-dockerd.sh /usr/local/bin/wait-for-dockerd.sh

ENV ANSIBLE_FORCE_COLOR=true
ENV ANSIBLE_HOST_KEY_CHECKING=False
ENV ANSIBLE_STDOUT_CALLBACK=${ansible_stdout_callback}
## Setting PYTHONUNBUFFERED to a non-empty value different from 0 ensures that the python output i.e. the stdout and stderr streams are sent straight to terminal (e.g. your container log) without being first buffered and that you can see the output of your application in real time.
ENV PYTHONUNBUFFERED=1

WORKDIR ${working_dir}

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["stevedore","version"]
HEALTHCHECK --interval=15s --timeout=2s --start-period=5s --retries=3 CMD nc -vz localhost 2376 || exit 1
