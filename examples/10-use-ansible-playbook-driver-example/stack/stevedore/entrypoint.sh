#!/bin/sh
set -eu

# allow docker client to interact with registry.stevedore.test
if [ -f /ssl/stevedore.test.crt ]; then
    mkdir -p /etc/docker/certs.d/registry.stevedore.test
    cp /ssl/stevedore.test.crt /etc/docker/certs.d/registry.stevedore.test/ca.crt
fi

/usr/local/bin/dockerd-entrypoint.sh 2> /dev/null &
/usr/local/bin/wait-for-dockerd.sh

exec "$@"
