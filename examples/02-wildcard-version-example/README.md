# Wildcard version example

This example serves as an introduction to the [wildcard version](https://gostevedore.github.io/docs/getting-started/concepts/#wildcard-version) feature in Stevedore. Building upon the foundation set by the [01-basic-example](https://github.com/gostevedore/stevedore/tree/main/examples/01-basic-example), it demonstrates the addition of a wildcard version to the `my-app` image definition. Furthermore, it guides you through the process of creating an image using the wildcard version as an example.

- [Wildcard version example](#wildcard-version-example)
  - [Requirements](#requirements)
  - [Stack](#stack)
  - [Usage](#usage)
  - [Expected Output](#expected-output)
    - [Starting the stack](#starting-the-stack)
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
The stack required to run this example is defined in that [Docker Compose file](./docker-compose.yml). The stack consists of three services: a Docker Registry, a Docker Registry authorization and a Stevedore service. The Docker registry is used to store the Docker images built by Stevedore during the example execution. The Stevedore service is where the example is executed.

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
Below is the expected output for the `make run` command, which starts the Docker stack, gets some information about the Stevedore configuration, builds and promotes a Docker image using Stevedore, and then cleans the stack up.

### Starting the stack
```sh
Starting the stack to run 02-wildcard-definition-example

[+] Building 8.6s (21/21) FINISHED
 => [internal] load .dockerignore                                                                                                                                    0.0s
 => => transferring context: 2B                                                                                                                                      0.0s
 => [internal] load build definition from Dockerfile                                                                                                                 0.0s
 => => transferring dockerfile: 989B                                                                                                                                 0.0s
 => [internal] load metadata for docker.io/library/docker:20.10-dind                                                                                                 1.1s
 => [internal] load metadata for docker.io/library/golang:1.19-alpine                                                                                                1.0s
 => [golang 1/8] FROM docker.io/library/golang:1.19-alpine@sha256:4147c2ad24a347e53a8a361b663dcf82fca4157a3f9a5136d696a4e53dd22b65                                   0.0s
 => => resolve docker.io/library/golang:1.19-alpine@sha256:4147c2ad24a347e53a8a361b663dcf82fca4157a3f9a5136d696a4e53dd22b65                                          0.0s
 => [stage-1 1/6] FROM docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                                   0.0s
 => => resolve docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                                           0.0s
 => [internal] load build context                                                                                                                                    0.0s
 => => transferring context: 94.72kB                                                                                                                                 0.0s
 => CACHED [golang 2/8] WORKDIR /usr/src/app                                                                                                                         0.0s
 => CACHED [golang 3/8] RUN apk add --no-cache make build-base                                                                                                       0.0s
 => CACHED [golang 4/8] COPY go.mod ./                                                                                                                               0.0s
 => CACHED [golang 5/8] COPY go.sum ./                                                                                                                               0.0s
 => CACHED [golang 6/8] RUN go mod download && go mod verify                                                                                                         0.0s
 => [golang 7/8] COPY . ./                                                                                                                                           0.5s
 => [golang 8/8] RUN go build -ldflags "-s -w" -v -o /usr/local/bin/stevedore ./cmd/stevedore.go                                                                     4.9s
 => CACHED [stage-1 2/6] COPY --from=golang /usr/local/go /usr/local/go                                                                                              0.0s
 => CACHED [stage-1 3/6] COPY --from=golang /usr/local/bin/stevedore /usr/local/bin/stevedore                                                                        0.0s
 => CACHED [stage-1 4/6] WORKDIR /go                                                                                                                                 0.0s
 => CACHED [stage-1 5/6] RUN mkdir -p "/go/src" "/go/bin" && chmod -R 777 "/go"                                                                                      0.0s
 => CACHED [stage-1 6/6] COPY test/stack/client/entrypoint.sh /usr/local/bin/entrypoint.sh                                                                           0.0s
 => exporting to docker image format                                                                                                                                 0.9s
 => => exporting layers                                                                                                                                              0.0s
 => => exporting manifest sha256:92e8e2d0cbf1e5d733c7b2b60e6e018313fbe66f785da8d4b3f386f515c8abba                                                                    0.0s
 => => exporting config sha256:abab447a0612fe36411f890747e176cbf37c44d5be04af96a5b5c47610664f6b                                                                      0.0s
 => => sending tarball                                                                                                                                               0.9s
 => importing to docker                                                                                                                                              0.0s
[+] Running 4/4
 ✔ Network 02-wildcard-definition-example_default         Created                                                                                                    0.1s
 ✔ Container 02-wildcard-definition-example-dockerauth-1  Started                                                                                                    0.5s
 ✔ Container 02-wildcard-definition-example-stevedore-1   Started                                                                                                    0.5s
 ✔ Container 02-wildcard-definition-example-registry-1    Started                                                                                                    0.7s
```

### Getting images
To view the images, run `stevedore get images`.

```sh
 Run example 02-wildcard-definition-example

 [02-wildcard-definition-example] Get images
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
/certs/server/cert.pem: OK
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
/certs/client/cert.pem: OK
NAME    VERSION                                                BUILDER IMAGE FULL NAME                                                                       PARENT
busybox 1.35                                                   -       busybox:1.35                                                                          -
busybox 1.36                                                   -       busybox:1.36                                                                          -
my-app  3.2.1-busybox1.35                                      my-app  registry.stevedore.test/my-app:3.2.1-busybox1.35                                      busybox:1.35
my-app  3.2.1-busybox1.36                                      my-app  registry.stevedore.test/my-app:3.2.1-busybox1.36                                      busybox:1.36
my-app  {{ .Version }}-{{ .Parent.Name }}{{ .Parent.Version }} my-app  registry.stevedore.test/my-app:{{ .Version }}-{{ .Parent.Name }}{{ .Parent.Version }} busybox:1.35
```

When you obtain the images list, you can notice that there is a version for `my-app` that is not explicitly defined. This version is indicated by the value `{{ .Version }}-{{ .Parent.Name }}{{ .Parent.Version }}`. The final version is resolved when you provide a specific value during the Docker image building process.

For further information on templating image attributes, please refer to the Stevedore [reference guide](https://gostevedore.github.io/docs/reference-guide/image/#templating-image-attributes).

### Building images
The example uses the command `stevedore build my-app --image-version example --push-after-build` to build and automatically promote the images to the Docker registry.

If you review the builder definition, you will notice that the source code for `my-app` is located in the [./builders/apps.yaml](builders/apps.yaml) file. This folder contains the necessary resources required for building the `my-app` Docker image.

Please note the presence of the `--image-version example` flag in the build command. In this case, since `my-app` does not have an explicit version definition for `example` and instead has a wildcard version defined, the build process utilizes the wildcard definition.

```sh
 [02-wildcard-definition-example] Build my-app and push images after build
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
/certs/server/cert.pem: OK
 Waiting for dockerd to be ready...
/certs/client/cert.pem: OK
registry.stevedore.test/my-app:example-busybox1.35 Step 1/7 : ARG image_from_name
registry.stevedore.test/my-app:example-busybox1.35 Step 2/7 : ARG image_from_tag
registry.stevedore.test/my-app:example-busybox1.35 Step 3/7 : FROM ${image_from_name}:${image_from_tag}
registry.stevedore.test/my-app:example-busybox1.35 ‣  1.35:  Pulling from library/busybox
registry.stevedore.test/my-app:example-busybox1.35 ‣  c15cbdab5f8e:  Pull complete
registry.stevedore.test/my-app:example-busybox1.35 ‣  Digest: sha256:b4e4a06de46acc0958cd93e2eeb769077d255f06a7c3a91196509c16b7bc989e
registry.stevedore.test/my-app:example-busybox1.35 ‣  Status: Downloaded newer image for busybox:1.35
registry.stevedore.test/my-app:example-busybox1.35 ---> dddc7578369a
registry.stevedore.test/my-app:example-busybox1.35 Step 4/7 : ARG whoami=unknown
registry.stevedore.test/my-app:example-busybox1.35 ---> Running in 9e6f0e4f81ca
registry.stevedore.test/my-app:example-busybox1.35 ---> a4053684b87f
registry.stevedore.test/my-app:example-busybox1.35 Step 5/7 : RUN echo "Hey there, I'm ${whoami}!" > /whoami.txt
registry.stevedore.test/my-app:example-busybox1.35 ---> Running in 1946b2a47c7c
registry.stevedore.test/my-app:example-busybox1.35 ---> 50ba5a0655b6
registry.stevedore.test/my-app:example-busybox1.35 Step 6/7 : CMD ["cat","/whoami.txt"]
registry.stevedore.test/my-app:example-busybox1.35 ---> Running in 074ce3ddde51
registry.stevedore.test/my-app:example-busybox1.35 ---> d428d8a4bb57
registry.stevedore.test/my-app:example-busybox1.35 Step 7/7 : LABEL created_at=2023-05-12T10:45:20.86119333Z
registry.stevedore.test/my-app:example-busybox1.35 ---> Running in 243f55598fe6
registry.stevedore.test/my-app:example-busybox1.35 ---> 6938c8420c3a
registry.stevedore.test/my-app:example-busybox1.35  ‣ sha256:6938c8420c3a1b97d8b71d1a09da095488e2e3c9332ae156deacbf8784a31a96
registry.stevedore.test/my-app:example-busybox1.35 Successfully built 6938c8420c3a
registry.stevedore.test/my-app:example-busybox1.35 Successfully tagged registry.stevedore.test/my-app:example-busybox1.35
registry.stevedore.test/my-app:example-busybox1.35 ‣  The push refers to repository [registry.stevedore.test/my-app]
registry.stevedore.test/my-app:example-busybox1.35 ‣  76cf1635ccd6:  Pushed
registry.stevedore.test/my-app:example-busybox1.35 ‣  42ef21f45b9a:  Pushed
registry.stevedore.test/my-app:example-busybox1.35 ‣  example-busybox1.35: digest: sha256:ba2647a7e9dff796918ef638c9c98200895bf06f90859837f3f1ba1c366cdbeb size: 735
```

Upon executing the build command, the resulting image `registry.stevedore.test/my-app:example-busybox1.35` is successfully created and pushed to the Docker registry.

### Cleaning the stack
```sh
Stopping the stack to run 02-wildcard-definition-example

[+] Running 6/6
 ✔ Container 02-wildcard-definition-example-registry-1                  Removed                                              0.3s
 ✔ Container 02-wildcard-definition-example-stevedore-1                 Removed                                             10.2s
 ✔ Container 02-wildcard-definition-example-stevedore-run-e5cf3245e4f1  Removed                                              0.0s
 ✔ Container 02-wildcard-definition-example-stevedore-run-c109f6b7425c  Removed                                              0.0s
 ✔ Container 02-wildcard-definition-example-dockerauth-1                Removed                                              0.2s
 ✔ Network 02-wildcard-definition-example_default                       Removed                                              0.4s
```

## Additional information
In addition to the core steps outlined in the example, the following section provides additional information and insights to further enhance your understanding of how this example uses Stevedore.

### Images
The wildcard version in Stevedore is denoted by the asterisk symbol `*`. In the image definitions, you can use `{{ .Version }}` to represent the version attribute. During the building process, when the `--image-version` flag is specified, all occurrences of `{{ .Version }}` within the image definitions are replaced with the value provided in the `--image-version` flag. This allows you to dynamically set the version of the Docker images based on the command-line input, providing flexibility and customization in your image building process.

Here is an example of how a wildcard version is defined in the Stevedore image definitions:
```yaml
images:
  my-app:
    "*":
      version: "{{ .Version }}-{{ .Parent.Name }}{{ .Parent.Version }}"
      registry: registry.stevedore.test
      builder: my-app
      parents:
        busybox:
          - "1.35"
      vars:
        whoami: "{{ .Version }}-{{ .Parent.Name }}{{ .Parent.Version }}"
```

When you execute the resulting image, you can see the output from the newly created image:
```sh
/app/examples/02-wildcard-version-example $ docker run --rm registry.stevedore.test/my-app:example-busybox1.35
Hey there, I'm example-busybox1.35!
```
This demonstrates how the wildcard version allows you to dynamically generate unique version identifiers for your Docker images based on the provided input through the `--image-version` flag.