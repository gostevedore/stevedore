# Inject Dockerfile To The Build Context Example

This example demonstrates how to create a unified build context using different sources. The source code for the applications is stored in a Git repository, but these repositories do not include a Dockerfile. With Stevedore's ability to merge multiple build context sources, the Dockerfile is injected from a local folder.

- [Inject Dockerfile To The Build Context Example](#inject-dockerfile-to-the-build-context-example)
  - [Requirements](#requirements)
  - [Stack](#stack)
  - [Usage](#usage)
  - [Example Execution Insights](#example-execution-insights)
    - [Starting the Stack](#starting-the-stack)
    - [Waiting for Dockerd to be Ready](#waiting-for-dockerd-to-be-ready)
    - [Getting Credentials](#getting-credentials)
    - [Getting Images](#getting-images)
    - [Building images](#building-images)
    - [Cleaning the stack](#cleaning-the-stack)
  - [Additional Information](#additional-information)
    - [Images](#images)


## Requirements

- Docker. _Tested on Docker server 20.10.21 and Docker API 1.41_
- Docker's Compose plugin or `docker-compose`. _Tested on Docker Compose version v2.17.3_
- `make` utility. _Tested on version 4.3-4.1build1_

## Stack

The stack required to run this example is defined in a [Docker Compose file](./docker-compose.yml). The stack consists of four services: a Docker Registry, a Docker Registry authorization, a Git server and a Stevedore service.

The Docker registry is used to store the Docker images built by Stevedore during the example execution.
The Git server service is used to store the source code of the applications that we pretend to build during this example.
The Stevedore service is built from a container which is defined in that [Dockerfile](stack/stevedore/Dockerfile).

## Usage

The example comes with a Makefile that can help you execute common actions, like starting the stack to run the example or attaching to a container in the stack to perform specific tasks.

Find below the available Makefile targets, as well as its description:

```sh
❯ make help
 Example basic-example:
  help                      Lists allowed targets
  run                       Runs start, example, and clean targets together
  start                     Starts the stack required to run the example
  clean                     Stops the stack required to run the example
  status                    Displays the status of the stack
  follow-logs               Shows the stack logs in follow mode
  attach                    Attaches to the Stevedore container
  example                   Executes the example (requires the stack to be started)
```

To execute the entire example, including starting and cleaning the stack, run the `run` target.

```sh
❯ make run
```

## Example Execution Insights

Below is the expected output for the `make run` command, which starts the Docker stack, gets some information about the Stevedore configuration, builds and promotes Docker images using Stevedore, and then cleans the stack up.

### Starting the Stack

When the run starts, it generates an SSH key pair.

```sh
❯ make run
 [07-inject-dockerfile-example] Create SSH keys
[+] Building 0.0s (0/0)
[+] Creating 2/1
 ✔ Network 07-inject-dockerfile-example_default  Created                                                                                                             0.1s
 ✔ Volume "07-inject-dockerfile-example_ssh"     Created                                                                                                             0.0s
[+] Building 0.0s (0/0)
```

Once the SSH key pair is generated, the services are started.

```sh
Starting the stack to run 07-inject-dockerfile-example

2023/05/29 08:05:50 http2: server connection error from localhost: connection error: PROTOCOL_ERROR
[+] Building 8.0s (37/37) FINISHED
 => [stevedore internal] load .dockerignore                                                                                                                          0.0s
 => => transferring context: 2B                                                                                                                                      0.0s
 => [stevedore internal] load build definition from Dockerfile                                                                                                       0.0s
 => => transferring dockerfile: 912B                                                                                                                                 0.0s
 => [gitserver internal] load .dockerignore                                                                                                                          0.0s
 => => transferring context: 2B                                                                                                                                      0.0s
 => [gitserver internal] load build definition from Dockerfile                                                                                                       0.0s
 => => transferring dockerfile: 1.41kB                                                                                                                               0.0s
 => [stevedore internal] load metadata for docker.io/library/docker:20.10-dind                                                                                       1.3s
 => [stevedore internal] load metadata for docker.io/library/golang:1.19-alpine                                                                                      1.8s
 => [gitserver internal] load metadata for docker.io/library/alpine:3.16                                                                                             1.7s
 => [ssh-keygen internal] load build definition from Dockerfile                                                                                                      0.0s
 => => transferring dockerfile: 129B                                                                                                                                 0.0s
 => [ssh-keygen internal] load .dockerignore                                                                                                                         0.0s
 => => transferring context: 2B                                                                                                                                      0.0s
 => [ssh-keygen internal] load metadata for docker.io/library/ubuntu:22.04                                                                                           0.4s
 => [ssh-keygen 1/2] FROM docker.io/library/ubuntu:22.04@sha256:dfd64a3b4296d8c9b62aa3309984f8620b98d87e47492599ee20739e8eb54fbf                                     0.0s
 => => resolve docker.io/library/ubuntu:22.04@sha256:dfd64a3b4296d8c9b62aa3309984f8620b98d87e47492599ee20739e8eb54fbf                                                0.0s
 => CACHED [ssh-keygen 2/2] RUN apt-get update     && apt-get install -y         openssh-client                                                                      0.0s
 => [ssh-keygen] exporting to docker image format                                                                                                                    0.3s
 => => exporting layers                                                                                                                                              0.0s
 => => exporting manifest sha256:aa203c7f59e7293f54c4b1ffe4e0c7ffe3f075a30f5a44f63fc8b90b0815bc69                                                                    0.0s
 => => exporting config sha256:f39ad726717f0db4b8b04eb3793e86760fe39635dd343a2733a9331dda584e50                                                                      0.0s
 => => sending tarball                                                                                                                                               0.3s
 => [gitserver 1/6] FROM docker.io/library/alpine:3.16@sha256:c2b622f6e510a0d25bccaffa9e67b75a6860cb09b74bb58cfc36a9ef4331109f                                       0.0s
 => => resolve docker.io/library/alpine:3.16@sha256:c2b622f6e510a0d25bccaffa9e67b75a6860cb09b74bb58cfc36a9ef4331109f                                                 0.0s
 => [gitserver internal] load build context                                                                                                                          0.0s
 => => transferring context: 1.41kB                                                                                                                                  0.0s
 => CACHED [gitserver 2/6] WORKDIR /git                                                                                                                              0.0s
 => CACHED [gitserver 3/6] RUN apk add --no-cache         openssh         git     && rm -rf /var/cache/apk/*     && ssh-keygen -A     && adduser -D --home /home/gi  0.0s
 => CACHED [gitserver 4/6] RUN apk add  --no-cache         nginx         git-daemon         fcgiwrap         spawn-fcgi     && rm -rf /var/cache/apk/*               0.0s
 => CACHED [gitserver 5/6] COPY files/nginx.conf /etc/nginx/nginx.conf                                                                                               0.0s
 => CACHED [gitserver 6/6] COPY entrypoint.sh /entrypoint.sh                                                                                                         0.0s
 => [gitserver] exporting to docker image format                                                                                                                     0.1s
 => => exporting layers                                                                                                                                              0.0s
 => => exporting manifest sha256:d1105f864f7ed842063697121fa40c2ab8092206ebb87c1db9fbf6e8ef0d68f5                                                                    0.0s
 => => exporting config sha256:86f36604c0a7e35afae5b84ddb7ec5ba8765d5967cdc0e433a6fe3b698f62699                                                                      0.0s
 => => sending tarball                                                                                                                                               0.1s
 => [ssh-keygen ssh-keygen] importing to docker                                                                                                                      0.0s
 => [gitserver gitserver] importing to docker                                                                                                                        0.0s
 => [stevedore stage-1 1/4] FROM docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                         0.0s
 => => resolve docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                                           0.0s
 => [stevedore golang 1/7] FROM docker.io/library/golang:1.19-alpine@sha256:4147c2ad24a347e53a8a361b663dcf82fca4157a3f9a5136d696a4e53dd22b65                         0.0s
 => => resolve docker.io/library/golang:1.19-alpine@sha256:4147c2ad24a347e53a8a361b663dcf82fca4157a3f9a5136d696a4e53dd22b65                                          0.0s
 => [stevedore internal] load build context                                                                                                                          0.1s
 => => transferring context: 397.45kB                                                                                                                                0.0s
 => CACHED [stevedore golang 2/7] WORKDIR /usr/src/app                                                                                                               0.0s
 => CACHED [stevedore golang 3/7] RUN apk add --no-cache make build-base                                                                                             0.0s
 => CACHED [stevedore golang 4/7] COPY go.mod go.sum ./                                                                                                              0.0s
 => CACHED [stevedore golang 5/7] RUN go mod download                                                                                                                0.0s
 => [stevedore golang 6/7] COPY . ./                                                                                                                                 0.4s
 => [stevedore golang 7/7] RUN go build -ldflags "-s -w" -v -o /usr/local/bin/stevedore ./cmd/stevedore.go                                                           5.1s
 => CACHED [stevedore stage-1 2/4] COPY --from=golang /usr/local/bin/stevedore /usr/local/bin/stevedore                                                              0.0s
 => CACHED [stevedore stage-1 3/4] COPY examples/07-inject-dockerfile-example/stack/stevedore/entrypoint.sh /usr/local/bin/entrypoint.sh                             0.0s
 => CACHED [stevedore stage-1 4/4] COPY examples/07-inject-dockerfile-example/stack/stevedore/wait-for-dockerd.sh /usr/local/bin/wait-for-dockerd.sh                 0.0s
 => [stevedore] exporting to docker image format                                                                                                                     0.5s
 => => exporting layers                                                                                                                                              0.0s
 => => exporting manifest sha256:2b2f6663bb8acf6bf37f44c26875d71babec1aac49a80e3e67442f556e798ea3                                                                    0.0s
 => => exporting config sha256:f3a70d2392a91a8ef5b6c075fa276b4e33c875c974f6546e4972e9a19fd207e0                                                                      0.0s
 => => sending tarball                                                                                                                                               0.5s
 => [stevedore stevedore] importing to docker                                                                                                                        0.0s
[+] Running 6/6
 ✔ Container 07-inject-dockerfile-example-ssh-keygen-1  Started                                                                                                      0.6s
 ✔ Container 07-inject-dockerfile-example-stevedore-1   Started                                                                                                      1.2s
 ✔ Container 07-inject-dockerfile-example-dockerauth-1  Started                                                                                                      0.6s
 ✔ Container 07-inject-dockerfile-example-gitserver-1   Started                                                                                                      1.2s
 ✔ Container 07-inject-dockerfile-example-worker-1      Started                                                                                                      0.7s
 ✔ Container 07-inject-dockerfile-example-registry-1    Started                                                                                                      1.2s
```

And finally, it generates a `known_hosts` file based on the SSH keys.

```sh
 [07-inject-dockerfile-example] Create known_hosts file
# gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
# gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
# gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
# gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
# gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
```

### Waiting for Dockerd to be Ready

Before starting the execution of the Stevedore command, it is important to ensure that the Docker daemon (dockerd) is ready. The stevedore service Docker image includes a script, [wait-for-dockerd.sh]((./stack/stevedore/wait-for-dockerd.sh)), which can be used to ensure the readiness of the Docker daemon.

```sh
 Run example 07-inject-dockerfile-example

 [07-inject-dockerfile-example] Waiting for dockerd
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
```

### Getting Credentials

To obtain the credentials, use the command: `stevedore get credentials`.

```sh
 [07-inject-dockerfile-example] Get credentials
ID                             TYPE              CREDENTIALS
registry.stevedore.test        username-password username=admin
https_gitserver.stevedore.test username-password username=admin
ssh_gitserver.stevedore.test   Private key file  private_key_file=/root/.ssh/id_rsa, protected by password
```

### Getting Images

To view the images in tree format, run `stevedore get images --tree`.

```sh
 [07-inject-dockerfile-example] Get images
├─── busybox:1.36
│  ├─── registry.stevedore.test/base:busybox-1.36
│  │  ├─── registry.stevedore.test/app2:v1-base-busybox-1.36
│  │  ├─── registry.stevedore.test/app3:v1-base-busybox-1.36
├─── busybox:1.35
│  ├─── registry.stevedore.test/base:busybox-1.35
│  │  ├─── registry.stevedore.test/app2:v1-base-busybox-1.35
│  │  ├─── registry.stevedore.test/app1:v1-base-busybox-1.35
```

### Building images

In this section, you can explore the output of the Docker image builds for `base` images and the applications `app1`, `app2` and `app3`. The image definitions for these applications reside in the [./image](./images) folder, the applications' source code is achieved from the Git server, however, the Dockerfile injected is defined in the [./images-src](./images-src/build/Dockerfile) folder.
Note that each Docker image is built independently of one another.

The example utilizes the command `stevedore build base --build-on-cascade --push-after-build` to build the images.

```sh
 [07-inject-dockerfile-example] Build the base image and its descendants, and push the images after build
registry.stevedore.test/base:busybox-1.35 Step 1/7 : ARG image_from_name
registry.stevedore.test/base:busybox-1.35 Step 2/7 : ARG image_from_tag
registry.stevedore.test/base:busybox-1.35 Step 3/7 : FROM ${image_from_name}:${image_from_tag}
registry.stevedore.test/base:busybox-1.36 Step 1/7 : ARG image_from_name
registry.stevedore.test/base:busybox-1.36 Step 2/7 : ARG image_from_tag
registry.stevedore.test/base:busybox-1.36 Step 3/7 : FROM ${image_from_name}:${image_from_tag}
registry.stevedore.test/base:busybox-1.35 ‣  1.35:  Pulling from library/busybox
registry.stevedore.test/base:busybox-1.35 ‣  c15cbdab5f8e:  Pull complete
registry.stevedore.test/base:busybox-1.35 ‣  Digest: sha256:b4e4a06de46acc0958cd93e2eeb769077d255f06a7c3a91196509c16b7bc989e
registry.stevedore.test/base:busybox-1.35 ‣  Status: Downloaded newer image for busybox:1.35
registry.stevedore.test/base:busybox-1.35 ---> dddc7578369a
registry.stevedore.test/base:busybox-1.35 Step 4/7 : RUN echo "anonymous:x:10001:10001:,,,:/app:/bin/sh" >> /etc/passwd &&     echo "anonymous:x:10001:" >> /etc/group &&     mkdir -p /app &&     chown 10001:10001 /app
registry.stevedore.test/base:busybox-1.35 ---> Running in 4e4307a4b9b7
registry.stevedore.test/base:busybox-1.35 ---> 4dd449afec45
registry.stevedore.test/base:busybox-1.35 Step 5/7 : USER anonymous
registry.stevedore.test/base:busybox-1.35 ---> Running in babb9f2d6272
registry.stevedore.test/base:busybox-1.35 ---> 2fc2d6f4d817
registry.stevedore.test/base:busybox-1.35 Step 6/7 : WORKDIR /app
registry.stevedore.test/base:busybox-1.35 ---> Running in 63409e3d42b0
registry.stevedore.test/base:busybox-1.35 ---> 079d4375f015
registry.stevedore.test/base:busybox-1.35 Step 7/7 : LABEL created_at=2023-05-29T06:06:02.188930103Z
registry.stevedore.test/base:busybox-1.35 ---> Running in 6f8212bf46e2
registry.stevedore.test/base:busybox-1.35 ---> 06e19bde9972
registry.stevedore.test/base:busybox-1.35  ‣ sha256:06e19bde99726c5f6cf7f85fe15586d810f2ca32d435e3ef08528e99304ccd55
registry.stevedore.test/base:busybox-1.35 Successfully built 06e19bde9972
registry.stevedore.test/base:busybox-1.35 Successfully tagged registry.stevedore.test/base:busybox-1.35
registry.stevedore.test/base:busybox-1.35 ‣  The push refers to repository [registry.stevedore.test/base]
registry.stevedore.test/base:busybox-1.36 ‣  1.36:  Pulling from library/busybox
registry.stevedore.test/base:busybox-1.36 ‣  325d69979d33:  Pull complete
registry.stevedore.test/base:busybox-1.36 ‣  Digest: sha256:560af6915bfc8d7630e50e212e08242d37b63bd5c1ccf9bd4acccf116e262d5b
registry.stevedore.test/base:busybox-1.36 ‣  Status: Downloaded newer image for busybox:1.36
registry.stevedore.test/base:busybox-1.36 ---> 8135583d97fe
registry.stevedore.test/base:busybox-1.36 Step 4/7 : RUN echo "anonymous:x:10001:10001:,,,:/app:/bin/sh" >> /etc/passwd &&     echo "anonymous:x:10001:" >> /etc/group &&     mkdir -p /app &&     chown 10001:10001 /app
registry.stevedore.test/base:busybox-1.36 ---> Running in a4718d8eb6ed
registry.stevedore.test/base:busybox-1.35 ‣  busybox-1.35: digest: sha256:ff134176adf74b7926e41eb8a0c877f3614e6d4433e8e7ce69e3a1635bc1a278 size: 735
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 1/8 : ARG image_from_name
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 2/8 : ARG image_from_tag
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 3/8 : ARG image_from_registry_host
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 4/8 : FROM ${image_from_registry_host}/${image_from_name}:${image_from_tag}
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> 06e19bde9972
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 5/8 : ARG app_name
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> Running in 9655da91a8ce
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> eec81b66b74e
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 6/8 : COPY ./${app_name}/app.sh /app.sh
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> 673e2b525541
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 7/8 : CMD ["/app.sh"]
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> Running in 082e1c9fb00a
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> 3858714cd258
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 8/8 : LABEL created_at=2023-05-29T06:06:02.188930103Z
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> Running in 73e498ee998d
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> 424054555793
registry.stevedore.test/app1:v1-base-busybox-1.35  ‣ sha256:424054555793562dfc8b86d98d5954ead24f068588dbea962da318ee98d497d9
registry.stevedore.test/app1:v1-base-busybox-1.35 Successfully built 424054555793
registry.stevedore.test/app1:v1-base-busybox-1.35 Successfully tagged registry.stevedore.test/app1:v1-base-busybox-1.35
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  The push refers to repository [registry.stevedore.test/app1]
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  f8f99cbf0a5a:  Preparing
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  4e8c7e0e829f:  Preparing
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  42ef21f45b9a:  Preparing
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 1/8 : ARG image_from_name
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 2/8 : ARG image_from_tag
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 3/8 : ARG image_from_registry_host
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 4/8 : FROM ${image_from_registry_host}/${image_from_name}:${image_from_tag}
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> 06e19bde9972
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 5/8 : ARG app_name
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> Using cache
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  f8f99cbf0a5a:  Pushing [==================================================>]     512B
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  4e8c7e0e829f:  Preparing
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  42ef21f45b9a:  Preparing
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> 906df6a59a2e
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 7/8 : CMD ["/app.sh"]
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> Running in 87aa6f651b84
registry.stevedore.test/base:busybox-1.36 ---> a7ac154c33f4
registry.stevedore.test/base:busybox-1.36 Step 5/7 : USER anonymous
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> 99df79468934
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 8/8 : LABEL created_at=2023-05-29T06:06:02.188930103Z
registry.stevedore.test/base:busybox-1.36 ---> Running in fbc50db5d7ba
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> Running in 47f8570813fa
registry.stevedore.test/base:busybox-1.36 ---> c377ac5e061b
registry.stevedore.test/base:busybox-1.36 Step 6/7 : WORKDIR /app
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> 22eb4609caa2
registry.stevedore.test/app2:v1-base-busybox-1.35  ‣ sha256:22eb4609caa2022e325e2e1514e903d436ec5f7020a9ac2c107f0cfb7986d758
registry.stevedore.test/app2:v1-base-busybox-1.35 Successfully built 22eb4609caa2
registry.stevedore.test/app2:v1-base-busybox-1.35 Successfully tagged registry.stevedore.test/app2:v1-base-busybox-1.35
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  The push refers to repository [registry.stevedore.test/app2]
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  99516c98f205:  Preparing
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  4e8c7e0e829f:  Preparing
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  42ef21f45b9a:  Waiting
registry.stevedore.test/base:busybox-1.36 ---> Running in 269ef244f44e
registry.stevedore.test/base:busybox-1.36 ---> 03bb6486d509
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  f8f99cbf0a5a:  Pushing [==================================================>]     512B
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  4e8c7e0e829f:  Mounted from base
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  42ef21f45b9a:  Mounted from base
registry.stevedore.test/base:busybox-1.36 ---> ca0cb73b3a77
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  f8f99cbf0a5a:  Pushed
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  4e8c7e0e829f:  Mounted from base
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  42ef21f45b9a:  Mounted from base
registry.stevedore.test/base:busybox-1.36 Successfully tagged registry.stevedore.test/base:busybox-1.36
registry.stevedore.test/base:busybox-1.36 ‣  The push refers to repository [registry.stevedore.test/base]
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  99516c98f205:  Pushed
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  4e8c7e0e829f:  Mounted from base
registry.stevedore.test/base:busybox-1.36 ‣  475721f5c964:  Pushed
registry.stevedore.test/base:busybox-1.36 ‣  9547b4c33213:  Pushed
registry.stevedore.test/base:busybox-1.36 ‣  busybox-1.36: digest: sha256:76bec13deb4b9c40cd618bd7da702a6f246b6114f64cf1a8201a4d601b81d4f2 size: 735
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 1/8 : ARG image_from_name
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 2/8 : ARG image_from_tag
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 3/8 : ARG image_from_registry_host
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 4/8 : FROM ${image_from_registry_host}/${image_from_name}:${image_from_tag}
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> ca0cb73b3a77
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 5/8 : ARG app_name
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> Running in 097aa951ee3f
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> 83937e1e5499
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 6/8 : COPY ./${app_name}/app.sh /app.sh
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> 5bd61b8aae40
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 7/8 : CMD ["/app.sh"]
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> Running in 1bd141fb9e3f
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> 5ddff94d0efc
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 8/8 : LABEL created_at=2023-05-29T06:06:02.188418492Z
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> Running in 8d672b0bb79c
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> cc9e86726a8a
registry.stevedore.test/app3:v1-base-busybox-1.36  ‣ sha256:cc9e86726a8ad440392f07bb5d96e929468cb5b0943d7253e8b2e9d840d3a2c2
registry.stevedore.test/app3:v1-base-busybox-1.36 Successfully built cc9e86726a8a
registry.stevedore.test/app3:v1-base-busybox-1.36 Successfully tagged registry.stevedore.test/app3:v1-base-busybox-1.36
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  The push refers to repository [registry.stevedore.test/app3]
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  77bc55a2b896:  Preparing
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  475721f5c964:  Preparing
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  9547b4c33213:  Preparing
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 1/8 : ARG image_from_name
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 2/8 : ARG image_from_tag
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 3/8 : ARG image_from_registry_host
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 4/8 : FROM ${image_from_registry_host}/${image_from_name}:${image_from_tag}
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> ca0cb73b3a77
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 5/8 : ARG app_name
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> Using cache
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> 83937e1e5499
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 6/8 : COPY ./${app_name}/app.sh /app.sh
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> 26df540e6ece
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  77bc55a2b896:  Pushing [==================================================>]     512B
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  475721f5c964:  Preparing
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  9547b4c33213:  Preparing
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> 1e3c83b3ab1a
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 8/8 : LABEL created_at=2023-05-29T06:06:02.188418492Z
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> Running in 51ae144ab52b
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> d991337255de
registry.stevedore.test/app2:v1-base-busybox-1.36  ‣ sha256:d991337255de6b378e57668169a9aa404961c88cefc82bfa0898cc78563cd3e2
registry.stevedore.test/app2:v1-base-busybox-1.36 Successfully built d991337255de
registry.stevedore.test/app2:v1-base-busybox-1.36 Successfully tagged registry.stevedore.test/app2:v1-base-busybox-1.36
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  The push refers to repository [registry.stevedore.test/app2]
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  77bc55a2b896:  Pushed
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  a8e187f3f916:  Pushed
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  475721f5c964:  Mounted from base
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  9547b4c33213:  Mounted from app3
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  v1-base-busybox-1.36: digest: sha256:3b7c04f9dede5788d0cb50c493a56d8b9dc3f96dbab1601e62cb3cd76ee4887d size: 942
```

### Cleaning the stack

```sh
Stopping the stack to run 07-inject-dockerfile-example

[+] Running 8/8
 ✔ Container 07-inject-dockerfile-example-registry-1    Removed                                                                                                      0.2s
 ✔ Container 07-inject-dockerfile-example-stevedore-1   Removed                                                                                                      3.3s
 ✔ Container 07-inject-dockerfile-example-ssh-keygen-1  Removed                                                                                                      3.5s
 ✔ Container 07-inject-dockerfile-example-worker-1      Removed                                                                                                      0.0s
 ✔ Container 07-inject-dockerfile-example-gitserver-1   Removed                                                                                                      3.6s
 ✔ Container 07-inject-dockerfile-example-dockerauth-1  Removed                                                                                                      0.3s
 ✔ Volume 07-inject-dockerfile-example_ssh              Removed                                                                                                      0.0s
 ✔ Network 07-inject-dockerfile-example_default         Removed                                                                                                      0.4s
```

## Additional Information

### Images

In this example, the image definition includes a builder with a list of build contexts. The first Docker build context is a [path context](https://gostevedore.github.io/docs/reference-guide/builder/docker/#path-context), which utilizes a local folder as the Docker build context. The Dockerfile used to build the applications' Docker images is located in this folder.

The second context is a [Git context](https://gostevedore.github.io/docs/reference-guide/builder/docker/#git-context), which utilizes a Git repository as the build context. This feature is available only when using the [docker driver](https://gostevedore.github.io/docs/reference-guide/driver/docker/).

Here you have an example of `app1` image definition.

```yaml
  app1:
    v1:
      version: "{{ .Version }}-{{ .Parent.Name }}-{{ .Parent.Version }}"
      registry: registry.stevedore.test
      builder:
        driver: docker
        options:
          context:
            - path: images-src/build
            - git:
                repository: http://gitserver.stevedore.test:/git/repos/app1.git
                reference: main
                auth:
                  credentials_id: https_gitserver.stevedore.test

      parents:
        base:
          - busybox-1.35
```
