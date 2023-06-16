# Git Build Context Example

This example illustrates how to utilize the [Git build context](https://gostevedore.github.io/docs/reference-guide/builder/docker/#git-context) feature in Stevedore. It demonstrates the ability to specify a Git repository as the build context, allowing you to directly build images from a specific branch, tag, or commit.

- [Git Build Context Example](#git-build-context-example)
  - [Requirements](#requirements)
  - [Stack](#stack)
  - [Usage](#usage)
  - [Example Execution Insights](#example-execution-insights)
    - [Starting the Stack](#starting-the-stack)
    - [Waiting for Dockerd to be Ready](#waiting-for-dockerd-to-be-ready)
    - [Getting images](#getting-images)
    - [Building images](#building-images)
    - [Cleaning the stack](#cleaning-the-stack)
  - [Additional Information](#additional-information)
    - [Credentials](#credentials)
    - [Images](#images)
    - [Git Server](#git-server)
      - [Preparing the Git Fixtures Required to Execute the Example](#preparing-the-git-fixtures-required-to-execute-the-example)
        - [Creating the Basic Auth Password](#creating-the-basic-auth-password)
        - [Creating the SSH Key Pair](#creating-the-ssh-key-pair)
        - [Creating the known\_hosts File](#creating-the-known_hosts-file)
        - [Preparing the Git Repository and the Application's Source Code](#preparing-the-git-repository-and-the-applications-source-code)
        - [Clone the application](#clone-the-application)


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
 [06-git-build-context-example] Create SSH keys
[+] Building 0.0s (0/0)
[+] Creating 2/0
 ✔ Network 06-git-build-context-example_default  Created                                                                                                             0.1s
 ✔ Volume "06-git-build-context-example_ssh"     Created                                                                                                             0.0s
[+] Building 0.0s (0/0)
```

Once the SSH key pair is generated, the services are started.
```sh
Starting the stack to run 06-git-build-context-example

[+] Building 6.4s (37/37) FINISHED
 => [stevedore internal] load build definition from Dockerfile                                                                                                       0.0s
 => => transferring dockerfile: 912B                                                                                                                                 0.0s
 => [stevedore internal] load .dockerignore                                                                                                                          0.0s
 => => transferring context: 2B                                                                                                                                      0.0s
 => [stevedore internal] load metadata for docker.io/library/docker:20.10-dind                                                                                       0.8s
 => [stevedore internal] load metadata for docker.io/library/golang:1.19-alpine                                                                                      0.8s
 => [ssh-keygen internal] load .dockerignore                                                                                                                         0.0s
 => => transferring context: 2B                                                                                                                                      0.0s
 => [ssh-keygen internal] load build definition from Dockerfile                                                                                                      0.0s
 => => transferring dockerfile: 129B                                                                                                                                 0.0s
 => [gitserver internal] load .dockerignore                                                                                                                          0.0s
 => => transferring context: 2B                                                                                                                                      0.0s
 => [gitserver internal] load build definition from Dockerfile                                                                                                       0.0s
 => => transferring dockerfile: 1.41kB                                                                                                                               0.0s
 => [ssh-keygen internal] load metadata for docker.io/library/ubuntu:22.04                                                                                           0.6s
 => [gitserver internal] load metadata for docker.io/library/alpine:3.16                                                                                             0.6s
 => [gitserver 1/6] FROM docker.io/library/alpine:3.16@sha256:c2b622f6e510a0d25bccaffa9e67b75a6860cb09b74bb58cfc36a9ef4331109f                                       0.0s
 => => resolve docker.io/library/alpine:3.16@sha256:c2b622f6e510a0d25bccaffa9e67b75a6860cb09b74bb58cfc36a9ef4331109f                                                 0.0s
 => [gitserver internal] load build context                                                                                                                          0.0s
 => => transferring context: 1.41kB                                                                                                                                  0.0s
 => [ssh-keygen 1/2] FROM docker.io/library/ubuntu:22.04@sha256:dfd64a3b4296d8c9b62aa3309984f8620b98d87e47492599ee20739e8eb54fbf                                     0.0s
 => => resolve docker.io/library/ubuntu:22.04@sha256:dfd64a3b4296d8c9b62aa3309984f8620b98d87e47492599ee20739e8eb54fbf                                                0.0s
 => CACHED [gitserver 2/6] WORKDIR /git                                                                                                                              0.0s
 => CACHED [gitserver 3/6] RUN apk add --no-cache         openssh         git     && rm -rf /var/cache/apk/*     && ssh-keygen -A     && adduser -D --home /home/gi  0.0s
 => CACHED [gitserver 4/6] RUN apk add  --no-cache         nginx         git-daemon         fcgiwrap         spawn-fcgi     && rm -rf /var/cache/apk/*               0.0s
 => CACHED [gitserver 5/6] COPY files/nginx.conf /etc/nginx/nginx.conf                                                                                               0.0s
 => CACHED [gitserver 6/6] COPY entrypoint.sh /entrypoint.sh                                                                                                         0.0s
 => [gitserver] exporting to docker image format                                                                                                                     0.1s
 => => exporting layers                                                                                                                                              0.0s
 => => exporting manifest sha256:ed529a96d72d0665a203d57b724f81cd6fcf6ff6a2a3b289e47edf0646deef72                                                                    0.0s
 => => exporting config sha256:35f4e457081993060a2c212564c2b87dd5baf601b57e5b780fb2c30786b04fb9                                                                      0.0s
 => => sending tarball                                                                                                                                               0.1s
 => CACHED [ssh-keygen 2/2] RUN apt-get update     && apt-get install -y         openssh-client                                                                      0.0s
 => [ssh-keygen] exporting to docker image format                                                                                                                    0.2s
 => => exporting layers                                                                                                                                              0.0s
 => => exporting manifest sha256:f4f6163e5f92e30135e2e9cf81d655098288082af34addc4789d958aa93de5cd                                                                    0.0s
 => => exporting config sha256:8eb337ed09d3afa6a01a3a35e1523542fbe30cb6a661ca994dd8ca84521e4a28                                                                      0.0s
 => => sending tarball                                                                                                                                               0.2s
 => [gitserver gitserver] importing to docker                                                                                                                        0.0s
 => [stevedore internal] load build context                                                                                                                          0.1s
 => => transferring context: 348.82kB                                                                                                                                0.1s
 => [stevedore golang 1/7] FROM docker.io/library/golang:1.19-alpine@sha256:4147c2ad24a347e53a8a361b663dcf82fca4157a3f9a5136d696a4e53dd22b65                         0.0s
 => => resolve docker.io/library/golang:1.19-alpine@sha256:4147c2ad24a347e53a8a361b663dcf82fca4157a3f9a5136d696a4e53dd22b65                                          0.0s
 => [stevedore stage-1 1/4] FROM docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                         0.0s
 => => resolve docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                                           0.0s
 => [ssh-keygen ssh-keygen] importing to docker                                                                                                                      0.0s
 => CACHED [stevedore golang 2/7] WORKDIR /usr/src/app                                                                                                               0.0s
 => CACHED [stevedore golang 3/7] RUN apk add --no-cache make build-base                                                                                             0.0s
 => CACHED [stevedore golang 4/7] COPY go.mod go.sum ./                                                                                                              0.0s
 => CACHED [stevedore golang 5/7] RUN go mod download                                                                                                                0.0s
 => [stevedore golang 6/7] COPY . ./                                                                                                                                 0.5s
 => [stevedore golang 7/7] RUN go build -ldflags "-s -w" -v -o /usr/local/bin/stevedore ./cmd/stevedore.go                                                           4.6s
 => CACHED [stevedore stage-1 2/4] COPY --from=golang /usr/local/bin/stevedore /usr/local/bin/stevedore                                                              0.0s
 => CACHED [stevedore stage-1 3/4] COPY examples/06-git-build-context-example/stack/stevedore/entrypoint.sh /usr/local/bin/entrypoint.sh                             0.0s
 => CACHED [stevedore stage-1 4/4] COPY examples/06-git-build-context-example/stack/stevedore/wait-for-dockerd.sh /usr/local/bin/wait-for-dockerd.sh                 0.0s
 => [stevedore] exporting to docker image format                                                                                                                     0.4s
 => => exporting layers                                                                                                                                              0.0s
 => => exporting manifest sha256:37501a85dd34a8740bd2b22b2adc35b055abdfdd8c65120c983011f5c12ff10d                                                                    0.0s
 => => exporting config sha256:835c3e7dc1597999596540ed0509ca34f20cda7765152ec4761d453dcf06511a                                                                      0.0s
 => => sending tarball                                                                                                                                               0.4s
 => [stevedore stevedore] importing to docker                                                                                                                        0.0s
[+] Running 6/6
 ✔ Container 06-git-build-context-example-ssh-keygen-1  Started                                                                                                      0.5s
 ✔ Container 06-git-build-context-example-gitserver-1   Started                                                                                                      1.0s
 ✔ Container 06-git-build-context-example-stevedore-1   Started                                                                                                      1.0s
 ✔ Container 06-git-build-context-example-dockerauth-1  Started                                                                                                      0.9s
 ✔ Container 06-git-build-context-example-registry-1    Started                                                                                                      1.1s
```

And finally, it generates a `known_hosts` file based on the SSH keys.

```sh
 [06-git-build-context-example] Create known_hosts file
# gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
# gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
# gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
# gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
# gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
```

### Waiting for Dockerd to be Ready

Before starting the execution of the Stevedore command, it is important to ensure that the Docker daemon (dockerd) is ready. The stevedore service Docker image includes a script, [wait-for-dockerd.sh]((./stack/stevedore/wait-for-dockerd.sh)), which can be used to ensure the readiness of the Docker daemon.

```sh
 Run example 06-git-build-context-example

 [06-git-build-context-example] Waiting for dockerd
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
```

### Getting images

To view the images in tree format, run `stevedore get images --tree`.

```sh
 [06-git-build-context-example] Get images
├─── busybox:1.35
│  ├─── registry.stevedore.test/base:busybox-1.35
│  │  ├─── registry.stevedore.test/app1:v1-base-busybox-1.35
│  │  ├─── registry.stevedore.test/app2:v1-base-busybox-1.35
├─── busybox:1.36
│  ├─── registry.stevedore.test/base:busybox-1.36
│  │  ├─── registry.stevedore.test/app2:v1-base-busybox-1.36
│  │  ├─── registry.stevedore.test/app3:v1-base-busybox-1.36
```

### Building images

In this section, you can explore the output of the Docker image builds for `base` images and the applications `app1`, `app2` and `app3`. The image definitions for these applications reside in the [./images](./images) folder, and the applications' source code is achieved from the Git server.
Note that each Docker image is built independently of one another.

```sh
 [06-git-build-context-example] Build the base image and its descendants, and push the images after build
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
registry.stevedore.test/base:busybox-1.36 ‣  1.36:  Pulling from library/busybox
registry.stevedore.test/base:busybox-1.36 ‣  325d69979d33:  Pull complete
registry.stevedore.test/base:busybox-1.36 ‣  Digest: sha256:560af6915bfc8d7630e50e212e08242d37b63bd5c1ccf9bd4acccf116e262d5b
registry.stevedore.test/base:busybox-1.36 ‣  Status: Downloaded newer image for busybox:1.36
registry.stevedore.test/base:busybox-1.36 ---> 8135583d97fe
registry.stevedore.test/base:busybox-1.36 Step 4/7 : RUN echo "anonymous:x:10001:10001:,,,:/app:/bin/sh" >> /etc/passwd &&     echo "anonymous:x:10001:" >> /etc/group &&     mkdir -p /app &&     chown 10001:10001 /app
registry.stevedore.test/base:busybox-1.36 ---> Running in 4897b128abb7
registry.stevedore.test/base:busybox-1.35 ---> e01093e4e79e
registry.stevedore.test/base:busybox-1.35 Step 5/7 : USER anonymous
registry.stevedore.test/base:busybox-1.35 ---> Running in 103f68bfd395
registry.stevedore.test/base:busybox-1.35 ---> f8b49d41d230
registry.stevedore.test/base:busybox-1.35 Step 6/7 : WORKDIR /app
registry.stevedore.test/base:busybox-1.35 ---> Running in da21646683a8
registry.stevedore.test/base:busybox-1.35 ---> 44afb4267261
registry.stevedore.test/base:busybox-1.35 Step 7/7 : LABEL created_at=2023-05-28T20:20:57.787559256Z
registry.stevedore.test/base:busybox-1.35 ---> Running in c42a29c67797
registry.stevedore.test/base:busybox-1.35 ---> ac09c9b5e8c5
registry.stevedore.test/base:busybox-1.35  ‣ sha256:ac09c9b5e8c54dca5df67db8f695494329f8ae8952ed6eec2fb23941c90628df
registry.stevedore.test/base:busybox-1.35 Successfully built ac09c9b5e8c5
registry.stevedore.test/base:busybox-1.35 Successfully tagged registry.stevedore.test/base:busybox-1.35
registry.stevedore.test/base:busybox-1.35 ‣  The push refers to repository [registry.stevedore.test/base]
registry.stevedore.test/base:busybox-1.35 ‣  58f17356d7ad:  Pushed
registry.stevedore.test/base:busybox-1.35 ‣  42ef21f45b9a:  Pushing [>                                                  ]  66.56kB/4.855MB
registry.stevedore.test/base:busybox-1.36 ---> e661c471d4e0
registry.stevedore.test/base:busybox-1.36 Step 5/7 : USER anonymous
registry.stevedore.test/base:busybox-1.36 ---> Running in 6287cb58960d
registry.stevedore.test/base:busybox-1.36 ---> 3083923d0d19
registry.stevedore.test/base:busybox-1.36 Step 6/7 : WORKDIR /app
registry.stevedore.test/base:busybox-1.35 ‣  58f17356d7ad:  Pushed
registry.stevedore.test/base:busybox-1.35 ‣  42ef21f45b9a:  Pushing [==================================>                ]  3.391MB/4.855MB
registry.stevedore.test/base:busybox-1.36 ---> 901a5a5b245d
registry.stevedore.test/base:busybox-1.36 Step 7/7 : LABEL created_at=2023-05-28T20:20:57.788632771Z
registry.stevedore.test/base:busybox-1.36 ---> Running in c85f73fef38d
registry.stevedore.test/base:busybox-1.36 ---> b5186750224a
registry.stevedore.test/base:busybox-1.36  ‣ sha256:b5186750224aab94d12bbdb532e39125824a08a84ffb25db9b9aaa9b99e446e7
registry.stevedore.test/base:busybox-1.36 Successfully built b5186750224a
registry.stevedore.test/base:busybox-1.36 Successfully tagged registry.stevedore.test/base:busybox-1.36
registry.stevedore.test/base:busybox-1.36 ‣  The push refers to repository [registry.stevedore.test/base]
registry.stevedore.test/base:busybox-1.35 ‣  58f17356d7ad:  Pushed
registry.stevedore.test/base:busybox-1.36 ‣  c7434056d8ca:  Pushed
registry.stevedore.test/base:busybox-1.36 ‣  9547b4c33213:  Pushing [>                                                  ]  66.56kB/4.863MB
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 1/7 : ARG image_from_name
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 2/7 : ARG image_from_tag
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 3/7 : ARG image_from_registry_host
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 4/7 : FROM ${image_from_registry_host}/${image_from_name}:${image_from_tag}
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> ac09c9b5e8c5
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 5/7 : COPY ./${app_name}/app.sh /app/run.sh
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> 5fa7f6e59200
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 6/7 : CMD ["/app/run.sh"]
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> Running in b27ae7c51fcb
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> 694226a77072
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 7/7 : LABEL created_at=2023-05-28T20:20:57.787559256Z
registry.stevedore.test/base:busybox-1.36 ‣  c7434056d8ca:  Pushed
registry.stevedore.test/base:busybox-1.36 ‣  9547b4c33213:  Pushing [============================>                      ]  2.809MB/4.863MB
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> ce22012aaa67
registry.stevedore.test/app1:v1-base-busybox-1.35  ‣ sha256:ce22012aaa678b15f0e028327a5a6676b3408125f1aab9ee7ca2c1bbd3f38677
registry.stevedore.test/app1:v1-base-busybox-1.35 Successfully built ce22012aaa67
registry.stevedore.test/app1:v1-base-busybox-1.35 Successfully tagged registry.stevedore.test/app1:v1-base-busybox-1.35
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  The push refers to repository [registry.stevedore.test/app1]
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  955d71cf9c97:  Pushing [==================================================>]     512B
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  58f17356d7ad:  Preparing
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  42ef21f45b9a:  Preparing
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 1/7 : ARG image_from_name
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 2/7 : ARG image_from_tag
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 3/7 : ARG image_from_registry_host
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 4/7 : FROM ${image_from_registry_host}/${image_from_name}:${image_from_tag}
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> ac09c9b5e8c5
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 5/7 : COPY ./${app_name}/app.sh /app/run.sh
registry.stevedore.test/base:busybox-1.36 ‣  busybox-1.36: digest: sha256:b483e41fdb4bc876e11cbf6dcd5d71e1fd6ef470584bb890cd77b7617a877bce size: 735
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 1/7 : ARG image_from_name
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 2/7 : ARG image_from_tag
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 3/7 : ARG image_from_registry_host
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 4/7 : FROM ${image_from_registry_host}/${image_from_name}:${image_from_tag}
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> b5186750224a
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 5/7 : COPY ./${app_name}/app.sh /app/run.sh
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> 8c9b65ceb2c5
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 6/7 : CMD ["/app/run.sh"]
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> Running in 68e3f34527e3
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> 81cf47e506b0
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 7/7 : LABEL created_at=2023-05-28T20:20:57.787559256Z
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> 0996329d3973
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 6/7 : CMD ["/app/run.sh"]
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> Running in 62ec6e7a4dbd
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> Running in e5530b398a8c
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> 728b00624e8c
registry.stevedore.test/app2:v1-base-busybox-1.35  ‣ sha256:728b00624e8c6da328b909765393cb930ea47a86249df4446590cbf94753417a
registry.stevedore.test/app2:v1-base-busybox-1.35 Successfully built 728b00624e8c
registry.stevedore.test/app2:v1-base-busybox-1.35 Successfully tagged registry.stevedore.test/app2:v1-base-busybox-1.35
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> dfcd7f30a04b
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 7/7 : LABEL created_at=2023-05-28T20:20:57.788632771Z
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  The push refers to repository [registry.stevedore.test/app2]
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  1aae2243f5ff:  Preparing
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  58f17356d7ad:  Preparing
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  42ef21f45b9a:  Waiting
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> Running in 36b055b4f793
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> 8f22004e05a6
registry.stevedore.test/app3:v1-base-busybox-1.36  ‣ sha256:8f22004e05a68fd7e064392697403cd1c172c5d03f074066df96b1dc12ea3935
registry.stevedore.test/app3:v1-base-busybox-1.36 Successfully built 8f22004e05a6
registry.stevedore.test/app3:v1-base-busybox-1.36 Successfully tagged registry.stevedore.test/app3:v1-base-busybox-1.36
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  The push refers to repository [registry.stevedore.test/app3]
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  1aae2243f5ff:  Pushing [==================================================>]     512B
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  58f17356d7ad:  Preparing
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  42ef21f45b9a:  Waiting
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 1/7 : ARG image_from_name
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 2/7 : ARG image_from_tag
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 3/7 : ARG image_from_registry_host
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 4/7 : FROM ${image_from_registry_host}/${image_from_name}:${image_from_tag}
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> b5186750224a
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 5/7 : COPY ./${app_name}/app.sh /app/run.sh
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  v1-base-busybox-1.35: digest: sha256:f8a3d658d997f5799511613cea26785403672afb39b344373c2dd970ac954414 size: 942
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> 4b509844da96
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  c712f1b34602:  Pushing [==================================================>]     512B
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  c7434056d8ca:  Waiting
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  9547b4c33213:  Waiting
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> 250c1451f40d
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 7/7 : LABEL created_at=2023-05-28T20:20:57.788632771Z
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> Running in f9f08b40e436
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> c7750130bc1c
registry.stevedore.test/app2:v1-base-busybox-1.36  ‣ sha256:c7750130bc1c8bad3e22c4fb58874ebe7339deac33d02e414368dfcb1d01d0b7
registry.stevedore.test/app2:v1-base-busybox-1.36 Successfully built c7750130bc1c
registry.stevedore.test/app2:v1-base-busybox-1.36 Successfully tagged registry.stevedore.test/app2:v1-base-busybox-1.36
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  The push refers to repository [registry.stevedore.test/app2]
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  1aae2243f5ff:  Pushed
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  c712f1b34602:  Pushed
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  b4db87d02fe9:  Pushed
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  c7434056d8ca:  Mounted from app3
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  9547b4c33213:  Mounted from base
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  v1-base-busybox-1.36: digest: sha256:868e355fe911ae284c71f8a7736ed532172458995c354e224074c043a0e75ea8 size: 942
```

### Cleaning the stack

```sh
Stopping the stack to run 06-git-build-context-example

[+] Running 7/7
 ✔ Container 06-git-build-context-example-gitserver-1   Removed                                                                                                      3.6s
 ✔ Container 06-git-build-context-example-stevedore-1   Removed                                                                                                      3.6s
 ✔ Container 06-git-build-context-example-ssh-keygen-1  Removed                                                                                                      3.4s
 ✔ Container 06-git-build-context-example-registry-1    Removed                                                                                                      0.2s
 ✔ Container 06-git-build-context-example-dockerauth-1  Removed                                                                                                      0.3s
 ✔ Volume 06-git-build-context-example_ssh              Removed                                                                                                      0.0s
 ✔ Network 06-git-build-context-example_default         Removed                                                                                                      0.5s
```

## Additional Information

### Credentials

In this example, environment variables are used as the [credentials store](https://gostevedore.github.io/docs/reference-guide/credentials/credentials-store/#envvars-storage), specifically the `envvars` storage type in Stevedore.

In addition to the credentials required to access the Docker registry, this example introduces two new credentials for accessing the Git server. The `ssh_gitserver.stevedore.test` credentials contain a reference to a private key that is protected by a password. These credentials enable you to clone Git repositories using SSH keys for authentication. On the other hand, the `https_gitserver.stevedore.test` credentials represent a username-password pair that allows you to authenticate with the Git server using the basic auth method.

```sh
/app/examples/06-git-build-context-example # stevedore get credentials
ID                             TYPE              CREDENTIALS
registry.stevedore.test        username-password username=admin
ssh_gitserver.stevedore.test   Private key file  private_key_file=/root/.ssh/id_rsa, protected by password
https_gitserver.stevedore.test username-password username=admin
```

Both credentials has been created using the Stevedore subcommand [create credentials](https://gostevedore.github.io/docs/reference-guide/cli/#create-credentials).

### Images

This example uses the `base` Docker images, which defines a shared configuration, as parent images to build the images for the applications `app1`, `app2` and `app3`. What outstands that example is how the Docker images for those applications are built using a Git repository as a Docker build context. This approach allows you to build Docker images from remote resources.

For the `app1` application, the Git repository is accessed using the basic auth through HTTP, and the `https_gitserver.stevedore.test` credentials are used for this case.

```yaml
app1:
  v1:
    version: "{{ .Version }}-{{ .Parent.Name }}-{{ .Parent.Version }}"
    registry: registry.stevedore.test
    builder:
      driver: docker
      options:
        context:
          - git:
              repository: http://gitserver.stevedore.test:/git/repos/app1.git
              reference: main
              auth:
                credentials_id: https_gitserver.stevedore.test
    parents:
      base:
        - busybox-1.35
```

The `app2` utilizes SSH to access to the Git server repositories. The repository is referenced as `git@gitserver.stevedore.test:/git/repos/app2.git`, and the `ssh_gitserver.stevedore.test` credentials contain a reference to an SSH key that allows you to read the application repository.

```yaml
  app2:
    v1:
      version: "{{ .Version }}-{{ .Parent.Name }}-{{ .Parent.Version }}"
      registry: registry.stevedore.test
      builder:
        driver: docker
        options:
          context:
            - git:
                repository: git@gitserver.stevedore.test:/git/repos/app2.git
                reference: main
                auth:
                  credentials_id: ssh_gitserver.stevedore.test
      parents:
        base:
          - busybox-1.35
          - busybox-1.36
```

These different approaches showcase the flexibility of using Git repositories as Docker build contexts, allowing for seamless integration of version-controlled source code into the Docker image building process.

### Git Server

The example utilizes a Git server with three repositories, each containing the source code of an individual application, `app1`, `app2` and `app3`. These repositories can be found in the [fixtures](./fixtures/gitserver/repos) folder.
Access to the Git server is available via the host `gitserver.stevedeore.test`. You can authenticate using either the basic auth method with the credentials user=`admin` and password=`admin` or by utilizing SSH keys. The necessary keys are created automatically executing the `make start` command. It is important to note that the private key is password-protected, with the password set to `password`.

#### Preparing the Git Fixtures Required to Execute the Example

The Git server uses a set of fixed configurations prepared in advance. The [fixtures](./fixtures/) folder contains the pre-configured elements such as the authorization mechanisms and the repositories.

##### Creating the Basic Auth Password

You can authenticate to the Git server using the basic auth method. 
The `htpasswd` utility has been used to create the password used in that example. This command creates a password file named `.gitpasswd` with the username `admin` and its corresponding password, which in this case is set as `admin`.

```sh
root@6a88123e81cc:/$ htpasswd -c .gitpasswd admin
New password:
Re-type new password:
Adding password for user admin
```

You can find the resulting `.gitpasswd` file [here](./fixtures/gitserver/repos/.gitpasswd).

An Nginx running within the Git server container that uses the `.gitpasswd` to manage the authorizations.

##### Creating the SSH Key Pair

Executing the command `/usr/bin/ssh-keygen -t rsa -q -N "password" -f id_rsa -C "apenella@stevedore.test"` you can create an SSH key pair used for authentication in the example. It generates an RSA key pair with a passphrase or `password` and saves the private key in the `id_rsa` file and the public key in the `id_rsa.pub`.

```sh
root@6a88123e81cc:/$ /usr/bin/ssh-keygen  -t rsa -q -N "password" -f id_rsa -C "apenella@stevedore.test"
```

##### Creating the known_hosts File

This example creates a `known_hosts` file that enables seamless and secure access to Git repositories, eliminating the need for manual host key verification. By including this file, you can establish a connection to the Git server without encountering any host verification prompts, ensuring a smooth and secure workflow.

```sh
$ ssh-keyscan -H gitserver.stevedore.test > ~/.ssh/known_hosts
gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0  
```

##### Preparing the Git Repository and the Application's Source Code

The first step involves creating the bare repositories on the Git server and configuring the necessary permissions. This can be accomplished by following these commands:

```sh
/git/repos $ git config --global init.defaultBranch main
/git/repos $ git init --bare /git/repos/app1.git
/git/repos $ chown -R git:git app1.git/
```

The git config `--global init.defaultBranch main` command sets the default branch to `main` for all newly created repositories. The `git init --bare /git/repos/app1.git` command initializes the bare repository for `app1.git` in the desired location. Finally, the `chown -R git:git app1.git/` command assigns the appropriate ownership and permissions to the `app1.git` repository, ensuring proper access control.

Once the empty repository is ready, you can proceed to clone it and start writing the code for the application. Here is an example of the commands to perform these steps:

```sh
root@9840c0dae81a:/apps $ git clone git@gitserver.stevedore.test:/git/repos/app1.git
Cloning into 'app1'...
Enter passphrase for key '/root/.ssh/id_rsa':
warning: You appear to have cloned an empty repository.
root@9840c0dae81a:/apps $ cd app1/
root@9840c0dae81a:/apps/app1 $ git checkout -b main
Switched to a new branch 'main'
root@9840c0dae81a:/apps/app1 $ cat << EOF >> app.sh
> #!/bin/sh

echo "[user \$(whoami)] I'm app1"
> EOF
root@9840c0dae81a:/apps/app1 $ chmod +x app.sh
root@9840c0dae81a:/apps/app1 $ git add .
root@9840c0dae81a:/apps/app1 $ git commit -m "Create app1"
[main (root-commit) a6df89d] Create app1
 1 file changed, 3 insertions(+)
 create mode 100755 app.sh
root@9840c0dae81a:/apps/app1 $ git push origin main
Enter passphrase for key '/root/.ssh/id_rsa':
Enumerating objects: 3, done.
Counting objects: 100% (3/3), done.
Writing objects: 100% (3/3), 251 bytes | 251.00 KiB/s, done.
Total 3 (delta 0), reused 0 (delta 0), pack-reused 0
To gitserver.stevedore.test:/git/repos/app1.git
 * [new branch]      main -> main
```

The git clone command clones the empty repository from the Git server, prompting you for the passphrase of the SSH key. Once cloned, navigate into the `app1` directory and create a new branch named `main` using the `git checkout -b main` command. Proceed by creating the script `app.sh` with the desired code using a heredoc. Make the script executable with `chmod +x app.sh` and then commit the code by running `git commit -m "Create app1"`.
Next, push the changes to the repository using `git push origin main`. You are prompted for the passphrase of the SSH key. The output shows the progress of enumerating, counting, and writing objects, indicating the successful push to the main branch of the Git repository.

##### Clone the application 

To verify that the repository content is as expected, you can clone its content using the following command:

```sh
root@9840c0dae81a:/tmp $ git clone http://admin:admin@gitserver.stevedore.test:/git/repos/app1.git
Cloning into 'app1'...
warning: redirecting to http://gitserver.stevedore.test/git/repos/app1.git/
remote: Enumerating objects: 6, done.
remote: Counting objects: 100% (6/6), done.
remote: Compressing objects: 100% (4/4), done.
remote: Total 6 (delta 0), reused 0 (delta 0), pack-reused 0
```

This command clones the repository into a local directory named `app1`. This confirms that the repository has been successfully cloned.
