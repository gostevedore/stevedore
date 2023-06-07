# Promote Images Example

This example showcases the image promotion feature in Stevedore, demonstrating how to promote Docker images after building them. Here you can see how to promote a Docker images in multiple way, using the semantic version tags generation, overwriting the namespace of the Docker registry or adding additional tags.

- [Promote Images Example](#promote-images-example)
  - [Requirements](#requirements)
  - [Stack](#stack)
  - [Usage](#usage)
  - [Example Execution Insights](#example-execution-insights)
    - [Starting the stack](#starting-the-stack)
    - [Waiting for Dockerd to be Ready](#waiting-for-dockerd-to-be-ready)
    - [Building the base image](#building-the-base-image)
    - [Promoting the base image](#promoting-the-base-image)
    - [Building the app1 image](#building-the-app1-image)
    - [Promoting the app1 image](#promoting-the-app1-image)
    - [Cleaning the stack](#cleaning-the-stack)

## Requirements
- Docker. _Tested on Docker server 20.10.21 and Docker API 1.41_
- Docker's Compose plugin or `docker-compose`. _Tested on Docker Compose version v2.17.3_
- `make` utility. _Tested on version 4.3-4.1build1_

## Stack
The stack required to run this example is defined in a [Docker Compose file](./docker-compose.yml). The stack consists of three services: a Docker Registry, a Docker Registry authorization and a Stevedore service. The Docker registry is used to store the Docker images built by Stevedore during the example execution. The Stevedore service is where the example is executed.

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
Below is the expected output for the `make run` command, which starts the Docker stack, gets some information about the Stevedore configuration, builds and promotes a Docker image using Stevedore, and then cleans the stack up.

### Starting the stack
```sh
Starting the stack to run 11-promote-images-example

[+] Building 7.6s (18/18) FINISHED
 => [stevedore internal] load .dockerignore                                                                                                                          0.0s
 => => transferring context: 2B                                                                                                                                      0.0s
 => [stevedore internal] load build definition from Dockerfile                                                                                                       0.0s
 => => transferring dockerfile: 888B                                                                                                                                 0.0s
 => [stevedore internal] load metadata for docker.io/library/docker:20.10-dind                                                                                       0.5s
 => [stevedore internal] load metadata for docker.io/library/golang:1.19-alpine                                                                                      0.5s
 => [stevedore golang 1/7] FROM docker.io/library/golang:1.19-alpine@sha256:106e12a19d4f8360c7dea4c08ef6ab7b62ec153972c77e099998eff9cb87e4f0                         0.0s
 => => resolve docker.io/library/golang:1.19-alpine@sha256:106e12a19d4f8360c7dea4c08ef6ab7b62ec153972c77e099998eff9cb87e4f0                                          0.0s
 => CACHED [stevedore stage-1 1/4] FROM docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                  0.0s
 => => resolve docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                                           0.0s
 => [stevedore internal] load build context                                                                                                                          0.1s
 => => transferring context: 342.07kB                                                                                                                                0.1s
 => CACHED [stevedore golang 2/7] WORKDIR /usr/src/app                                                                                                               0.0s
 => CACHED [stevedore golang 3/7] RUN apk add --no-cache make build-base                                                                                             0.0s
 => CACHED [stevedore golang 4/7] COPY go.mod go.sum ./                                                                                                              0.0s
 => CACHED [stevedore golang 5/7] RUN go mod download                                                                                                                0.0s
 => [stevedore golang 6/7] COPY . ./                                                                                                                                 0.5s
 => [stevedore golang 7/7] RUN go build -ldflags "-s -w" -v -o /usr/local/bin/stevedore ./cmd/stevedore.go                                                           5.1s
 => [stevedore stage-1 2/4] COPY --from=golang /usr/local/bin/stevedore /usr/local/bin/stevedore                                                                     0.0s
 => [stevedore stage-1 3/4] COPY examples/01-basic-example/stack/stevedore/entrypoint.sh /usr/local/bin/entrypoint.sh                                                0.0s
 => [stevedore stage-1 4/4] COPY examples/01-basic-example/stack/stevedore/wait-for-dockerd.sh /usr/local/bin/wait-for-dockerd.sh                                    0.0s
 => [stevedore] exporting to docker image format                                                                                                                     1.3s
 => => exporting layers                                                                                                                                              0.6s
 => => exporting manifest sha256:7fab208b7054a73397ae955984a4a0e61c383a98d01bb6ea443e6eae542c95a2                                                                    0.0s
 => => exporting config sha256:49abb63bc160befb384bde5dbdbb882bf5f849b3cff5ae4daa65a4e70adbc8db                                                                      0.0s
 => => sending tarball                                                                                                                                               0.6s
 => [stevedore stevedore] importing to docker                                                                                                                        0.2s
[+] Running 4/4
 ✔ Network 11-promote-images-example_default         Created                                                                                                         0.1s
 ✔ Container 11-promote-images-example-stevedore-1   Started                                                                                                         0.6s
 ✔ Container 11-promote-images-example-dockerauth-1  Started                                                                                                         0.5s
 ✔ Container 11-promote-images-example-registry-1    Started                                                                                                         0.8s
```

### Waiting for Dockerd to be Ready
Before starting the execution of the Stevedore command, it is important to ensure that the Docker daemon (dockerd) is ready. The stevedore service Docker image includes a script, [wait-for-dockerd.sh](./stack/stevedore/wait-for-dockerd.sh), which can be used to ensure the readiness of the Docker daemon.
```sh
 Run example 11-promote-images-example

 [11-promote-images-example] Waiting for dockerd
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
```

### Building the base image
In this example, the command `stevedore build base` creates the `base` Docker image. As a result, the locally available images are `registry.stevedore.test/base:busybox1.35` and `registry.stevedore.test/base:2.4.6`.
```sh
 [11-promote-images-example] Building base image
registry.stevedore.test/base:2.4.6 Step 1/7 : ARG image_from_name
registry.stevedore.test/base:2.4.6 Step 2/7 : ARG image_from_tag
registry.stevedore.test/base:2.4.6 Step 3/7 : FROM ${image_from_name}:${image_from_tag}
registry.stevedore.test/base:2.4.6 ‣  1.35:  Pulling from library/busybox
registry.stevedore.test/base:2.4.6 ‣  c15cbdab5f8e:  Pull complete
registry.stevedore.test/base:2.4.6 ‣  Digest: sha256:b4e4a06de46acc0958cd93e2eeb769077d255f06a7c3a91196509c16b7bc989e
registry.stevedore.test/base:2.4.6 ‣  Status: Downloaded newer image for busybox:1.35
registry.stevedore.test/base:2.4.6 ---> dddc7578369a
registry.stevedore.test/base:2.4.6 Step 4/7 : RUN echo "anonymous:x:10001:10001:,,,:/app:/bin/sh" >> /etc/passwd &&     echo "anonymous:x:10001:" >> /etc/group &&     mkdir -p /app &&     chown 10001:10001 /app
registry.stevedore.test/base:2.4.6 ---> Running in 612b6713c9c9
registry.stevedore.test/base:2.4.6 ---> d5469204929c
registry.stevedore.test/base:2.4.6 Step 5/7 : USER anonymous
registry.stevedore.test/base:2.4.6 ---> Running in c6d96a50b862
registry.stevedore.test/base:2.4.6 ---> 2b5da65619c0
registry.stevedore.test/base:2.4.6 Step 6/7 : WORKDIR /app
registry.stevedore.test/base:2.4.6 ---> Running in 016e72905614
registry.stevedore.test/base:2.4.6 ---> 2bbdaebff40c
registry.stevedore.test/base:2.4.6 Step 7/7 : LABEL created_at=2023-06-07T05:54:35.257558888Z
registry.stevedore.test/base:2.4.6 ---> Running in f36bc8cbaaaa
registry.stevedore.test/base:2.4.6 ---> d088c654453a
registry.stevedore.test/base:2.4.6  ‣ sha256:d088c654453aa784c099feda862e05e2d1f4a14e12b5efb9bab6e5a6a75681a9
registry.stevedore.test/base:2.4.6 Successfully built d088c654453a
registry.stevedore.test/base:2.4.6 Successfully tagged registry.stevedore.test/base:busybox1.35
registry.stevedore.test/base:2.4.6 Successfully tagged registry.stevedore.test/base:2.4.6
```

### Promoting the base image
After the images are ready and stored locally, the command `stevedore promote registry.stevedore.test/base:2.4.6 --promote-image-registry-namespace stable` is used to promote the image from the local environment to the Docker registry.
It is important to note the use of the `--promote-image-registry-namespace stable` flag, which allows updating the registry namespace of the Docker image to `stable` before pushing it to the Docker registry.
```sh
 [11-promote-images-example] Promoting building image
registry.stevedore.test/stable/base:2.4.6 ‣  The push refers to repository [registry.stevedore.test/stable/base]
registry.stevedore.test/stable/base:2.4.6 ‣  1b3060abf396:  Pushed
registry.stevedore.test/stable/base:2.4.6 ‣  42ef21f45b9a:  Pushed
registry.stevedore.test/stable/base:2.4.6 ‣  2.4.6: digest: sha256:b074b006d5183fd2075a1a5141f2e2d3d47559b78272512c732ddbe7256ef056 size: 735
```
That case can be a starting point of how to manage the lifecycle of the Docker images.

### Building the app1 image
When building the `app1` Docker image, you have the flexibility to set the desired registry namespace. Even though the `app1` image definition itself does not specify a specific namespace, you can use the `--image-from-namespace` flag to achieve this customization.
By executing the command `stevedore build app1 --image-version 0.1.0 --image-from-namespace stable --pull-parent-image`, you can utilize the previously promoted image `registry.stevedore.test/stable/base:2.4.6` as the base image for `app1`.
```sh
 [11-promote-images-example] Building app1 image
registry.stevedore.test/app1:0.1.0-base2.4.6 Step 1/8 : ARG image_from_name
registry.stevedore.test/app1:0.1.0-base2.4.6 Step 2/8 : ARG image_from_registry_host
registry.stevedore.test/app1:0.1.0-base2.4.6 Step 3/8 : ARG image_from_registry_namespace
registry.stevedore.test/app1:0.1.0-base2.4.6 Step 4/8 : ARG image_from_tag
registry.stevedore.test/app1:0.1.0-base2.4.6 Step 5/8 : FROM ${image_from_registry_host}/${image_from_registry_namespace}/${image_from_name}:${image_from_tag}
registry.stevedore.test/app1:0.1.0-base2.4.6 ---> d088c654453a
registry.stevedore.test/app1:0.1.0-base2.4.6 Step 6/8 : COPY ./${app_name}/app.sh /app.sh
registry.stevedore.test/app1:0.1.0-base2.4.6 ---> 9a3349846df5
registry.stevedore.test/app1:0.1.0-base2.4.6 Step 7/8 : CMD ["/app.sh"]
registry.stevedore.test/app1:0.1.0-base2.4.6 ---> Running in dcf34dc6ab48
registry.stevedore.test/app1:0.1.0-base2.4.6 ---> eea871458f7c
registry.stevedore.test/app1:0.1.0-base2.4.6 Step 8/8 : LABEL created_at=2023-06-07T05:54:40.872038442Z
registry.stevedore.test/app1:0.1.0-base2.4.6 ---> Running in 2e26421e6ce0
registry.stevedore.test/app1:0.1.0-base2.4.6 ---> fdbc238604d2
registry.stevedore.test/app1:0.1.0-base2.4.6  ‣ sha256:fdbc238604d27d33420392d4fd824745de24c313cc1695d0d4f6699121eb620a
registry.stevedore.test/app1:0.1.0-base2.4.6 Successfully built fdbc238604d2
registry.stevedore.test/app1:0.1.0-base2.4.6 Successfully tagged registry.stevedore.test/app1:0.1.0-base2.4.6
```
Take note of the usage of the `--pull-parent-image` flag, which ensures that the parent image is pulled from the registry before the build process starts.

The resulting image, `registry.stevedore.test/app1:0.1.0-base2.4.6`, is now stored locally and pending promotion.

### Promoting the app1 image
Finally, to promote the image `registry.stevedore.test/app1:0.1.0-base2.4.6`, execute the command `stevedore promote registry.stevedore.test/app1:0.1.0-base2.4.6 --promote-image-tag latest --force-promote-source-image --enable-semver-tags`. This promotion involves several actions:
- The `--promote-image-tag` flag overwrites the source image tag with the specified tag `latest`.
- The `--force-promote-source-image` flag ensures that the tag of the source image is pushed.
- The `--enable-semver-tags` flag enables automatic tag generation based on the semantic version.

The output below shows the result of the promotion:
```sh
 [11-promote-images-example] Promoting app1 image
registry.stevedore.test/app1:latest ‣  The push refers to repository [registry.stevedore.test/app1]
registry.stevedore.test/app1:latest ‣  0d1764f937d0:  Pushed
registry.stevedore.test/app1:latest ‣  1b3060abf396:  Mounted from stable/base
registry.stevedore.test/app1:latest ‣  42ef21f45b9a:  Mounted from stable/base
registry.stevedore.test/app1:latest ‣  0.1: digest: sha256:4c935e2e4eb9f50a622017e38c6331483afbd3cb5c96dd8cbd694c4fe93d93c9 size: 942
registry.stevedore.test/app1:latest ‣  The push refers to repository [registry.stevedore.test/app1]
registry.stevedore.test/app1:latest ‣  0d1764f937d0:  Layer already exists
registry.stevedore.test/app1:latest ‣  1b3060abf396:  Layer already exists
registry.stevedore.test/app1:latest ‣  42ef21f45b9a:  Layer already exists
registry.stevedore.test/app1:latest ‣  0.1.0: digest: sha256:4c935e2e4eb9f50a622017e38c6331483afbd3cb5c96dd8cbd694c4fe93d93c9 size: 942
registry.stevedore.test/app1:latest ‣  The push refers to repository [registry.stevedore.test/app1]
registry.stevedore.test/app1:latest ‣  0d1764f937d0:  Layer already exists
registry.stevedore.test/app1:latest ‣  1b3060abf396:  Layer already exists
registry.stevedore.test/app1:latest ‣  42ef21f45b9a:  Layer already exists
registry.stevedore.test/app1:latest ‣  0.1.0-base2.4.6: digest: sha256:4c935e2e4eb9f50a622017e38c6331483afbd3cb5c96dd8cbd694c4fe93d93c9 size: 942
registry.stevedore.test/app1:latest ‣  The push refers to repository [registry.stevedore.test/app1]
registry.stevedore.test/app1:latest ‣  0d1764f937d0:  Layer already exists
registry.stevedore.test/app1:latest ‣  1b3060abf396:  Layer already exists
registry.stevedore.test/app1:latest ‣  42ef21f45b9a:  Layer already exists
registry.stevedore.test/app1:latest ‣  latest: digest: sha256:4c935e2e4eb9f50a622017e38c6331483afbd3cb5c96dd8cbd694c4fe93d93c9 size: 942
```
As a result, the following tags are automatically created and pushed simultaneously:  `registry.stevedore.test/app1:0.1`, `registry.stevedore.test/app1:0.1.0`, `registry.stevedore.test/app1:0.1.0-base2.4.6`, `registry.stevedore.test/app1:latest`.

### Cleaning the stack
```sh
Stopping the stack to run 11-promote-images-example

[+] Running 4/4
 ✔ Container 11-promote-images-example-stevedore-1   Removed                                                                                                         3.3s
 ✔ Container 11-promote-images-example-registry-1    Removed                                                                                                         0.2s
 ✔ Container 11-promote-images-example-dockerauth-1  Removed                                                                                                         0.3s
 ✔ Network 11-promote-images-example_default         Removed                                                                                                         0.5s
```
