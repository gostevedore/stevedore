# Build for Multiple Applications Example

This example showcases how to define and manage multiple applications in Stevedore, allowing you to build and push their Docker images.

- [Build for Multiple Applications Example](#build-for-multiple-applications-example)
  - [Requirements](#requirements)
  - [Stack](#stack)
  - [Usage](#usage)
  - [Example Execution Insights](#example-execution-insights)
    - [Getting images](#getting-images)
    - [Building images](#building-images)
      - [app1](#app1)
      - [app2](#app2)
      - [app3](#app3)
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

```sh
Starting the stack to run 04-multiple-apps-example

[+] Building 1.6s (21/21) FINISHED                                                                                                                               
 => [internal] load .dockerignore                                                                                                                           0.0s
 => => transferring context: 2B                                                                                                                             0.0s
 => [internal] load build definition from Dockerfile                                                                                                        0.0s
 => => transferring dockerfile: 989B                                                                                                                        0.0s
 => [internal] load metadata for docker.io/library/docker:20.10-dind                                                                                        0.6s
 => [internal] load metadata for docker.io/library/golang:1.19-alpine                                                                                       0.6s
 => [golang 1/8] FROM docker.io/library/golang:1.19-alpine@sha256:4147c2ad24a347e53a8a361b663dcf82fca4157a3f9a5136d696a4e53dd22b65                          0.0s
 => => resolve docker.io/library/golang:1.19-alpine@sha256:4147c2ad24a347e53a8a361b663dcf82fca4157a3f9a5136d696a4e53dd22b65                                 0.0s
 => [internal] load build context                                                                                                                           0.0s
 => => transferring context: 179.38kB                                                                                                                       0.0s
 => [stage-1 1/6] FROM docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                          0.0s
 => => resolve docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                                  0.0s
 => CACHED [golang 2/8] WORKDIR /usr/src/app                                                                                                                0.0s
 => CACHED [golang 3/8] RUN apk add --no-cache make build-base                                                                                              0.0s
 => CACHED [golang 4/8] COPY go.mod ./                                                                                                                      0.0s
 => CACHED [golang 5/8] COPY go.sum ./                                                                                                                      0.0s
 => CACHED [golang 6/8] RUN go mod download && go mod verify                                                                                                0.0s
 => CACHED [golang 7/8] COPY . ./                                                                                                                           0.0s
 => CACHED [golang 8/8] RUN go build -ldflags "-s -w" -v -o /usr/local/bin/stevedore ./cmd/stevedore.go                                                     0.0s
 => CACHED [stage-1 2/6] COPY --from=golang /usr/local/go /usr/local/go                                                                                     0.0s
 => CACHED [stage-1 3/6] COPY --from=golang /usr/local/bin/stevedore /usr/local/bin/stevedore                                                               0.0s
 => CACHED [stage-1 4/6] WORKDIR /go                                                                                                                        0.0s
 => CACHED [stage-1 5/6] RUN mkdir -p "/go/src" "/go/bin" && chmod -R 777 "/go"                                                                             0.0s
 => CACHED [stage-1 6/6] COPY test/stack/client/entrypoint.sh /usr/local/bin/entrypoint.sh                                                                  0.0s
 => exporting to docker image format                                                                                                                        0.9s
 => => exporting layers                                                                                                                                     0.0s
 => => exporting manifest sha256:bfdf8be5dd1d180066e7d635bcada076ef197f16b20046716747409af60ea56e                                                           0.0s
 => => exporting config sha256:9657515b4d35fcde7bd356f94f6ada2a27904b3c160827d3e394a4c947672dbf                                                             0.0s
 => => sending tarball                                                                                                                                      0.9s
 => importing to docker                                                                                                                                     0.0s
[+] Running 4/4
 ✔ Network 04-multiple-apps-example_default         Created                                                                                                 0.1s 
 ✔ Container 04-multiple-apps-example-stevedore-1   Started                                                                                                 0.7s 
 ✔ Container 04-multiple-apps-example-dockerauth-1  Started                                                                                                 0.7s 
 ✔ Container 04-multiple-apps-example-registry-1    Started                                                                                                 0.9s 
 ```

### Getting images
To view the images in tree format, run `stevedore get images --tree`.

```sh
 [04-multiple-apps-example] Get images
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
/certs/server/cert.pem: OK
 Waiting for dockerd to be ready...
/certs/client/cert.pem: OK
├─── busybox:1.35
│  ├─── registry.stevedore.test/app1:v1-busybox1.35
│  ├─── registry.stevedore.test/app2:v1-busybox1.35
├─── busybox:1.36
│  ├─── registry.stevedore.test/app2:v1-busybox1.36
│  ├─── registry.stevedore.test/app3:v1-busybox1.36
```

### Building images
In this section, you can explore the output of the Docker image builds for the applications `app1`, `app2` and `app3`. The image definitions for these applications reside in the [./images](./images) folder, while their source code can be found in the [./apps](./apps) folder.
Note that each application is built independently of one another.

#### app1
To initiate the build process for `app1`, execute the command `stevedore build app1 --push-after-build`.
```sh
 [04-multiple-apps-example] Build app1 and push images after build
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
/certs/server/cert.pem: OK
 Waiting for dockerd to be ready...
/certs/client/cert.pem: OK
registry.stevedore.test/app1:v1-busybox1.35 Step 1/6 : ARG image_from_name
registry.stevedore.test/app1:v1-busybox1.35 Step 2/6 : ARG image_from_tag
registry.stevedore.test/app1:v1-busybox1.35 Step 3/6 : FROM ${image_from_name}:${image_from_tag}
registry.stevedore.test/app1:v1-busybox1.35 ‣  1.35:  Pulling from library/busybox 
registry.stevedore.test/app1:v1-busybox1.35 ‣  c15cbdab5f8e:  Pull complete 
registry.stevedore.test/app1:v1-busybox1.35 ‣  Digest: sha256:b4e4a06de46acc0958cd93e2eeb769077d255f06a7c3a91196509c16b7bc989e 
registry.stevedore.test/app1:v1-busybox1.35 ‣  Status: Downloaded newer image for busybox:1.35 
registry.stevedore.test/app1:v1-busybox1.35 ---> dddc7578369a
registry.stevedore.test/app1:v1-busybox1.35 Step 4/6 : COPY ./${app_name}/app.sh /app.sh
registry.stevedore.test/app1:v1-busybox1.35 ---> b94ccb3691f8
registry.stevedore.test/app1:v1-busybox1.35 Step 5/6 : CMD ["/app.sh"]
registry.stevedore.test/app1:v1-busybox1.35 ---> Running in adba76ed484f
registry.stevedore.test/app1:v1-busybox1.35 ---> a3e394dad1ca
registry.stevedore.test/app1:v1-busybox1.35 Step 6/6 : LABEL created_at=2023-05-19T13:06:03.410337541Z
registry.stevedore.test/app1:v1-busybox1.35 ---> Running in 086fc1384186
registry.stevedore.test/app1:v1-busybox1.35 ---> caf0c68df4b8
registry.stevedore.test/app1:v1-busybox1.35  ‣ sha256:caf0c68df4b80cc65320fd68a19a9f1d32013e008a65776e8392cfc1c8965f46
registry.stevedore.test/app1:v1-busybox1.35 Successfully built caf0c68df4b8
registry.stevedore.test/app1:v1-busybox1.35 Successfully tagged registry.stevedore.test/app1:v1-busybox1.35
registry.stevedore.test/app1:v1-busybox1.35 ‣  The push refers to repository [registry.stevedore.test/app1] 
registry.stevedore.test/app1:v1-busybox1.35 ‣  56a0dfcff812:  Pushed 
registry.stevedore.test/app1:v1-busybox1.35 ‣  42ef21f45b9a:  Pushed 
registry.stevedore.test/app1:v1-busybox1.35 ‣  v1-busybox1.35: digest: sha256:fd3fe99fd5f8c1aebd2d3feb38989d1a901af9029555dd81af51845401125ae4 size: 735
```

#### app2
To generate the image for app2, use the following command: `stevedore build app2 --push-after-build`. Keep in mind that two separate images of `app2` will be built, each with its respective parent image.
```sh
 [04-multiple-apps-example] Build app2 and push images after build
 Waiting for dockerd to be ready...
/certs/server/cert.pem: OK
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
/certs/client/cert.pem: OK
registry.stevedore.test/app2:v1-busybox1.36 Step 1/6 : ARG image_from_name
registry.stevedore.test/app2:v1-busybox1.35 Step 1/6 : ARG image_from_name
registry.stevedore.test/app2:v1-busybox1.36 Step 2/6 : ARG image_from_tag
registry.stevedore.test/app2:v1-busybox1.35 Step 2/6 : ARG image_from_tag
registry.stevedore.test/app2:v1-busybox1.35 Step 3/6 : FROM ${image_from_name}:${image_from_tag}
registry.stevedore.test/app2:v1-busybox1.36 Step 3/6 : FROM ${image_from_name}:${image_from_tag}
registry.stevedore.test/app2:v1-busybox1.35 ‣  1.35:  Pulling from library/busybox 
registry.stevedore.test/app2:v1-busybox1.35 ‣  c15cbdab5f8e:  Pull complete 
registry.stevedore.test/app2:v1-busybox1.35 ‣  Digest: sha256:b4e4a06de46acc0958cd93e2eeb769077d255f06a7c3a91196509c16b7bc989e 
registry.stevedore.test/app2:v1-busybox1.35 ‣  Status: Downloaded newer image for busybox:1.35 
registry.stevedore.test/app2:v1-busybox1.35 ---> dddc7578369a
registry.stevedore.test/app2:v1-busybox1.35 Step 4/6 : COPY ./${app_name}/app.sh /app.sh
registry.stevedore.test/app2:v1-busybox1.35 ---> 89523bafde14
registry.stevedore.test/app2:v1-busybox1.35 Step 5/6 : CMD ["/app.sh"]
registry.stevedore.test/app2:v1-busybox1.35 ---> Running in e1fd74b02e82
registry.stevedore.test/app2:v1-busybox1.35 ---> 1f02583a2af5
registry.stevedore.test/app2:v1-busybox1.35 Step 6/6 : LABEL created_at=2023-05-19T13:06:09.756209989Z
registry.stevedore.test/app2:v1-busybox1.35 ---> Running in 8da5f065ecd3
registry.stevedore.test/app2:v1-busybox1.35 ---> 4defe77ffaba
registry.stevedore.test/app2:v1-busybox1.35  ‣ sha256:4defe77ffaba4ed5a002085c367b9f8515064e302347bea58cabb4aebf47006d
registry.stevedore.test/app2:v1-busybox1.35 Successfully built 4defe77ffaba
registry.stevedore.test/app2:v1-busybox1.35 Successfully tagged registry.stevedore.test/app2:v1-busybox1.35
registry.stevedore.test/app2:v1-busybox1.35 ‣  The push refers to repository [registry.stevedore.test/app2] 
registry.stevedore.test/app2:v1-busybox1.36 ‣  1.36:  Pulling from library/busybox 
registry.stevedore.test/app2:v1-busybox1.36 ‣  a58ecd4f0c86:  Pull complete 
registry.stevedore.test/app2:v1-busybox1.35 ‣  v1-busybox1.35: digest: sha256:0cbc0000e62372f12883b16b89b05a9895ff96a2b2d8f441a552116879ca88c3 size: 735 
registry.stevedore.test/app2:v1-busybox1.36 ‣  Digest: sha256:9e2bbca079387d7965c3a9cee6d0c53f4f4e63ff7637877a83c4c05f2a666112 
registry.stevedore.test/app2:v1-busybox1.36 ‣  Status: Downloaded newer image for busybox:1.36 
registry.stevedore.test/app2:v1-busybox1.36 ---> af2c3e96bcf1
registry.stevedore.test/app2:v1-busybox1.36 Step 4/6 : COPY ./${app_name}/app.sh /app.sh
registry.stevedore.test/app2:v1-busybox1.36 ---> 35e88d8a038b
registry.stevedore.test/app2:v1-busybox1.36 Step 5/6 : CMD ["/app.sh"]
registry.stevedore.test/app2:v1-busybox1.36 ---> Running in 0f949fe3a66f
registry.stevedore.test/app2:v1-busybox1.36 ---> ace1dc5b86f5
registry.stevedore.test/app2:v1-busybox1.36 Step 6/6 : LABEL created_at=2023-05-19T13:06:09.755854242Z
registry.stevedore.test/app2:v1-busybox1.36 ---> Running in 4c7ac1499bc6
registry.stevedore.test/app2:v1-busybox1.36 ---> 9493698a1583
registry.stevedore.test/app2:v1-busybox1.36  ‣ sha256:9493698a1583325c281bfab6db57075c81bd52b9d257a4d7085b27d049b5f923
registry.stevedore.test/app2:v1-busybox1.36 Successfully built 9493698a1583
registry.stevedore.test/app2:v1-busybox1.36 Successfully tagged registry.stevedore.test/app2:v1-busybox1.36
registry.stevedore.test/app2:v1-busybox1.36 ‣  The push refers to repository [registry.stevedore.test/app2] 
registry.stevedore.test/app2:v1-busybox1.36 ‣  115f06ec0a74:  Layer already exists 
registry.stevedore.test/app2:v1-busybox1.36 ‣  1f1d08b81bbe:  Pushed 
registry.stevedore.test/app2:v1-busybox1.36 ‣  v1-busybox1.36: digest: sha256:f302125ff8d7cc62f38c0175bb3be36a0a7af1364fd6cc47caed7d46c199fcc4 size: 735 
```

#### app3
For `app3`, execute `stevedore build app3 --push-after-build` to trigger the build and image creation.
```sh
 [04-multiple-apps-example] Build app3 and push images after build
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
/certs/server/cert.pem: OK
 Waiting for dockerd to be ready...
/certs/client/cert.pem: OK
registry.stevedore.test/app3:v1-busybox1.36 Step 1/6 : ARG image_from_name
registry.stevedore.test/app3:v1-busybox1.36 Step 2/6 : ARG image_from_tag
registry.stevedore.test/app3:v1-busybox1.36 Step 3/6 : FROM ${image_from_name}:${image_from_tag}
registry.stevedore.test/app3:v1-busybox1.36 ‣  1.36:  Pulling from library/busybox 
registry.stevedore.test/app3:v1-busybox1.36 ‣  a58ecd4f0c86:  Pull complete 
registry.stevedore.test/app3:v1-busybox1.36 ‣  Digest: sha256:9e2bbca079387d7965c3a9cee6d0c53f4f4e63ff7637877a83c4c05f2a666112 
registry.stevedore.test/app3:v1-busybox1.36 ‣  Status: Downloaded newer image for busybox:1.36 
registry.stevedore.test/app3:v1-busybox1.36 ---> af2c3e96bcf1
registry.stevedore.test/app3:v1-busybox1.36 Step 4/6 : COPY ./${app_name}/app.sh /app.sh
registry.stevedore.test/app3:v1-busybox1.36 ---> bf88e91d0c29
registry.stevedore.test/app3:v1-busybox1.36 Step 5/6 : CMD ["/app.sh"]
registry.stevedore.test/app3:v1-busybox1.36 ---> Running in 7b54d6a9eb0d
registry.stevedore.test/app3:v1-busybox1.36 ---> cf4f7a9aae6e
registry.stevedore.test/app3:v1-busybox1.36 Step 6/6 : LABEL created_at=2023-05-19T13:06:16.592129227Z
registry.stevedore.test/app3:v1-busybox1.36 ---> Running in 0d0182cab3a6
registry.stevedore.test/app3:v1-busybox1.36 ---> d638d679d615
registry.stevedore.test/app3:v1-busybox1.36  ‣ sha256:d638d679d615eb3ea277c5b6fe335e88be7a8d35b0e505309cf6a5e14421f33f
registry.stevedore.test/app3:v1-busybox1.36 Successfully built d638d679d615
registry.stevedore.test/app3:v1-busybox1.36 Successfully tagged registry.stevedore.test/app3:v1-busybox1.36
registry.stevedore.test/app3:v1-busybox1.36 ‣  The push refers to repository [registry.stevedore.test/app3] 
registry.stevedore.test/app3:v1-busybox1.36 ‣  caa4e783cf93:  Pushed 
registry.stevedore.test/app3:v1-busybox1.36 ‣  1f1d08b81bbe:  Pushed 
registry.stevedore.test/app3:v1-busybox1.36 ‣  v1-busybox1.36: digest: sha256:5086c7e5f3b4d2b3e31a89b944cda09c625eadef2127fc49ef0eff92c302bbcd size: 735 
```

### Cleaning the stack
```sh
Stopping the stack to run 04-multiple-apps-example

[+] Running 8/8
 ✔ Container 04-multiple-apps-example-registry-1                  Removed                                                                                   0.3s 
 ✔ Container 04-multiple-apps-example-stevedore-run-ab79586e3977  Removed                                                                                   0.0s 
 ✔ Container 04-multiple-apps-example-stevedore-run-900d67839ff2  Removed                                                                                   0.0s 
 ✔ Container 04-multiple-apps-example-stevedore-run-b60aadd80da6  Removed                                                                                   0.0s 
 ✔ Container 04-multiple-apps-example-stevedore-run-b31f6304a414  Removed                                                                                   0.0s 
 ✔ Container 04-multiple-apps-example-stevedore-1                 Removed                                                                                  10.2s 
 ✔ Container 04-multiple-apps-example-dockerauth-1                Removed                                                                                   0.2s 
 ✔ Network 04-multiple-apps-example_default                       Removed                                                                                   0.4s 
```