#!/bin/sh
set -e

GOLANG_VERSION=1.19
ALPINE_VERSION=3.16
UBUNTU_VERSION=20.04

docker pull golang:${GOLANG_VERSION}-alpine
docker tag golang:${GOLANG_VERSION}-alpine docker-hub.stevedore.test:5000/library/golang:${GOLANG_VERSION}-alpine
docker push docker-hub.stevedore.test:5000/library/golang:${GOLANG_VERSION}-alpine

docker pull busybox
docker tag busybox docker-hub.stevedore.test:5000/library/busybox
docker push docker-hub.stevedore.test:5000/library/busybox

docker pull alpine:${ALPINE_VERSION}
docker tag alpine:${ALPINE_VERSION} docker-hub.stevedore.test:5000/library/alpine:${ALPINE_VERSION}
docker push docker-hub.stevedore.test:5000/library/alpine:${ALPINE_VERSION}

docker pull ubuntu:${UBUNTU_VERSION}
docker tag ubuntu:${UBUNTU_VERSION} docker-hub.stevedore.test:5000/library/ubuntu:${UBUNTU_VERSION}
docker push docker-hub.stevedore.test:5000/library/ubuntu:${UBUNTU_VERSION}

docker pull ubuntu:latest
docker tag ubuntu:latest docker-hub.stevedore.test:5000/library/ubuntu:latest
docker push docker-hub.stevedore.test:5000/library/ubuntu:latest
