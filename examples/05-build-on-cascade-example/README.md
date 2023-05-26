# Build On Cascade Example

This example demonstrates the concept of building images on a cascade in Stevedore. It showcases how to create a foundational image with a common configuration and then build multiple applications using that image as a base, streamlining the image building process and ensuring consistency across the applications.

- [Build On Cascade Example](#build-on-cascade-example)
  - [Requirements](#requirements)
  - [Stack](#stack)
  - [Usage](#usage)
  - [Expected Output](#expected-output)
    - [Getting images](#getting-images)
    - [Building images](#building-images)
    - [Cleaning the stack](#cleaning-the-stack)
  - [Additional information](#additional-information)
    - [Images](#images)


## Requirements
- Docker. _Tested on Docker server 20.10.21 and Docker API 1.41_
- Docker's Compose plugin or `docker-compose`. _Tested on Docker Compose version v2.17.3_
- `make` utility. _Tested on version 4.3-4.1build1_

## Stack
The stack required to run this example is defined in a [Docker Compose file](./docker-compose.yml). The stack consists of three services: a Docker Registry, a Docker Registry authorization and a Stevedore service. The Docker registry is used to store the Docker images built by Stevedore during the example execution. The Stevedore service is where the example is executed.

The Stevedore service is built from a container which is defined in the [Dockerfile](https://github.com/gostevedore/stevedore/blob/main/test/stack/client/Dockerfile) present in the `test/stack/client` directory.

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
Below is the expected output for the `make run` command, which starts the Docker stack, gets some information about the Stevedore configuration, builds and promotes a Docker images using Stevedore, and then cleans the stack up.

```sh
Starting the stack to run 05-build-on-cascade-example

[+] Building 33.5s (22/22) FINISHED
 => [internal] booting buildkit                                                                                                                                     11.5s
 => => pulling image moby/buildkit:buildx-stable-1                                                                                                                  10.8s
 => => creating container buildx_buildkit_buildkit0                                                                                                                  0.7s
 => [internal] load .dockerignore                                                                                                                                    0.0s
 => => transferring context: 2B                                                                                                                                      0.0s
 => [internal] load build definition from Dockerfile                                                                                                                 0.0s
 => => transferring dockerfile: 989B                                                                                                                                 0.0s
 => [internal] load metadata for docker.io/library/docker:20.10-dind                                                                                                 6.8s
 => [internal] load metadata for docker.io/library/golang:1.19-alpine                                                                                               13.4s
 => [internal] load build context                                                                                                                                    0.1s
 => => transferring context: 7.16MB                                                                                                                                  0.1s
 => [stage-1 1/6] FROM docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                                   0.0s
 => => resolve docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                                           0.0s
 => [golang 1/8] FROM docker.io/library/golang:1.19-alpine@sha256:4147c2ad24a347e53a8a361b663dcf82fca4157a3f9a5136d696a4e53dd22b65                                   0.0s
 => => resolve docker.io/library/golang:1.19-alpine@sha256:4147c2ad24a347e53a8a361b663dcf82fca4157a3f9a5136d696a4e53dd22b65                                          0.0s
 => CACHED [golang 2/8] WORKDIR /usr/src/app                                                                                                                         0.0s
 => CACHED [golang 3/8] RUN apk add --no-cache make build-base                                                                                                       0.0s
 => CACHED [golang 4/8] COPY go.mod ./                                                                                                                               0.0s
 => CACHED [golang 5/8] COPY go.sum ./                                                                                                                               0.0s
 => CACHED [golang 6/8] RUN go mod download && go mod verify                                                                                                         0.0s
 => [golang 7/8] COPY . ./                                                                                                                                           0.7s
 => [golang 8/8] RUN go build -ldflags "-s -w" -v -o /usr/local/bin/stevedore ./cmd/stevedore.go                                                                     5.1s
 => CACHED [stage-1 2/6] COPY --from=golang /usr/local/go /usr/local/go                                                                                              0.0s
 => CACHED [stage-1 3/6] COPY --from=golang /usr/local/bin/stevedore /usr/local/bin/stevedore                                                                        0.0s
 => CACHED [stage-1 4/6] WORKDIR /go                                                                                                                                 0.0s
 => CACHED [stage-1 5/6] RUN mkdir -p "/go/src" "/go/bin" && chmod -R 777 "/go"                                                                                      0.0s
 => CACHED [stage-1 6/6] COPY test/stack/client/entrypoint.sh /usr/local/bin/entrypoint.sh                                                                           0.0s
 => exporting to docker image format                                                                                                                                 1.0s
 => => exporting layers                                                                                                                                              0.0s
 => => exporting manifest sha256:8c9795111721911e24ac0c2652bda6b35647a39af933be2b486f1cc3335ce31b                                                                    0.0s
 => => exporting config sha256:998dfdd399ab14b8a4620fffea07c0a6caed4f8a2752d34a7fa68c0b74a9e3b9                                                                      0.0s
 => => sending tarball                                                                                                                                               1.0s
 => importing to docker                                                                                                                                              0.0s
[+] Running 4/4
 ✔ Network 05-build-on-cascade-example_default         Created                                                                                                       0.1s
 ✔ Container 05-build-on-cascade-example-stevedore-1   Started                                                                                                       0.6s
 ✔ Container 05-build-on-cascade-example-dockerauth-1  Started                                                                                                       0.4s
 ✔ Container 05-build-on-cascade-example-registry-1    Started                                                                                                       0.8s
```

### Getting images
To view the images in a tree format, execute the command `stevedore get images --tree`. This command will display a hierarchical representation of the images. You will observe the `base` images nested under the `busybox` images, along with the `app1`, `app2`, and `app3` images defined from the `base`.

```sh
 [05-build-on-cascade-example] Get images
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
The example utilizes the command `stevedore build base --build-on-cascade --push-after-build` to build the `base` images. This command triggers the automatic building of their descendants once the `base` images are successfully built.
The promotion of the images to the Docker registry is initiated automatically once each image is ready. Since all three images are being built concurrently, the output may display a mixture of these outputs.

```sh
 [05-build-on-cascade-example] Build the base image and its descendants, and push the images after build
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
/certs/server/cert.pem: OK
/certs/client/cert.pem: OK
 Waiting for dockerd to be ready...
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
registry.stevedore.test/base:busybox-1.35 ---> Running in 20d146cbd446
registry.stevedore.test/base:busybox-1.35 ---> 050b72f9afb3
registry.stevedore.test/base:busybox-1.35 Step 5/7 : USER anonymous
registry.stevedore.test/base:busybox-1.35 ---> Running in e06c1006dba5
registry.stevedore.test/base:busybox-1.35 ---> 5c7a9fa446b9
registry.stevedore.test/base:busybox-1.35 Step 6/7 : WORKDIR /app
registry.stevedore.test/base:busybox-1.35 ---> Running in 8c0fc43db1a6
registry.stevedore.test/base:busybox-1.35 ---> 4633390e596f
registry.stevedore.test/base:busybox-1.35 Step 7/7 : LABEL created_at=2023-05-22T05:33:55.180068046Z
registry.stevedore.test/base:busybox-1.35 ---> Running in 388dc43bf6c0
registry.stevedore.test/base:busybox-1.35 ---> 8f191579b101
registry.stevedore.test/base:busybox-1.35  ‣ sha256:8f191579b101efd918b5c115cf0bfd79919aaf543594cb8d930070d3ec7fa1f4
registry.stevedore.test/base:busybox-1.35 Successfully built 8f191579b101
registry.stevedore.test/base:busybox-1.35 Successfully tagged registry.stevedore.test/base:busybox-1.35
registry.stevedore.test/base:busybox-1.35 ‣  The push refers to repository [registry.stevedore.test/base]
registry.stevedore.test/base:busybox-1.35 ‣  d7c5e4244af0:  Pushed
registry.stevedore.test/base:busybox-1.35 ‣  42ef21f45b9a:  Pushed
registry.stevedore.test/base:busybox-1.35 ‣  busybox-1.35: digest: sha256:7ae2498b0fcaa80517a1b18ba616b72631a86c827b9de1af6bbc95f2d49191b0 size: 735
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 1/7 : ARG image_from_name
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 2/7 : ARG image_from_tag
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 3/7 : ARG image_from_registry_host
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 4/7 : FROM ${image_from_registry_host}/${image_from_name}:${image_from_tag}
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 1/7 : ARG image_from_name
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 2/7 : ARG image_from_tag
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 3/7 : ARG image_from_registry_host
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 4/7 : FROM ${image_from_registry_host}/${image_from_name}:${image_from_tag}
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> 8f191579b101
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 5/7 : COPY ./${app_name}/app.sh /app/run.sh
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> 8f191579b101
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 5/7 : COPY ./${app_name}/app.sh /app/run.sh
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> 3c858048ab5b
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 6/7 : CMD ["/app/run.sh"]
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> a63e51a95a1c
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 6/7 : CMD ["/app/run.sh"]
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> Running in f9f0e182c6ee
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> Running in ddfa5811767d
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> 41ad30af0ef5
registry.stevedore.test/app2:v1-base-busybox-1.35 Step 7/7 : LABEL created_at=2023-05-22T05:33:55.180068046Z
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> d49923280689
registry.stevedore.test/app1:v1-base-busybox-1.35 Step 7/7 : LABEL created_at=2023-05-22T05:33:55.180068046Z
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> Running in ac7b9ce8b009
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> Running in 6cfa1a5e9be2
registry.stevedore.test/app2:v1-base-busybox-1.35 ---> a768ff607b09
registry.stevedore.test/app2:v1-base-busybox-1.35  ‣ sha256:a768ff607b09eb3666e30142082b95f0476766321e1ab3cc0c2fd3914b97c09e
registry.stevedore.test/app2:v1-base-busybox-1.35 Successfully built a768ff607b09
registry.stevedore.test/app1:v1-base-busybox-1.35 ---> 40c05a077121
registry.stevedore.test/app1:v1-base-busybox-1.35  ‣ sha256:40c05a077121a69767e29130b433ae853d4fd9900a56e0e27c880b076827b1e1
registry.stevedore.test/app1:v1-base-busybox-1.35 Successfully built 40c05a077121
registry.stevedore.test/app2:v1-base-busybox-1.35 Successfully tagged registry.stevedore.test/app2:v1-base-busybox-1.35
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  The push refers to repository [registry.stevedore.test/app2]
registry.stevedore.test/app1:v1-base-busybox-1.35 Successfully tagged registry.stevedore.test/app1:v1-base-busybox-1.35
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  The push refers to repository [registry.stevedore.test/app1]
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  20684b00e2f9:  Preparing
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  d7c5e4244af0:  Preparing
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  20684b00e2f9:  Pushed
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  d7c5e4244af0:  Mounted from base
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  42ef21f45b9a:  Mounted from base
registry.stevedore.test/app2:v1-base-busybox-1.35 ‣  v1-base-busybox-1.35: digest: sha256:40d68e3eade0f2e12d8dde52bd15ebf5ba1922e8af6edc172e929c9ee4b47e86 size: 942
registry.stevedore.test/app1:v1-base-busybox-1.35 ‣  v1-base-busybox-1.35: digest: sha256:811157e3070b448d0c62a2d81a25d814181d4caa25b75df02239395a4bbcda78 size: 942
registry.stevedore.test/base:busybox-1.36 ‣  1.36:  Pulling from library/busybox
registry.stevedore.test/base:busybox-1.36 ‣  325d69979d33:  Pull complete
registry.stevedore.test/base:busybox-1.36 ‣  Digest: sha256:560af6915bfc8d7630e50e212e08242d37b63bd5c1ccf9bd4acccf116e262d5b
registry.stevedore.test/base:busybox-1.36 ‣  Status: Downloaded newer image for busybox:1.36
registry.stevedore.test/base:busybox-1.36 ---> 8135583d97fe
registry.stevedore.test/base:busybox-1.36 Step 4/7 : RUN echo "anonymous:x:10001:10001:,,,:/app:/bin/sh" >> /etc/passwd &&     echo "anonymous:x:10001:" >> /etc/group &&     mkdir -p /app &&     chown 10001:10001 /app
registry.stevedore.test/base:busybox-1.36 ---> Running in 04f6bbb1fca9
registry.stevedore.test/base:busybox-1.36 ---> 2332038a725a
registry.stevedore.test/base:busybox-1.36 Step 5/7 : USER anonymous
registry.stevedore.test/base:busybox-1.36 ---> Running in 42a1771efccd
registry.stevedore.test/base:busybox-1.36 ---> 5c8a80488120
registry.stevedore.test/base:busybox-1.36 Step 6/7 : WORKDIR /app
registry.stevedore.test/base:busybox-1.36 ---> Running in e864a04a4dbe
registry.stevedore.test/base:busybox-1.36 ---> a6d66cef4388
registry.stevedore.test/base:busybox-1.36 Step 7/7 : LABEL created_at=2023-05-22T05:33:55.179477523Z
registry.stevedore.test/base:busybox-1.36 ---> Running in 3af519ba934a
registry.stevedore.test/base:busybox-1.36 ---> 3f1a26b4b251
registry.stevedore.test/base:busybox-1.36  ‣ sha256:3f1a26b4b251d1539ecc05100158190d4a097177374995ebab8911d30735d473
registry.stevedore.test/base:busybox-1.36 Successfully built 3f1a26b4b251
registry.stevedore.test/base:busybox-1.36 Successfully tagged registry.stevedore.test/base:busybox-1.36
registry.stevedore.test/base:busybox-1.36 ‣  The push refers to repository [registry.stevedore.test/base]
registry.stevedore.test/base:busybox-1.36 ‣  2054dc80e2f0:  Pushed
registry.stevedore.test/base:busybox-1.36 ‣  9547b4c33213:  Pushed
registry.stevedore.test/base:busybox-1.36 ‣  busybox-1.36: digest: sha256:5f064c61ca61983506c5b6dea7b2f3d627c7be5064c5f8f66c733b38a05e232c size: 735
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 1/7 : ARG image_from_name
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 2/7 : ARG image_from_tag
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 3/7 : ARG image_from_registry_host
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 4/7 : FROM ${image_from_registry_host}/${image_from_name}:${image_from_tag}
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> 3f1a26b4b251
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 5/7 : COPY ./${app_name}/app.sh /app/run.sh
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 1/7 : ARG image_from_name
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 2/7 : ARG image_from_tag
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 3/7 : ARG image_from_registry_host
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 4/7 : FROM ${image_from_registry_host}/${image_from_name}:${image_from_tag}
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> 3f1a26b4b251
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 5/7 : COPY ./${app_name}/app.sh /app/run.sh
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> a6426501a096
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 6/7 : CMD ["/app/run.sh"]
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> 8ac13fcb5de6
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 6/7 : CMD ["/app/run.sh"]
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> Running in 20c366c42b67
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> Running in 425223be931d
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> a0f7d0be7419
registry.stevedore.test/app2:v1-base-busybox-1.36 Step 7/7 : LABEL created_at=2023-05-22T05:33:55.179477523Z
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> 5f49976bad79
registry.stevedore.test/app3:v1-base-busybox-1.36 Step 7/7 : LABEL created_at=2023-05-22T05:33:55.179477523Z
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> Running in ded577a8b1cf
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> Running in 0d687765f5d1
registry.stevedore.test/app3:v1-base-busybox-1.36 ---> 698a68e508b2
registry.stevedore.test/app3:v1-base-busybox-1.36  ‣ sha256:698a68e508b2b10064ffa1a412e2cdc448d04c8d8de3c21860fca77d2531bf51
registry.stevedore.test/app3:v1-base-busybox-1.36 Successfully built 698a68e508b2
registry.stevedore.test/app2:v1-base-busybox-1.36 ---> 32cc202ba0a7
registry.stevedore.test/app2:v1-base-busybox-1.36  ‣ sha256:32cc202ba0a70c7f7d1517beba89e957962f064fd2a2e50add0d4998f34db717
registry.stevedore.test/app2:v1-base-busybox-1.36 Successfully built 32cc202ba0a7
registry.stevedore.test/app3:v1-base-busybox-1.36 Successfully tagged registry.stevedore.test/app3:v1-base-busybox-1.36
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  The push refers to repository [registry.stevedore.test/app3]
registry.stevedore.test/app2:v1-base-busybox-1.36 Successfully tagged registry.stevedore.test/app2:v1-base-busybox-1.36
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  The push refers to repository [registry.stevedore.test/app2]
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  8f0ff8370a82:  Preparing
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  2054dc80e2f0:  Preparing
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  7eca5c1dd819:  Pushed
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  2054dc80e2f0:  Mounted from base
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  9547b4c33213:  Mounted from base
registry.stevedore.test/app3:v1-base-busybox-1.36 ‣  v1-base-busybox-1.36: digest: sha256:adbe2c9d5bfed7cbab58e9c298f0999b3fbdac1a10380a4902cd3c98e9b9fba6 size: 942
registry.stevedore.test/app2:v1-base-busybox-1.36 ‣  v1-base-busybox-1.36: digest: sha256:cd6a0e2ec02c0a16fe7b99d8dcc3711afd57b0e64b1776a7c5480b493c9e74e7 size: 942
```

### Cleaning the stack
```sh
Stopping the stack to run 05-build-on-cascade-example

[+] Running 6/6
 ✔ Container 05-build-on-cascade-example-registry-1                  Removed                                                                                         0.4s
 ✔ Container 05-build-on-cascade-example-stevedore-1                 Removed                                                                                        10.4s
 ✔ Container 05-build-on-cascade-example-stevedore-run-ee2456c22aa0  Removed                                                                                         0.0s
 ✔ Container 05-build-on-cascade-example-stevedore-run-4f10db448376  Removed                                                                                         0.0s
 ✔ Container 05-build-on-cascade-example-dockerauth-1                Removed                                                                                         0.3s
 ✔ Network 05-build-on-cascade-example_default                       Removed                                                                                         0.4s
```

## Additional information
In addition to the core steps outlined in the example, the following section provides additional information and insights to further enhance your understanding of how this example uses Stevedore.

### Images
This example showcases the process of defining a base image that establishes a shared configuration for your images, including the user used to run the containers. As a result, all images derived from this base will inherit and utilize the common configuration.

The following configuration is available by default for all the images created from the `base` image.
```dockerfile
# Create a new user
RUN echo "anonymous:x:10001:10001:,,,:/app:/bin/sh" >> /etc/passwd && \
    echo "anonymous:x:10001:" >> /etc/group && \
    mkdir -p /app && \
    chown 10001:10001 /app

# Set the user as the default user
USER anonymous
WORKDIR /app
```

Additionally, the example highlights the use of the `build-on-cascade` plan, allowing you to rebuild all descendant images of the base whenever changes need to be applied to the foundational image.
