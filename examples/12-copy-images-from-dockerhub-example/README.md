# Copying Images From DockerHub Example

This example illustrates how to copy Docker images from Docker Hub to your local Docker registry using Stevedore, allowing you to have a local copy of the desired images for offline or restricted environments.

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
Starting the stack to run 12-copy-images-from-dockerhub-example

[+] Building 7.5s (18/18) FINISHED
 => [stevedore internal] load .dockerignore                                                                                                                          0.0s
 => => transferring context: 2B                                                                                                                                      0.0s
 => [stevedore internal] load build definition from Dockerfile                                                                                                       0.0s
 => => transferring dockerfile: 888B                                                                                                                                 0.0s
 => [stevedore internal] load metadata for docker.io/library/docker:20.10-dind                                                                                       1.3s
 => [stevedore internal] load metadata for docker.io/library/golang:1.19-alpine                                                                                      1.8s
 => [stevedore golang 1/7] FROM docker.io/library/golang:1.19-alpine@sha256:106e12a19d4f8360c7dea4c08ef6ab7b62ec153972c77e099998eff9cb87e4f0                         0.0s
 => => resolve docker.io/library/golang:1.19-alpine@sha256:106e12a19d4f8360c7dea4c08ef6ab7b62ec153972c77e099998eff9cb87e4f0                                          0.0s
 => [stevedore internal] load build context                                                                                                                          0.1s
 => => transferring context: 346.26kB                                                                                                                                0.0s
 => [stevedore stage-1 1/4] FROM docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                         0.0s
 => => resolve docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                                           0.0s
 => CACHED [stevedore golang 2/7] WORKDIR /usr/src/app                                                                                                               0.0s
 => CACHED [stevedore golang 3/7] RUN apk add --no-cache make build-base                                                                                             0.0s
 => CACHED [stevedore golang 4/7] COPY go.mod go.sum ./                                                                                                              0.0s
 => CACHED [stevedore golang 5/7] RUN go mod download                                                                                                                0.0s
 => [stevedore golang 6/7] COPY . ./                                                                                                                                 0.5s
 => [stevedore golang 7/7] RUN go build -ldflags "-s -w" -v -o /usr/local/bin/stevedore ./cmd/stevedore.go                                                           4.6s
 => CACHED [stevedore stage-1 2/4] COPY --from=golang /usr/local/bin/stevedore /usr/local/bin/stevedore                                                              0.0s
 => CACHED [stevedore stage-1 3/4] COPY examples/01-basic-example/stack/stevedore/entrypoint.sh /usr/local/bin/entrypoint.sh                                         0.0s
 => CACHED [stevedore stage-1 4/4] COPY examples/01-basic-example/stack/stevedore/wait-for-dockerd.sh /usr/local/bin/wait-for-dockerd.sh                             0.0s
 => [stevedore] exporting to docker image format                                                                                                                     0.4s
 => => exporting layers                                                                                                                                              0.0s
 => => exporting manifest sha256:07ea2209dd5b11ad4035e9bfc1b2636ec28bb81b1feec1e99ec14c7d718880a3                                                                    0.0s
 => => exporting config sha256:dae71832d549322f79b10ecfce3316c2fe24f138db41e5921f7921c20b628e13                                                                      0.0s
 => => sending tarball                                                                                                                                               0.4s
 => [stevedore stevedore] importing to docker                                                                                                                        0.0s
[+] Running 4/4
 ✔ Network 12-copy-images-from-dockerhub-example_default         Created                                                                                             0.1s
 ✔ Container 12-copy-images-from-dockerhub-example-dockerauth-1  Started                                                                                             0.6s
 ✔ Container 12-copy-images-from-dockerhub-example-stevedore-1   Started                                                                                             0.6s
 ✔ Container 12-copy-images-from-dockerhub-example-registry-1    Started                                                                                             0.8s
```

### Waiting for Dockerd to be Ready
Before starting the execution of the Stevedore command, it is important to ensure that the Docker daemon (dockerd) is ready. The stevedore service Docker image includes a script, [wait-for-dockerd.sh](./stack/stevedore/wait-for-dockerd.sh), which can be used to ensure the readiness of the Docker daemon.
```sh
 Run example 12-copy-images-from-dockerhub-example

 [12-copy-images-from-dockerhub-example] Waiting for dockerd
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
```

### Copying the Busybox Image from DockerHub to a Private Docker Registry
To start, the example demonstrates the process of obtaining the busybox:1.35 image from Docker Hub and copying it to the private registry running in the example's stack. This is achieved by executing the following command:
`stevedore copy busybox:1.35 --promote-image-registry-host registry.stevedore.test --use-source-image-from-remote --remove-local-images-after-push`

To ensure that a locally stored image is not pushed to the registry, is set the `--use-source-image-from-remote` flag, which enforces the use of the image from a remote source.

Note that the `copy` subcommand is used, which is an alias for the `promote` subcommand in the [Stevedore CLI](https://gostevedore.github.io/docs/reference-guide/cli/#promote). 

```sh
registry.stevedore.test/library/busybox:1.35 ‣  1.35:  Pulling from library/busybox
registry.stevedore.test/library/busybox:1.35 ‣  c15cbdab5f8e:  Pull complete
registry.stevedore.test/library/busybox:1.35 ‣  Digest: sha256:b4e4a06de46acc0958cd93e2eeb769077d255f06a7c3a91196509c16b7bc989e
registry.stevedore.test/library/busybox:1.35 ‣  Status: Downloaded newer image for busybox:1.35
registry.stevedore.test/library/busybox:1.35 ‣  The push refers to repository [registry.stevedore.test/library/busybox]
registry.stevedore.test/library/busybox:1.35 ‣  42ef21f45b9a:  Pushed
registry.stevedore.test/library/busybox:1.35 ‣  1.35: digest: sha256:2197ffa9bd16c893488bc26712a9dd28826daf2abb1a1dabf554fe32615a541d size: 528
registry.stevedore.test/library/busybox:1.35 untagged:  registry.stevedore.test/library/busybox:1.35
registry.stevedore.test/library/busybox:1.35 untagged:  registry.stevedore.test/library/busybox@sha256:2197ffa9bd16c893488bc26712a9dd28826daf2abb1a1dabf554fe32615a541d
```

It is worth noting that when referring to `busybox:1.35`, the image is stored in the hidden `library` Docker Hub namespace, resulting in the image being tagged as `registry.stevedore.test/library/busybox:1.35`.

### Building the app1 image
Once the copied image is available in the private Docker registry, the next step is to build the `app1` Docker image. You can initiate the build process using the following command:
`stevedore build app1 --push-after-build --remove-local-images-after-push`

This command not only builds the image but also automatically pushes it to the Docker registry and removes the local images after the push operation.
```sh
 [12-copy-images-from-dockerhub-example] Build my-app and push images after build
