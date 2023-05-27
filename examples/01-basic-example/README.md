# Basic example
This example aims to introduce you to Stevedore and its basic concepts and commands. It follows the [quickstart guide](https://gostevedore.github.io/docs/getting-started/quickstart/) from the documentation, which serves as a starting point to get familiar with Stevedore. By following this example, you can learn how to create several image definitions into multiple files and a builder, add credentials to the credentials store, and build and promote a Docker image to a Docker registry using Stevedore.

- [Basic example](#basic-example)
  - [Requirements](#requirements)
  - [Stack](#stack)
  - [Usage](#usage)
  - [Expected Output](#expected-output)
    - [Starting the stack](#starting-the-stack)
    - [Getting Credentials](#getting-credentials)
    - [Getting Builders](#getting-builders)
    - [Getting images](#getting-images)
    - [Building images](#building-images)
    - [Cleaning the stack](#cleaning-the-stack)
  - [Additional information](#additional-information)
    - [Docker Registry](#docker-registry)
    - [Builders](#builders)
    - [Credentials](#credentials)
    - [Images](#images)

## Requirements
- Docker. _Tested on Docker server 20.10.21 and Docker API 1.41_
- Docker's Compose plugin or `docker-compose`. _Tested on Docker Compose version v2.17.3_
- `make` utility. _Tested on version 4.3-4.1build1_

## Stack
The stack required to run this example is defined in a [Docker Compose file](./docker-compose.yml). The stack consists of three services: a Docker Registry, a Docker Registry authorization and a Stevedore service. The Docker registry is used to store the Docker images built by Stevedore during the example execution. The Stevedore service is where the example is executed.

The Stevedore service is built from a container which is defined in that [Dockerfile](https://github.com/gostevedore/stevedore/blob/main/examples/01-basic-example/stack/stevedore/Dockerfile).

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
Below is the expected output for the `make run` command, which starts the Docker stack, gets some information about the Stevedore configuration, builds and promotes a Docker image using Stevedore, and then cleans the stack up.

### Starting the stack
```sh
Starting the stack to run basic-example

[+] Building 10.7s (21/21) FINISHED
 => [internal] load .dockerignore                                                                                                                                    0.0s
 => => transferring context: 2B                                                                                                                                      0.0s
 => [internal] load build definition from Dockerfile                                                                                                                 0.0s
 => => transferring dockerfile: 989B                                                                                                                                 0.0s
 => [internal] load metadata for docker.io/library/docker:20.10-dind                                                                                                 1.3s
 => [internal] load metadata for docker.io/library/golang:1.19-alpine                                                                                                1.9s
 => [golang 1/8] FROM docker.io/library/golang:1.19-alpine@sha256:9668643a2e62d8bd298ef3663a96de4a70ceb2865b9b7cadd1d5e08387745103                                   0.0s
 => => resolve docker.io/library/golang:1.19-alpine@sha256:9668643a2e62d8bd298ef3663a96de4a70ceb2865b9b7cadd1d5e08387745103                                          0.0s
 => [internal] load build context                                                                                                                                    0.0s
 => => transferring context: 94.05kB                                                                                                                                 0.0s
 => [stage-1 1/6] FROM docker.io/library/docker:20.10-dind@sha256:b848ea12a51f9be34b5ad6774a93a015fee1c2017d1896414c2f8fbaeb0c87d3                                   0.0s
 => => resolve docker.io/library/docker:20.10-dind@sha256:b848ea12a51f9be34b5ad6774a93a015fee1c2017d1896414c2f8fbaeb0c87d3                                           0.0s
 => CACHED [golang 2/8] WORKDIR /usr/src/app                                                                                                                         0.0s
 => CACHED [golang 3/8] RUN apk add --no-cache make build-base                                                                                                       0.0s
 => CACHED [golang 4/8] COPY go.mod ./                                                                                                                               0.0s
 => CACHED [golang 5/8] COPY go.sum ./                                                                                                                               0.0s
 => CACHED [golang 6/8] RUN go mod download && go mod verify                                                                                                         0.0s
 => [golang 7/8] COPY . ./                                                                                                                                           0.4s
 => [golang 8/8] RUN go build -ldflags "-s -w" -v -o /usr/local/bin/stevedore ./cmd/stevedore.go                                                                     5.5s
 => CACHED [stage-1 2/6] COPY --from=golang /usr/local/go /usr/local/go                                                                                              0.0s
 => [stage-1 3/6] COPY --from=golang /usr/local/bin/stevedore /usr/local/bin/stevedore                                                                               0.0s
 => [stage-1 4/6] WORKDIR /go                                                                                                                                        0.0s
 => [stage-1 5/6] RUN mkdir -p "/go/src" "/go/bin" && chmod -R 777 "/go"                                                                                             0.0s
 => [stage-1 6/6] COPY test/stack/client/entrypoint.sh /usr/local/bin/entrypoint.sh                                                                                  0.0s
 => exporting to docker image format                                                                                                                                 1.7s
 => => exporting layers                                                                                                                                              0.6s
 => => exporting manifest sha256:82d4394dc7e95ac3feca101b2c00bbac0421fac46d4388bfe1650bb021f8fcfb                                                                    0.0s
 => => exporting config sha256:7f4e00bd20f3a739646b7b9f9dbf589294ff7ecd56c276b535e4a47f0bdf9119                                                                      0.0s
 => => sending tarball                                                                                                                                               1.1s
 => importing to docker                                                                                                                                              0.2s
[+] Running 4/4
 ✔ Network basic-example_default         Created                                                                                                                     0.1s
 ✔ Container basic-example-stevedore-1   Started                                                                                                                     0.5s
 ✔ Container basic-example-dockerauth-1  Started                                                                                                                     0.5s
 ✔ Container basic-example-registry-1    Started                                                                                                                     0.7s
```

### Getting Credentials
To obtain the credentials, use the command: `stevedore get credentials`.

```sh
 Run example basic-example

 [basic-example] Get credentials
 Waiting for dockerd to be ready...
/certs/server/cert.pem: OK
/certs/client/cert.pem: OK
 Waiting for dockerd to be ready...
ID                      TYPE              CREDENTIALS
registry.stevedore.test username-password username=admin
```

### Getting Builders
Use the command `stevedore get builders` to get the builders.

```sh
 [basic-example] Get builders

 Waiting for dockerd to be ready...
/certs/server/cert.pem: OK
/certs/client/cert.pem: OK
 Waiting for dockerd to be ready...
NAME   DRIVER
my-app docker
```

### Getting images
To view the images in tree format, run `stevedore get images --tree`.

```sh
 [basic-example] Get images
 Waiting for dockerd to be ready...
/certs/server/cert.pem: OK
/certs/client/cert.pem: OK
 Waiting for dockerd to be ready...
├─── busybox:1.35
│  ├─── registry.stevedore.test/my-app:2.1.0-busybox1.35
│  ├─── registry.stevedore.test/my-app:3.2.1-busybox1.35
├─── busybox:1.36
│  ├─── registry.stevedore.test/my-app:3.2.1-busybox1.36
```

### Building images
The example uses the command `stevedore build my-app --push-after-build` to build and automatically promote the images to the Docker registry. Because the three defined images are being built at the same time, the output shows these outputs mixed.

If you review the builder definition, you will notice that the source code for `my-app` is located in the [./builders/apps.yaml](builders/apps.yaml) file. This folder contains the necessary resources required for building the `my-app` Docker image.

```sh
 [basic-example] Build my-app and push images after build
 Waiting for dockerd to be ready...
 /certs/server/cert.pem: OK
 /certs/client/cert.pem: OK
 Waiting for dockerd to be ready...
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 1/7 : ARG image_from_name
registry.stevedore.test/my-app:2.1.0-busybox1.35 Step 1/7 : ARG image_from_name
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 2/7 : ARG image_from_tag
registry.stevedore.test/my-app:2.1.0-busybox1.35 Step 2/7 : ARG image_from_tag
registry.stevedore.test/my-app:2.1.0-busybox1.35 Step 3/7 : FROM ${image_from_name}:${image_from_tag}
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 3/7 : FROM ${image_from_name}:${image_from_tag}
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 1/7 : ARG image_from_name
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 2/7 : ARG image_from_tag
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 3/7 : FROM ${image_from_name}:${image_from_tag}
registry.stevedore.test/my-app:2.1.0-busybox1.35 ‣  1.35:  Pulling from library/busybox
registry.stevedore.test/my-app:2.1.0-busybox1.35 ‣  880dcab25a96:  Pull complete
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  Digest: sha256:223ae047b1065bd069aac01ae3ac8088b3ca4a527827e283b85112f29385fb1b
registry.stevedore.test/my-app:2.1.0-busybox1.35 ‣  Digest: sha256:223ae047b1065bd069aac01ae3ac8088b3ca4a527827e283b85112f29385fb1b
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  Status: Downloaded newer image for busybox:1.35
registry.stevedore.test/my-app:2.1.0-busybox1.35 ‣  Status: Downloaded newer image for busybox:1.35
registry.stevedore.test/my-app:2.1.0-busybox1.35 ---> 12b6f68a826b
registry.stevedore.test/my-app:2.1.0-busybox1.35 Step 4/7 : ARG whoami=unknown
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> 12b6f68a826b
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 4/7 : ARG whoami=unknown
registry.stevedore.test/my-app:2.1.0-busybox1.35 ---> Running in 178523232c17
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  1.36:  Pulling from library/busybox
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  4b35f584bb4f:  Pull complete
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  Digest: sha256:b5d6fe0712636ceb7430189de28819e195e8966372edfc2d9409d79402a0dc16
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  Status: Downloaded newer image for busybox:1.36
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> 7cfbbec8963d
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 4/7 : ARG whoami=unknown
registry.stevedore.test/my-app:2.1.0-busybox1.35 ---> 50eb222ba4dd
registry.stevedore.test/my-app:2.1.0-busybox1.35 Step 5/7 : RUN echo "Hey there, I'm ${whoami}!" > /whoami.txt
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> Running in f11a6eb03c94
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> 03d5c0e6221d
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 5/7 : RUN echo "Hey there, I'm ${whoami}!" > /whoami.txt
registry.stevedore.test/my-app:2.1.0-busybox1.35 ---> Running in 69fa230982cf
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> Running in 3d76a3ff01ff
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> 4a7667b9e405
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 5/7 : RUN echo "Hey there, I'm ${whoami}!" > /whoami.txt
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> Running in f402e66aae25
registry.stevedore.test/my-app:2.1.0-busybox1.35 ---> 8c004675c87e
registry.stevedore.test/my-app:2.1.0-busybox1.35 Step 6/7 : CMD ["cat","/whoami.txt"]
registry.stevedore.test/my-app:2.1.0-busybox1.35 ---> Running in 6691d9e37ae6
registry.stevedore.test/my-app:2.1.0-busybox1.35 ---> 00d6c9b88dc7
registry.stevedore.test/my-app:2.1.0-busybox1.35 Step 7/7 : LABEL created_at=2023-05-08T06:19:28.936204297Z
registry.stevedore.test/my-app:2.1.0-busybox1.35 ---> Running in a5afd4d63240
registry.stevedore.test/my-app:2.1.0-busybox1.35 ---> 3c0a9004bd64
registry.stevedore.test/my-app:2.1.0-busybox1.35  ‣ sha256:3c0a9004bd6467786d660310ad7618d2934385169b74e10cd548f4a9ef612f0e
registry.stevedore.test/my-app:2.1.0-busybox1.35 Successfully built 3c0a9004bd64
registry.stevedore.test/my-app:2.1.0-busybox1.35 Successfully tagged registry.stevedore.test/my-app:2.1.0-busybox1.35
```

Here you can see how Stevedore starts to push images automatically.
```sh
registry.stevedore.test/my-app:2.1.0-busybox1.35 ‣  The push refers to repository [registry.stevedore.test/my-app]
registry.stevedore.test/my-app:2.1.0-busybox1.35 ‣  cb6098d4c396:  Pushed
registry.stevedore.test/my-app:2.1.0-busybox1.35 ‣  f9b26e1dcefb:  Pushing [============================>                      ]  2.734MB/4.859MB
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> 08a0654be778
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 6/7 : CMD ["cat","/whoami.txt"]
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> Running in 4650b3b8234b
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> 4ccd37b93885
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 7/7 : LABEL created_at=2023-05-08T06:19:28.935611217Z
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> Running in d3da5b431c3d
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> 95e7efdbd5d2
registry.stevedore.test/my-app:3.2.1-busybox1.36  ‣ sha256:95e7efdbd5d25c746d4a4ff841758120b873423996be9826aff148d15ccbc40c
registry.stevedore.test/my-app:3.2.1-busybox1.36 Successfully built 95e7efdbd5d2
registry.stevedore.test/my-app:3.2.1-busybox1.36 Successfully tagged registry.stevedore.test/my-app:3.2.1-busybox1.36
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  The push refers to repository [registry.stevedore.test/my-app]
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> 738c4666ab0d
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 6/7 : CMD ["cat","/whoami.txt"]
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  d2f9965848ab:  Preparing
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  d2f9965848ab:  Pushed
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  baacf561cfff:  Pushing [>                                                  ]  66.56kB/4.863MB
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> da5bb1a2adbf
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 7/7 : LABEL created_at=2023-05-08T06:19:28.936204297Z
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> Running in bd88293b7374
registry.stevedore.test/my-app:2.1.0-busybox1.35 ‣  2.1.0-busybox1.35: digest: sha256:0bcfe06148e16fedd343c73e751ce00ee098f9bb80c1a867aa2bf97e17bdc44f size: 735
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> 695fe33875ed
registry.stevedore.test/my-app:3.2.1-busybox1.35  ‣ sha256:695fe33875ed6dfd3641f0c05f17e17a51dc9cc129dd5143dc41150167d8d20d
registry.stevedore.test/my-app:3.2.1-busybox1.35 Successfully built 695fe33875ed
registry.stevedore.test/my-app:3.2.1-busybox1.35 Successfully tagged registry.stevedore.test/my-app:3.2.1-busybox1.35
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  The push refers to repository [registry.stevedore.test/my-app]
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  d2f9965848ab:  Pushed
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  d2f9965848ab:  Pushed
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  baacf561cfff:  Pushed
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  3.2.1-busybox1.36: digest: sha256:1e0521ea9319a2bd1dd56eb87ff39bda6b78566bc76cd451d451627729a7c99b size: 735
```

### Cleaning the stack
```sh
Stopping the stack to run basic-example

[+] Running 8/8
 ✔ Container basic-example-stevedore-1                 Removed                                                                                                      10.2s
 ✔ Container basic-example-registry-1                  Removed                                                                                                       0.3s
 ✔ Container basic-example-stevedore-run-e1a847da670a  Removed                                                                                                       0.0s
 ✔ Container basic-example-stevedore-run-2060c0f06724  Removed                                                                                                       0.0s
 ✔ Container basic-example-stevedore-run-00680eb3c6e9  Removed                                                                                                       0.0s
 ✔ Container basic-example-stevedore-run-78389cc51f7a  Removed                                                                                                       0.0s
 ✔ Container basic-example-dockerauth-1                Removed                                                                                                       0.2s
 ✔ Network basic-example_default                       Removed                                                                                                       0.4s
```

## Additional information
In addition to the core steps outlined in the example, the following section provides additional information and insights to further enhance your understanding of how this example uses Stevedore.

### Docker Registry
To illustrate how Stevedore can push Docker images after building them, this example starts a Docker registry. The Docker registry is launched as a Docker Compose service, accompanied by an authentication service. In this example, Cesanta's [docker_auth](https://github.com/cesanta/docker_auth) is utilized.

The Docker registry service is accessible internally through the host `registry.stevedore.test`. To log in to the Docker registry service, you can use the credentials: user=`admin`, password=`admin`. Additionally, the authentication service is accessible through the host `auth.stevedore.test`.

### Builders
The example uses a [global builder](https://gostevedore.github.io/docs/reference-guide/builder/#global-builder), a builder that can be applied to any image definition.

### Credentials
The current example uses the credentials [local store](https://gostevedore.github.io/docs/reference-guide/credentials/credentials-store/#local-storage). It is the default backend storage for the credentials which stores them locally, on your local system disk.
You can see the credentials store configuration on Stevedore's configuration file [./stevedore.yaml](stevedore.yaml), within the example's folder.

The example already provides the `./credentials` folder with the credentials to log in to the Docker registry available.
Using the command `stevedore get credentials` and the `--show-secrets` flag, you can see all the details about the stored credentials.
```sh
/app/examples/01-basic-example # stevedore get credentials --show-secrets
ID                      TYPE              CREDENTIALS
registry.stevedore.test username-password username=admin, password=admin
```

The credentials already present there were generated by the command:
```sh
stevedore create credentials registry.stevedore.test --username admin
```

### Images
The example includes image definitions for different versions of the application `my-app`. These image definitions are being built from the Docker images `busybox:1.35` and `busybox:1.36`. You can find these foundational image definitions in the file [./images/foundational.yaml](images/foundational.yaml).

It is important to note that the foundational image definitions include a persistent label named `created_at`, which is inherited by all images defined below them.

For the specific image definitions of `my-app`, you can refer to the file [./images/applications.yaml](images/applications.yaml).
