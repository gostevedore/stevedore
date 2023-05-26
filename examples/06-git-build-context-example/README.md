# Git Build Context Example
This example illustrates how to utilize the [Git build context](https://gostevedore.github.io/docs/reference-guide/builder/docker/#git-context) feature in Stevedore. It demonstrates the ability to specify a Git repository as the build context, allowing you to directly build images from a specific branch, tag, or commit.

- [Git Build Context Example](#git-build-context-example)
  - [Requirements](#requirements)
  - [Stack](#stack)
  - [Usage](#usage)
  - [Expected Output](#expected-output)
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
        - [Preparing the Application Source Code Repository code](#preparing-the-application-source-code-repository-code)
        - [Clone the application](#clone-the-application)


## Requirements
- Docker. _Tested on Docker server 20.10.21 and Docker API 1.41_
- Docker's Compose plugin or `docker-compose`. _Tested on Docker Compose version v2.17.3_
- `make` utility. _Tested on version 4.3-4.1build1_

## Stack
The stack required to run this example is defined in a [Docker Compose file](./docker-compose.yml). The stack consists of four services: a Docker Registry, a Docker Registry authorization, a Git server and a Stevedore service.

The Docker registry is used to store the Docker images built by Stevedore during the example execution.
The Git server service is used to store the source code of the applications that we pretend to build during this example.
The Stevedore service is where the example is executed and it is built from a container which is defined in the [Dockerfile](https://github.com/gostevedore/stevedore/blob/main/test/stack/client/Dockerfile) present in the `test/stack/client` directory.

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

## Expected Output
Below is the expected output for the `make run` command, which starts the Docker stack, gets some information about the Stevedore configuration, builds and promotes Docker images using Stevedore, and then cleans the stack up.

```sh
Starting the stack to run 06-git-build-context-example

[+] Building 13.5s (33/33) FINISHED
 => [gitserver internal] load build definition from Dockerfile                                                                                                       0.0s
 => => transferring dockerfile: 1.41kB                                                                                                                               0.0s
 => [gitserver internal] load .dockerignore                                                                                                                          0.0s
 => => transferring context: 2B                                                                                                                                      0.0s
 => [gitserver internal] load metadata for docker.io/library/alpine:3.16                                                                                             5.5s
 => [gitserver internal] load build context                                                                                                                          0.0s
 => => transferring context: 96B                                                                                                                                     0.0s
 => [gitserver 1/6] FROM docker.io/library/alpine:3.16@sha256:c2b622f6e510a0d25bccaffa9e67b75a6860cb09b74bb58cfc36a9ef4331109f                                       0.0s
 => => resolve docker.io/library/alpine:3.16@sha256:c2b622f6e510a0d25bccaffa9e67b75a6860cb09b74bb58cfc36a9ef4331109f                                                 0.0s
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
 => [gitserver gitserver] importing to docker                                                                                                                        0.0s
 => [stevedore internal] load build definition from Dockerfile                                                                                                       0.0s
 => => transferring dockerfile: 989B                                                                                                                                 0.0s
 => [stevedore internal] load .dockerignore                                                                                                                          0.0s
 => => transferring context: 2B                                                                                                                                      0.0s
 => [stevedore internal] load metadata for docker.io/library/docker:20.10-dind                                                                                       0.2s
 => [stevedore internal] load metadata for docker.io/library/golang:1.19-alpine                                                                                      0.4s
 => [stevedore internal] load build context                                                                                                                          0.0s
 => => transferring context: 218.87kB                                                                                                                                0.0s
 => [stevedore stage-1 1/6] FROM docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                         0.0s
 => => resolve docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                                           0.0s
 => [stevedore golang 1/8] FROM docker.io/library/golang:1.19-alpine@sha256:4147c2ad24a347e53a8a361b663dcf82fca4157a3f9a5136d696a4e53dd22b65                         0.0s
 => => resolve docker.io/library/golang:1.19-alpine@sha256:4147c2ad24a347e53a8a361b663dcf82fca4157a3f9a5136d696a4e53dd22b65                                          0.0s
 => CACHED [stevedore golang 2/8] WORKDIR /usr/src/app                                                                                                               0.0s
 => CACHED [stevedore golang 3/8] RUN apk add --no-cache make build-base                                                                                             0.0s
 => CACHED [stevedore golang 4/8] COPY go.mod ./                                                                                                                     0.0s
 => CACHED [stevedore golang 5/8] COPY go.sum ./                                                                                                                     0.0s
 => CACHED [stevedore golang 6/8] RUN go mod download && go mod verify                                                                                               0.0s
 => [stevedore golang 7/8] COPY . ./                                                                                                                                 0.4s
 => [stevedore golang 8/8] RUN go build -ldflags "-s -w" -v -o /usr/local/bin/stevedore ./cmd/stevedore.go                                                           4.8s
 => CACHED [stevedore stage-1 2/6] COPY --from=golang /usr/local/go /usr/local/go                                                                                    0.0s
 => CACHED [stevedore stage-1 3/6] COPY --from=golang /usr/local/bin/stevedore /usr/local/bin/stevedore                                                              0.0s
 => CACHED [stevedore stage-1 4/6] WORKDIR /go                                                                                                                       0.0s
 => CACHED [stevedore stage-1 5/6] RUN mkdir -p "/go/src" "/go/bin" && chmod -R 777 "/go"                                                                            0.0s
 => CACHED [stevedore stage-1 6/6] COPY test/stack/client/entrypoint.sh /usr/local/bin/entrypoint.sh                                                                 0.0s
 => [stevedore] exporting to docker image format                                                                                                                     0.9s
 => => exporting layers                                                                                                                                              0.0s
 => => exporting manifest sha256:4fccbb1e18be27632f7bfbf50912c465a8dcb6535fa0ac41021ef76070c378a0                                                                    0.0s
 => => exporting config sha256:19e63d4df31a13d448f723680c52828091eaeaa47079c1438da20fefbe9d855b                                                                      0.0s
 => => sending tarball                                                                                                                                               0.9s
 => [stevedore stevedore] importing to docker                                                                                                                        0.0s
[+] Running 6/6
 ✔ Network 06-git-build-context-example_default         Created                                                                                                      0.1s
 ✔ Container 06-git-build-context-example-worker-1      Started                                                                                                      0.6s
 ✔ Container 06-git-build-context-example-gitserver-1   Started                                                                                                      0.6s
 ✔ Container 06-git-build-context-example-dockerauth-1  Started                                                                                                      0.6s
 ✔ Container 06-git-build-context-example-registry-1    Started                                                                                                      1.1s
 ✔ Container 06-git-build-context-example-stevedore-1   Started                                                                                                      1.2s
 ```

### Getting images
To view the images in tree format, run `stevedore get images --tree`.

```sh
 Run example 06-git-build-context-example

 [06-git-build-context-example] Get images
[+] Building 0.0s (0/0)
[+] Creating 1/1
 ✔ Container 06-git-build-context-example-gitserver-1  Recreated                                                                                                    10.4s
[+] Running 1/1
 ✔ Container 06-git-build-context-example-gitserver-1  Started                                                                                                       0.3s
[+] Building 0.0s (0/0)
# gitserver:22 SSH-2.0-OpenSSH_9.0
# gitserver:22 SSH-2.0-OpenSSH_9.0
# gitserver:22 SSH-2.0-OpenSSH_9.0
# gitserver:22 SSH-2.0-OpenSSH_9.0
# gitserver:22 SSH-2.0-OpenSSH_9.0
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
/certs/server/cert.pem: OK
 Waiting for dockerd to be ready...
/certs/client/cert.pem: OK
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
 [06-git-build-context-example] Build the base image and its descendaant, and push the images after build
[+] Building 0.0s (0/0)
[+] Creating 1/0
 ✔ Container 06-git-build-context-example-gitserver-1  Running                                                                                                       0.0s
[+] Building 0.0s (0/0)
# gitserver:22 SSH-2.0-OpenSSH_9.0
# gitserver:22 SSH-2.0-OpenSSH_9.0
# gitserver:22 SSH-2.0-OpenSSH_9.0
# gitserver:22 SSH-2.0-OpenSSH_9.0
# gitserver:22 SSH-2.0-OpenSSH_9.0
 Waiting for dockerd to be ready...
/certs/server/cert.pem: OK
 Waiting for dockerd to be ready...
/certs/client/cert.pem: OK
registry.stevedore.test/base:busybox-1.35 Step 1/7 : ARG image_from_name
registry.stevedore.test/base:busybox-1.35 Step 2/7 : ARG image_from_tag
registry.stevedore.test/base:busybox-1.35 Step 3/7 : FROM ${image_from_name}:${image_from_tag}
registry.stevedore.test/base:busybox-1.36 Step 1/7 : ARG image_from_name
registry.stevedore.test/base:busybox-1.36 Step 2/7 : ARG image_from_tag
registry.stevedore.test/base:busybox-1.36 Step 3/7 : FROM ${image_from_name}:${image_from_tag}
registry.stevedore.test/base:busybox-1.36 ‣  1.36:  Pulling from library/busybox
registry.stevedore.test/base:busybox-1.36 ‣  325d69979d33:  Pull complete
registry.stevedore.test/base:busybox-1.36 ‣  Digest: sha256:560af6915bfc8d7630e50e212e08242d37b63bd5c1ccf9bd4acccf116e262d5b
registry.stevedore.test/base:busybox-1.36 ‣  Status: Downloaded newer image for busybox:1.36
registry.stevedore.test/base:busybox-1.36 ---> 8135583d97fe
registry.stevedore.test/base:busybox-1.36 Step 4/7 : RUN echo "anonymous:x:10001:10001:,,,:/app:/bin/sh" >> /etc/passwd &&     echo "anonymous:x:10001:" >> /etc/group &&     mkdir -p /app &&     chown 10001:10001 /app
registry.stevedore.test/base:busybox-1.36 ---> Running in 370a35bb6c83
registry.stevedore.test/base:busybox-1.36 ---> 22e37b81fe8c
registry.stevedore.test/base:busybox-1.36 Step 5/7 : USER anonymous
registry.stevedore.test/base:busybox-1.36 ---> Running in 73c3521c3031
registry.stevedore.test/base:busybox-1.36 ---> 6a3c1bc863b2
registry.stevedore.test/base:busybox-1.36 Step 6/7 : WORKDIR /app
registry.stevedore.test/base:busybox-1.36 ---> Running in 867392b8528c
registry.stevedore.test/base:busybox-1.36 ---> 266ecfbcf5fa
registry.stevedore.test/base:busybox-1.36 Step 7/7 : LABEL created_at=2023-05-23T16:14:59.043679729Z
registry.stevedore.test/base:busybox-1.36 ---> Running in f10235741c70
registry.stevedore.test/base:busybox-1.36 ---> 91e6f924c3fb
registry.stevedore.test/base:busybox-1.36  ‣ sha256:91e6f924c3fb41d6c92e9ecb597e67aa9e54d5df935fe9314094ac8c2c796dd5
registry.stevedore.test/base:busybox-1.36 Successfully built 91e6f924c3fb
registry.stevedore.test/base:busybox-1.36 Successfully tagged registry.stevedore.test/base:busybox-1.36
registry.stevedore.test/base:busybox-1.36 ‣  The push refers to repository [registry.stevedore.test/base]
registry.stevedore.test/base:busybox-1.36 ‣  14a9c01685cb:  Pushed
registry.stevedore.test/base:busybox-1.36 ‣  9547b4c33213:  Pushed
registry.stevedore.test/base:busybox-1.36 ‣  busybox-1.36: digest: sha256:0259c988d53535a243c4135f444edc4622d45c1d2b1ed8ae9818525d0ec7e4b4 size: 735
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 1/7 : ARG image_from_name
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 2/7 : ARG image_from_tag
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 3/7 : ARG image_from_registry_host
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 4/7 : FROM ${image_from_registry_host}/${image_from_name}:${image_from_tag}
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> 91e6f924c3fb
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 5/7 : COPY ./${app_name}/app.sh /app/run.sh
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> 76a2921067fe
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 6/7 : CMD ["/app/run.sh"]
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> Running in 2f538893af8e
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> 6eb1a5539e9c
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 7/7 : LABEL created_at=2023-05-23T16:14:59.043679729Z
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> Running in 08431b029d5d
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> cae9264d8769
registry.stevedore.test/app3:v1-base-busybox-1.36  ‣ sha256:cae9264d8769272112eb72696123239164639d3341af64df145c793423a70e47
registry.stevedore.test/app3:v1-base-busybox-1.36 Successfully built cae9264d8769
registry.stevedore.test/app3:v1-base-busybox-1.36 Successfully tagged registry.stevedore.test/app3:v1-base-busybox-1.36
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  The push refers to repository [registry.stevedore.test/app3]
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  325fa5d553b3:  Pushed
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  14a9c01685cb:  Mounted from base
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  9547b4c33213:  Mounted from base
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 1/7 : ARG image_from_name
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 2/7 : ARG image_from_tag
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 3/7 : ARG image_from_registry_host
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 4/7 : FROM ${image_from_registry_host}/${image_from_name}:${image_from_tag}
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> 91e6f924c3fb
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 5/7 : COPY ./${app_name}/app.sh /app/run.sh
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  v1-base-busybox-1.36: digest: sha256:41b67331b4961ed525aa3b98293d73ef2def20d0f9f29742300dbf3cfc1efdc8 size: 942
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> 84fa8c05368b
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 6/7 : CMD ["/app/run.sh"]
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> Running in 1672c065c001
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> fc55544dc1f3
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 7/7 : LABEL created_at=2023-05-23T16:14:59.043679729Z
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> Running in cffbc8e08be1
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> 02ad7ca8fd37
registry.stevedore.test/app2:v1-base-busybox-1.36  ‣ sha256:02ad7ca8fd37df07a9ac029db3010c2a07b4ab8d2cf41335f34048ea01ed9c60
registry.stevedore.test/app2:v1-base-busybox-1.36 Successfully built 02ad7ca8fd37
registry.stevedore.test/app2:v1-base-busybox-1.36 Successfully tagged registry.stevedore.test/app2:v1-base-busybox-1.36
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  The push refers to repository [registry.stevedore.test/app2]
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  9bacd26a0896:  Pushed
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  14a9c01685cb:  Mounted from app3
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  9547b4c33213:  Mounted from app3
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  v1-base-busybox-1.36: digest: sha256:48bd716e778d18f412572e2d746b326dbcaa0433ea609f3da0e2f3873a296145 size: 942
registry.stevedore.test/base:busybox-1.35 ‣  1.35:  Pulling from library/busybox
registry.stevedore.test/base:busybox-1.35 ‣  c15cbdab5f8e:  Pull complete
registry.stevedore.test/base:busybox-1.35 ‣  Digest: sha256:b4e4a06de46acc0958cd93e2eeb769077d255f06a7c3a91196509c16b7bc989e
registry.stevedore.test/base:busybox-1.35 ‣  Status: Downloaded newer image for busybox:1.35
registry.stevedore.test/base:busybox-1.35 ---> dddc7578369a
registry.stevedore.test/base:busybox-1.35 Step 4/7 : RUN echo "anonymous:x:10001:10001:,,,:/app:/bin/sh" >> /etc/passwd &&     echo "anonymous:x:10001:" >> /etc/group &&     mkdir -p /app &&     chown 10001:10001 /app
registry.stevedore.test/base:busybox-1.35 ---> Running in e30d2d089c17
registry.stevedore.test/base:busybox-1.35 ---> 34d524b928c0
registry.stevedore.test/base:busybox-1.35 Step 5/7 : USER anonymous
registry.stevedore.test/base:busybox-1.35 ---> Running in 2a552e33da99
registry.stevedore.test/base:busybox-1.35 ---> 85d20392a1cb
registry.stevedore.test/base:busybox-1.35 Step 6/7 : WORKDIR /app
registry.stevedore.test/base:busybox-1.35 ---> Running in 4181365b1441
registry.stevedore.test/base:busybox-1.35 ---> 1e838f1213a3
registry.stevedore.test/base:busybox-1.35 Step 7/7 : LABEL created_at=2023-05-23T16:14:59.044265696Z
registry.stevedore.test/base:busybox-1.35 ---> Running in b17c4d9e7647
registry.stevedore.test/base:busybox-1.35 ---> 4893e409d6fb
registry.stevedore.test/base:busybox-1.35  ‣ sha256:4893e409d6fb65196ec422d3f5156a121ebff78be4427f2b6f503c11a161aff9
registry.stevedore.test/base:busybox-1.35 Successfully built 4893e409d6fb
registry.stevedore.test/base:busybox-1.35 Successfully tagged registry.stevedore.test/base:busybox-1.35
registry.stevedore.test/base:busybox-1.35 ‣  The push refers to repository [registry.stevedore.test/base]
registry.stevedore.test/base:busybox-1.35 ‣  9ef89241914f:  Pushed
registry.stevedore.test/base:busybox-1.35 ‣  42ef21f45b9a:  Pushed
registry.stevedore.test/base:busybox-1.35 ‣  busybox-1.35: digest: sha256:d151271cedfb1de502d4cfee11cce79b13cec61e1076f63163379564efa36208 size: 735
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 1/7 : ARG image_from_name
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 2/7 : ARG image_from_tag
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 3/7 : ARG image_from_registry_host
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 4/7 : FROM ${image_from_registry_host}/${image_from_name}:${image_from_tag}
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> 4893e409d6fb
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 5/7 : COPY ./${app_name}/app.sh /app/run.sh
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> 6c3fbab23afd
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 6/7 : CMD ["/app/run.sh"]
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> Running in 6ce0fd8c0f50
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> ad9b04491ff5
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 7/7 : LABEL created_at=2023-05-23T16:14:59.044265696Z
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> Running in d8747e3842e6
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> b539d0f618aa
registry.stevedore.test/app1:v1-base-busybox-1.35  ‣ sha256:b539d0f618aa9801943dee45fe2553c171cb8ef64da420d98709be2ae0149d30
registry.stevedore.test/app1:v1-base-busybox-1.35 Successfully built b539d0f618aa
registry.stevedore.test/app1:v1-base-busybox-1.35 Successfully tagged registry.stevedore.test/app1:v1-base-busybox-1.35
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  The push refers to repository [registry.stevedore.test/app1]
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  55b603b59147:  Pushed
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  9ef89241914f:  Mounted from base
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  42ef21f45b9a:  Mounted from base
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  v1-base-busybox-1.35: digest: sha256:153a184f64e62b942fc3f1f595db69f103ab7342a4fbec3e5a78542f527e294c size: 942
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 1/7 : ARG image_from_name
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 2/7 : ARG image_from_tag
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 3/7 : ARG image_from_registry_host
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 4/7 : FROM ${image_from_registry_host}/${image_from_name}:${image_from_tag}
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> 4893e409d6fb
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 5/7 : COPY ./${app_name}/app.sh /app/run.sh
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> 8324840d198e
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 6/7 : CMD ["/app/run.sh"]
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> Running in d94074c2b986
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> 97247f28394b
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 7/7 : LABEL created_at=2023-05-23T16:14:59.044265696Z
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> Running in 37e8d81cc919
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> bbfd64928232
registry.stevedore.test/app2:v1-base-busybox-1.35  ‣ sha256:bbfd6492823204de670b9ae6193c4de7fbe812edab3509f62d2b24d4cf6aa595
registry.stevedore.test/app2:v1-base-busybox-1.35 Successfully built bbfd64928232
registry.stevedore.test/app2:v1-base-busybox-1.35 Successfully tagged registry.stevedore.test/app2:v1-base-busybox-1.35
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  The push refers to repository [registry.stevedore.test/app2]
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  d22a8ac4e949:  Pushed
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  9ef89241914f:  Mounted from app1
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  42ef21f45b9a:  Mounted from app1
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  v1-base-busybox-1.35: digest: sha256:953e0442ffc3441362bb37fe9e7d8ce487f836cb59768795302a65c5a7b413fd size: 942
```

### Cleaning the stack
```sh
Stopping the stack to run 06-git-build-context-example

[+] Running 8/8
 ✔ Container 06-git-build-context-example-worker-1                    Removed                                                                                        0.0s
 ✔ Container 06-git-build-context-example-stevedore-run-9bf1e474dd60  Removed                                                                                        0.0s
 ✔ Container 06-git-build-context-example-stevedore-1                 Removed                                                                                       10.3s
 ✔ Container 06-git-build-context-example-registry-1                  Removed                                                                                        0.3s
 ✔ Container 06-git-build-context-example-stevedore-run-a51dc04cf27b  Removed                                                                                        0.0s
 ✔ Container 06-git-build-context-example-dockerauth-1                Removed                                                                                        0.2s
 ✔ Container 06-git-build-context-example-gitserver-1                 Removed                                                                                       10.3s
 ✔ Network 06-git-build-context-example_default                       Removed                                                                                        0.4s
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
Access to the Git server is available via the host `gitserver.stevedeore.test`. You can authenticate using either the basic auth method with the credentials user=`admin` and password=`admin` or by utilizing SSH keys. If using SSH keys, the necessary keys are stored in the [fixtures](./fixtures/ssh) folder as well. It is important to note that the private key is password-protected, with the password set to `password`.

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

An Nginx running within the Git server container that uses the `.gitpasswd` fix to manage the authorizations.

##### Creating the SSH Key Pair
Executing the command `/usr/bin/ssh-keygen -t rsa -q -N "password" -f id_rsa -C "apenella@stevedore.test"` you can create an SSH key pair used for authentication in the example. It generates an RSA key pair with a passphrase or `password` and saves the private key in the `id_rsa` file and the public key in the `id_rsa.pub`.
```sh
root@6a88123e81cc:/$ /usr/bin/ssh-keygen  -t rsa -q -N "password" -f id_rsa -C "apenella@stevedore.test"
```

You can find the used keys [here](./fixtures/ssh/).

##### Creating the known_hosts File
This example provides a preconfigured `known_hosts` file that enables seamless and secure access to Git repositories, eliminating the need for manual host key verification. By including this file, you can establish a connection to the Git server without encountering any host verification prompts, ensuring a smooth and secure workflow.

```sh
$ ssh-keyscan -H gitserver.stevedore.test > ~/.ssh/known_hosts
gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0  
```

##### Preparing the Application Source Code Repository code
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
