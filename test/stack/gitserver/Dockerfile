FROM alpine:3.16

ARG git_dir="/git"
ARG repos_dir="${git_dir}/repos"
ARG keys_dir="${git_dir}/keys"
ARG git_user_home="/home/git"

ENV GIT_DIR=${git_dir}
ENV REPOS_DIR=${repos_dir}
ENV KEYS_DIR=${keys_dir}

WORKDIR ${git_dir}

RUN apk add --no-cache \
        openssh \
        git \
    && rm -rf /var/cache/apk/* \
    && ssh-keygen -A \
    && adduser -D --home ${git_user_home} --shel /usr/bin/git-shell git \ 
    # user password is required in order to access via ssh
    && echo git:$(xxd -l4 /dev/random -p) | chpasswd \
    && mkdir -p ${repos_dir} \
    && mkdir -p ${keys_dir} \
    && mkdir -p ${git_user_home}/.ssh \
    && touch ${git_user_home}/.ssh/authorized_keys \
    && chown -R git:git ${git_user_home} \
    && chmod 700 ${git_user_home}/.ssh \
    && chmod 600 ${git_user_home}/.ssh/authorized_keys \
    && sed -i "s/.*PubkeyAuthentication.*/PubkeyAuthentication yes/g" /etc/ssh/sshd_config \
    && sed -i "s/.*PasswordAuthentication.*/PasswordAuthentication no/g" /etc/ssh/sshd_config

RUN apk add  --no-cache \
        nginx \
        git-daemon \
        fcgiwrap \
        spawn-fcgi \
    && rm -rf /var/cache/apk/*

COPY files/nginx.conf /etc/nginx/nginx.conf

EXPOSE 22
EXPOSE 80

COPY entrypoint.sh /entrypoint.sh
CMD ["/entrypoint.sh"]