# Notes for Stevedore functional test

- [Notes for Stevedore functional test](#notes-for-stevedore-functional-test)
  - [SSL](#ssl)
    - [Generate wildcard self-signed certificate from scratch](#generate-wildcard-self-signed-certificate-from-scratch)
      - [Generate private key and csr](#generate-private-key-and-csr)
      - [Generate certificate](#generate-certificate)
  - [Docker registry](#docker-registry)
    - [docker-hub.stevedore.test](#docker-hubstevedoretest)
      - [How to load a preload docker image](#how-to-load-a-preload-docker-image)
    - [registry.stevedore.test](#registrystevedoretest)
  - [Git](#git)
    - [Users](#users)
    - [How to create a bare repository](#how-to-create-a-bare-repository)
    - [How to clone a gitserver repositories](#how-to-clone-a-gitserver-repositories)

## SSL

### Generate wildcard self-signed certificate from scratch

#### Generate private key and csr

```sh
openssl req -newkey rsa:2048 -nodes -keyout stevedore.test.key -out stevedore.test.csr -config stevedore.test.cnf
```

#### Generate certificate

```sh
openssl x509 -signkey stevedore.test.key -in stevedore.test.csr -req -days 365 -out stevedore.test.crt -extensions req_ext -extfile stevedore.test.cnf
```

## Docker registry

### docker-hub.stevedore.test

`docker-hub.stevedore.test` is the docker-registry which uses the dataset located at `files/docker-hub-storage`. On that dataset are preloaded those images that are used during Stevedore functional test


#### How to load a preload docker image

- Start testing stack

```sh
make attach-client
```

- Download `alpine:3.16` from docker-hub registry

```sh
/app/test/client/stevedore # docker pull alpine:3.16
3.16: Pulling from library/alpine
213ec9aee27d: Pull complete
Digest: sha256:bc41182d7ef5ffc53a40b044e725193bc10142a1243f395ee852a8d9730fc2ad
Status: Downloaded newer image for alpine:3.16
docker.io/library/alpine:3.16
/app/test/client/stevedore # docker images
REPOSITORY   TAG       IMAGE ID       CREATED       SIZE
alpine       3.16      9c6f07244728   6 weeks ago   5.54MB
```

- Tag `alpine:3.16` image to `docker-hub.stevedore.test:5000/library/alpine:3.16`

```sh
/app/test/client/stevedore # docker tag alpine:3.16 docker-hub.stevedore.test:5000/library/alpine:3.16
/app/test/client/stevedore # docker images
REPOSITORY                                      TAG       IMAGE ID       CREATED       SIZE
docker-hub.stevedore.test:5000/library/alpine   3.16      9c6f07244728   6 weeks ago   5.54MB
alpine                                          3.16      9c6f07244728   6 weeks ago   5.54MB
/app/test/client/stevedore #
```

- Push the image `docker-hub.stevedore.test:5000/library/alpine:3.16` 

```sh
/app/test/client/stevedore # docker push docker-hub.stevedore.test:5000/library/alpine:3.16
The push refers to repository [docker-hub.stevedore.test:5000/library/alpine]
994393dc58e7: Pushed
3.16: digest: sha256:1304f174557314a7ed9eddb4eab12fed12cb0cd9809e4c28f29af86979a3c870 size: 528
/app/test/client/stevedore #
```

### registry.stevedore.test

`registry.stevedore.test` is the docker registry used during functional test to push and validate the images.

## Git

### Users

Git could be accessed by HTTP or SSH.
To access by SSH, use the ssh key pair located at `files/ssl`. To private key is protected by a password, which is `password`.

HTTP access requires authorization to be used. User: `admin` Password: `admin`

### How to create a bare repository

- Start testing stack

```sh
make start
```

- Attach to git container

```sh
docker compose exec gitserver sh
```

```sh
/git # git config --global init.defaultBranch main
/git # cd repos
/git/repos # git init --bare /git/repos/app2.git
Initialized empty Git repository in /git/repos/app2.git/
/git/repos # 
```

- Set permission properly to the new bare repository

```sh
/git/repos # chown -R git:git app2.git/
```

### How to clone a gitserver repositories

You must be authorized before clone or update any repository on gitserver. The `client` container is provided by a key pair, on `root` user, that are allowed to manage git repositories.

```sh
/app/test/client/stevedore # ls -l ~/.ssh/
total 12
-rw-------    1 root     root          2655 Sep 26 11:11 id_rsa
-rw-------    1 root     root           577 Sep 26 11:11 id_rsa.pub
-rw-------    1 root     root           978 Sep 26 11:12 known_hosts
```

- Install git package on `client` because it is not installed by default

```sh
/app/test/client/stevedore # apk --update add git
fetch https://dl-cdn.alpinelinux.org/alpine/v3.16/main/x86_64/APKINDEX.tar.gz
fetch https://dl-cdn.alpinelinux.org/alpine/v3.16/community/x86_64/APKINDEX.tar.gz
(1/6) Installing brotli-libs (1.0.9-r6)
(2/6) Installing nghttp2-libs (1.47.0-r0)
(3/6) Installing libcurl (7.83.1-r3)
(4/6) Installing expat (2.4.9-r0)
(5/6) Installing pcre2 (10.40-r0)
(6/6) Installing git (2.36.2-r0)
Executing busybox-1.35.0-r17.trigger
OK: 43 MiB in 61 packages
```

```sh
/app/test/client/stevedore # cd /tmp
/tmp # git clone git@gitserver:/git/repos/app2.git
Cloning into 'go-docker-builder-alpine'...
Enter passphrase for key '/root/.ssh/id_rsa':
remote: Enumerating objects: 14, done.
remote: Counting objects: 100% (14/14), done.
remote: Compressing objects: 100% (11/11), done.
remote: Total 14 (delta 2), reused 0 (delta 0), pack-reused 0
Receiving objects: 100% (14/14), done.
Resolving deltas: 100% (2/2), done.
```