registry.stevedore.test/app1:v1-busybox1.35 Step 1/8 : ARG image_from_name
registry.stevedore.test/app1:v1-busybox1.35 Step 2/8 : ARG image_from_registry_host
registry.stevedore.test/app1:v1-busybox1.35 Step 3/8 : ARG image_from_registry_namespace
registry.stevedore.test/app1:v1-busybox1.35 Step 4/8 : ARG image_from_tag
registry.stevedore.test/app1:v1-busybox1.35 Step 5/8 : FROM ${image_from_registry_host}/${image_from_registry_namespace}/${image_from_name}:${image_from_tag}
registry.stevedore.test/app1:v1-busybox1.35 ‣  1.35:  Pulling from library/busybox
registry.stevedore.test/app1:v1-busybox1.35 ‣  Digest: sha256:2197ffa9bd16c893488bc26712a9dd28826daf2abb1a1dabf554fe32615a541d
registry.stevedore.test/app1:v1-busybox1.35 ‣  Status: Downloaded newer image for registry.stevedore.test/library/busybox:1.35
registry.stevedore.test/app1:v1-busybox1.35 ---> dddc7578369a
registry.stevedore.test/app1:v1-busybox1.35 Step 6/8 : COPY ./${app_name}/app.sh /app.sh
registry.stevedore.test/app1:v1-busybox1.35 ---> 17d7326e6386
registry.stevedore.test/app1:v1-busybox1.35 Step 7/8 : CMD ["/app.sh"]
registry.stevedore.test/app1:v1-busybox1.35 ---> Running in fc50765d3e5e
registry.stevedore.test/app1:v1-busybox1.35 ---> 6465f8bcfe96
registry.stevedore.test/app1:v1-busybox1.35 Step 8/8 : LABEL created_at=2023-06-09T13:56:39.351801142Z
registry.stevedore.test/app1:v1-busybox1.35 ---> Running in 34f16c56e0cb
registry.stevedore.test/app1:v1-busybox1.35 ---> c4a5075087f1
registry.stevedore.test/app1:v1-busybox1.35  ‣ sha256:c4a5075087f164133621a8e3230a9b544d60f4619aa4e0474fc9899fc39ddd8b
registry.stevedore.test/app1:v1-busybox1.35 Successfully built c4a5075087f1
registry.stevedore.test/app1:v1-busybox1.35 Successfully tagged registry.stevedore.test/app1:v1-busybox1.35
registry.stevedore.test/app1:v1-busybox1.35 ‣  The push refers to repository [registry.stevedore.test/app1]
registry.stevedore.test/app1:v1-busybox1.35 ‣  019c59cb31af:  Pushed
registry.stevedore.test/app1:v1-busybox1.35 ‣  42ef21f45b9a:  Mounted from library/busybox
registry.stevedore.test/app1:v1-busybox1.35 ‣  v1-busybox1.35: digest: sha256:db93b691c976470053eaf884cc819120974c78d2a29f6070cb3a9a16f3d4a398 size: 735
registry.stevedore.test/app1:v1-busybox1.35 untagged:  registry.stevedore.test/app1:v1-busybox1.35
registry.stevedore.test/app1:v1-busybox1.35 untagged:  registry.stevedore.test/app1@sha256:db93b691c976470053eaf884cc819120974c78d2a29f6070cb3a9a16f3d4a398
registry.stevedore.test/app1:v1-busybox1.35 deleted:  sha256:c4a5075087f164133621a8e3230a9b544d60f4619aa4e0474fc9899fc39ddd8b
```

### Cleaning the stack
```sh
Stopping the stack to run 12-copy-images-from-dockerhub-example

[+] Running 4/4
 ✔ Container 12-copy-images-from-dockerhub-example-registry-1    Removed                                                                                             0.2s
 ✔ Container 12-copy-images-from-dockerhub-example-stevedore-1   Removed                                                                                             3.3s
 ✔ Container 12-copy-images-from-dockerhub-example-dockerauth-1  Removed                                                                                             0.4s
 ✔ Network 12-copy-images-from-dockerhub-example_default         Removed                                                                                             0.4s
```