#!/bin/sh

keys=${GIT_KEYS:-/git/keys}
git_ssh_folder="${GIT_SSH_FOLDER:-/home/git/.ssh}"

find "${keys}" -type f -name '*.pub' | while read key
do
    echo "Loading key: $key"
    cat "${key}" >> "${git_ssh_folder}/authorized_keys"
done

spawn-fcgi -U git -G git -s /var/run/fcgiwrap.socket /usr/bin/fcgiwrap &
nginx -g "daemon off;" &

/usr/sbin/sshd -D
