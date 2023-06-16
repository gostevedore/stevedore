# Custom Variables Mapping Example

This example showcases the utilization of custom variables mapping in Stevedore, enabling flexible configuration and customization of build arguments in the Dockerfile.

By default, Stevedore automatically provides the `image_from_name` and `image_from_tag` arguments to the Docker build, containing information about the parent image. However, in this example, these default arguments are overwritten within the [builder](https://gostevedore.github.io/docs/reference-guide/builder/) configuration. Instead, the variables  `parent_name` and `parent_tag` are used to customize the build arguments.

- [Custom Variables Mapping Example](#custom-variables-mapping-example)
  - [Requirements](#requirements)
  - [Stack](#stack)
  - [Usage](#usage)
  - [Example Execution Insights](#example-execution-insights)
    - [Starting the stack](#starting-the-stack)
    - [Waiting for Dockerd to be Ready](#waiting-for-dockerd-to-be-ready)
    - [Execute Build in dry-run mode](#execute-build-in-dry-run-mode)
    - [Cleaning the stack](#cleaning-the-stack)
  - [Additional information](#additional-information)
    - [Builder](#builder)


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
Starting the stack to run 09-custom-variables-mapping-example

[+] Building 8.3s (18/18) FINISHED
 => [stevedore internal] load build definition from Dockerfile                                                                                                       0.0s
 => => transferring dockerfile: 888B                                                                                                                                 0.0s
 => [stevedore internal] load .dockerignore                                                                                                                          0.0s
 => => transferring context: 2B                                                                                                                                      0.0s
 => [stevedore internal] load metadata for docker.io/library/docker:20.10-dind                                                                                       1.4s
 => [stevedore internal] load metadata for docker.io/library/golang:1.19-alpine                                                                                      1.8s
 => [stevedore internal] load build context                                                                                                                          0.1s
 => => transferring context: 372.54kB                                                                                                                                0.1s
 => [stevedore golang 1/7] FROM docker.io/library/golang:1.19-alpine@sha256:4147c2ad24a347e53a8a361b663dcf82fca4157a3f9a5136d696a4e53dd22b65                         0.0s
 => => resolve docker.io/library/golang:1.19-alpine@sha256:4147c2ad24a347e53a8a361b663dcf82fca4157a3f9a5136d696a4e53dd22b65                                          0.0s
 => [stevedore stage-1 1/4] FROM docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                         0.0s
 => => resolve docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                                           0.0s
 => CACHED [stevedore golang 2/7] WORKDIR /usr/src/app                                                                                                               0.0s
 => CACHED [stevedore golang 3/7] RUN apk add --no-cache make build-base                                                                                             0.0s
 => CACHED [stevedore golang 4/7] COPY go.mod go.sum ./                                                                                                              0.0s
 => CACHED [stevedore golang 5/7] RUN go mod download                                                                                                                0.0s
 => [stevedore golang 6/7] COPY . ./                                                                                                                                 0.5s
 => [stevedore golang 7/7] RUN go build -ldflags "-s -w" -v -o /usr/local/bin/stevedore ./cmd/stevedore.go                                                           5.3s
 => CACHED [stevedore stage-1 2/4] COPY --from=golang /usr/local/bin/stevedore /usr/local/bin/stevedore                                                              0.0s
 => CACHED [stevedore stage-1 3/4] COPY examples/01-basic-example/stack/stevedore/entrypoint.sh /usr/local/bin/entrypoint.sh                                         0.0s
 => CACHED [stevedore stage-1 4/4] COPY examples/01-basic-example/stack/stevedore/wait-for-dockerd.sh /usr/local/bin/wait-for-dockerd.sh                             0.0s
 => [stevedore] exporting to docker image format                                                                                                                     0.5s
 => => exporting layers                                                                                                                                              0.0s
 => => exporting manifest sha256:6dcc99e5af0a53c9a62a0a4cf0e04d7bcfff62858991c70df12d3856947b0a44                                                                    0.0s
 => => exporting config sha256:fe4178704bcb16dc9a4f5228529ea424410f9b258e694d90a73a563168328886                                                                      0.0s
 => => sending tarball                                                                                                                                               0.5s
 => [stevedore stevedore] importing to docker                                                                                                                        0.0s
[+] Running 4/4
 ✔ Network 09-custom-variables-mapping-example_default         Created                                                                                               0.1s
 ✔ Container 09-custom-variables-mapping-example-stevedore-1   Started                                                                                               0.5s
 ✔ Container 09-custom-variables-mapping-example-dockerauth-1  Started                                                                                               0.4s
 ✔ Container 09-custom-variables-mapping-example-registry-1    Started                                                                                               0.7s
```

### Waiting for Dockerd to be Ready

Before starting the execution of the Stevedore command, it is important to ensure that the Docker daemon (dockerd) is ready. The stevedore service Docker image includes a script, [wait-for-dockerd.sh](./stack/stevedore/wait-for-dockerd.sh), which can be used to ensure the readiness of the Docker daemon.

```sh
 Run example 09-custom-variables-mapping-example

 [09-custom-variables-mapping-example] Waiting for dockerd
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
```

### Execute Build in dry-run mode

To begin, the example performs a Docker image build in dry-run mode using the following command: `stevedore build my-app --image-version 3.2.1-busybox1.36 --dry-run`. The dry-run mode allows for the exploration of all the parameters utilized during the build process, including the [variables-mapping](https://gostevedore.github.io/docs/getting-started/concepts/#variables-mapping) configuration.

Here is the output generated during the dry-run:

```sh
 [09-custom-variables-mapping-example] Build my-app and push images after build

 builder:       my-app
 lables: map[]
 name: my-app
 parent:
 - busybox:1.36
 presistent labels: map[created_at:2023-06-01T05:25:36.730005762Z]
 presistent vars: map[]
 registry host: registry.stevedore.test
 registry namespace:
 tags: []
 vars: map[whoami:3.2.1-busybox1.36]
 version: 3.2.1-busybox1.36
 options:
  ansible_connection_local: false
  ansible_intermediate_container_name: builder_docker__my-app_3.2.1-busybox1.36
  ansible_inventory_path: ""
  ansible_limit: ""
  builder_options:
    playbook: ""
    inventory: ""
    dockerfile: ""
    context:
    - path: ./my-app
  builder_variables_mapping:
    image_builder_label_key: image_builder_label
    image_builder_name_key: image_builder_name
    image_builder_registry_host_key: image_builder_registry_host
    image_builder_registry_namespace_key: image_builder_registry_namespace
    image_builder_tag_key: image_builder_tag
    image_extra_tags_key: image_extra_tags
    image_from_name_key: parent_name
    image_from_registry_host_key: image_from_registry_host
    image_from_registry_namespace_key: image_from_registry_namespace
    image_from_tag_key: parent_version
    image_lables_key: image_labels
    image_name_key: image_name
    image_registry_host_key: image_registry_host
    image_registry_namespace_key: image_registry_namespace
    image_tag_key: image_tag
    pull_parent_image_key: pull_parent_image
    push_image_key: push_image
  output_prefix: ""
  pull_auth_username: ""
  pull_parent_image: false
  push_auth_username: admin
  push_image_after_build: false
  remove_image_after_build: false
 parent builder vars mapping:
  parent_name: busybox
  parent_version: 1.36
```

Within the `builder_variables_mapping` block, you can observe the argument names that will be passed to the Docker API for image creation. In this case, the customized arguments are `image_from_name_key: parent_name` and `image_from_tag_key: parent_version`.

Furthermore, in the `parent builder vars mapping` block, you can find the corresponding values that will be passed to the Docker API:

```yaml
parent builder vars mapping:
  parent_name: busybox
  parent_version: 1.36
```

This demonstrates how the custom variables mapping allows for flexible configuration and customization of build arguments in Stevedore.

After understanding the customization of the variable mapping, you can proceed with the build of the images.

```sh
 [09-custom-variables-mapping-example] Build my-app and push images after build
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 1/7 : ARG parent_name
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 2/7 : ARG parent_version
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 3/7 : FROM ${parent_name}:${parent_version}
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 1/7 : ARG parent_name
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 2/7 : ARG parent_version
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 3/7 : FROM ${parent_name}:${parent_version}

registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  1.36:  Pulling from library/busybox
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  325d69979d33:  Pull complete
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  Digest: sha256:b4e4a06de46acc0958cd93e2eeb769077d255f06a7c3a91196509c16b7bc989e
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  Digest: sha256:560af6915bfc8d7630e50e212e08242d37b63bd5c1ccf9bd4acccf116e262d5b
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  Status: Downloaded newer image for busybox:1.35
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  Status: Downloaded newer image for busybox:1.36
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> dddc7578369a
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 4/7 : ARG whoami=unknown
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> 8135583d97fe
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 4/7 : ARG whoami=unknown
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> Running in 6addf8af10bf
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> Running in 36f7fa62cd0b
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> 3d0b214e075d
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 5/7 : RUN echo "Hey there, I'm ${whoami}!" > /whoami.txt
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> c67ccbfbabf7
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 5/7 : RUN echo "Hey there, I'm ${whoami}!" > /whoami.txt
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> Running in fe7e92ee442e
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> Running in 715c090d3889
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> dd420075f33c
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 6/7 : CMD ["cat","/whoami.txt"]
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> Running in 8df85864bc82
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> 4dbd84131040
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 7/7 : LABEL created_at=2023-06-01T05:25:36.860355092Z
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> Running in 0e9cb61d79c8
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> c1ff52972445
registry.stevedore.test/my-app:3.2.1-busybox1.35  ‣ sha256:c1ff5297244589ac3bef9fd308d0e4e3fa4b843a1d94942eeaaa696889c6b5fa
registry.stevedore.test/my-app:3.2.1-busybox1.35 Successfully built c1ff52972445
registry.stevedore.test/my-app:3.2.1-busybox1.35 Successfully tagged registry.stevedore.test/my-app:3.2.1-busybox1.35
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  The push refers to repository [registry.stevedore.test/my-app]
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  f8339ba34ace:  Preparing
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  42ef21f45b9a:  Preparing
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> 91795c96264c
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 6/7 : CMD ["cat","/whoami.txt"]
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  f8339ba34ace:  Pushed
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  42ef21f45b9a:  Pushing [>                                                  ]  66.56kB/4.855MB
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> c9317cb18668
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 7/7 : LABEL created_at=2023-06-01T05:25:36.861921762Z
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  f8339ba34ace:  Pushed
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  42ef21f45b9a:  Pushing [========================>                          ]  2.406MB/4.855MB
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> b53ea67449c6
registry.stevedore.test/my-app:3.2.1-busybox1.36  ‣ sha256:b53ea67449c620c431e86fd11aeb810bb76233fe94bcb1f8bcd35dbd97c73910
registry.stevedore.test/my-app:3.2.1-busybox1.36 Successfully built b53ea67449c6
registry.stevedore.test/my-app:3.2.1-busybox1.36 Successfully tagged registry.stevedore.test/my-app:3.2.1-busybox1.36
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  The push refers to repository [registry.stevedore.test/my-app]
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  ad6e7ee3b6ad:  Pushed
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  ad6e7ee3b6ad:  Pushed
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  9547b4c33213:  Pushed
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  3.2.1-busybox1.36: digest: sha256:67e87f2a624c9294028c1af2126c0b1c7a62d09c912285d640bba3769e5482a9 size: 735
```

### Cleaning the stack

```sh
Stopping the stack to run 09-custom-variables-mapping-example

[+] Running 4/4
 ✔ Container 09-custom-variables-mapping-example-stevedore-1   Removed                                                                                               3.3s
 ✔ Container 09-custom-variables-mapping-example-registry-1    Removed                                                                                               0.2s
 ✔ Container 09-custom-variables-mapping-example-dockerauth-1  Removed                                                                                               0.2s
 ✔ Network 09-custom-variables-mapping-example_default         Removed                                                                                               0.4s
```

## Additional information

In addition to the core steps outlined in the example, the following section provides additional information and insights to further enhance your understanding of how this example uses Stevedore.

### Builder
The `variables_mapping` parameter in the [builder](https://gostevedore.github.io/docs/reference-guide/builder/) configuration allows you to customize the automatic arguments passed to the Docker API and consumed by the Dockerfile during the image building process. In this example, the `image_from_name_key` and `image_from_tag_key` arguments are customized to `parent_name` and `parent_version`, respectively.

Here is an example of the builder configuration showcasing the customization of the variables mapping:

```yaml
my-app:
  driver: docker
  options:
    context:
      - path: ./my-app
  variables_mapping:
    image_from_name_key: parent_name
    image_from_tag_key: parent_version
```

This configuration ensures that the customized arguments `parent_name` and `parent_version` are used when communicating with the Docker API during the image building process.

Additionally, the Dockerfile in this example utilizes the customized build arguments `parent_name` and `parent_version`:

```dockerfile
ARG parent_name
ARG parent_version

FROM ${parent_name}:${parent_version}

ARG whoami=unknown

RUN echo "Hey there, I'm ${whoami}!" > /whoami.txt

CMD ["cat","/whoami.txt"]
```

The values of `parent_name` and `parent_version` will be substituted with the corresponding values provided during the build process.
