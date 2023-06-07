#!/bin/sh
set -eu

# allow docker client to interact with registry.stevedore.test
if [ -f /ssl/stevedore.test.crt ]; then
    mkdir -p /etc/docker/certs.d/registry.stevedore.test
    cp /ssl/stevedore.test.crt /etc/docker/certs.d/registry.stevedore.test/ca.crt
fi

/usr/local/bin/dockerd-entrypoint.sh 2> /dev/null &
/usr/local/bin/wait-for-dockerd.sh
# round=50
# # while ! nc -z localhost 2376; do
# while ! docker info > /dev/null 2>&1 ; do
#     if [ $round -eq 0 ]; then
#         >&2 echo "dockerd is not ready"
#         exit 1
#     fi

#     >&2 echo " Waiting for dockerd to be ready..."
#     sleep 0.5 # wait for 1/2 of the second before check again
#     round=$(( round - 1 ))
# done

exec "$@"
