# Use Ansible Playbook Driver Example

This example showcases the usage of the [Ansible playbook driver](https://gostevedore.github.io/docs/reference-guide/driver/ansible-playbook/) in Stevedore, enabling you to build Docker images using Ansible playbooks. By using as a base image an image build with the `ansible-playbook`, the example demostrates the flexibility of using multiple drivers within your image definitions.

- [Use Ansible Playbook Driver Example](#use-ansible-playbook-driver-example)
  - [Requirements](#requirements)
  - [Stack](#stack)
  - [Usage](#usage)
  - [Example Execution Insights](#example-execution-insights)
    - [Starting the stack](#starting-the-stack)
    - [Waiting for Dockerd to be Ready](#waiting-for-dockerd-to-be-ready)
    - [Building the base image](#building-the-base-image)
    - [Building the app1 image](#building-the-app1-image)
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
Below is the expected output for the `make run` command, which starts the Docker stack, gets some information about the Stevedore configuration, builds and promotes a Docker image using Stevedore, and then cleans the stack up.

Starting the stack to run 10-use-ansible-playbook-driver-example

[+] Building 74.1s (24/24) FINISHED
 => [stevedore internal] load .dockerignore                                                                                                                          0.0s
 => => transferring context: 2B                                                                                                                                      0.0s
 => [stevedore internal] load build definition from Dockerfile                                                                                                       0.0s
 => => transferring dockerfile: 2.57kB                                                                                                                               0.0s
 => [stevedore internal] load metadata for docker.io/library/debian:bookworm                                                                                         1.0s
 => [stevedore internal] load metadata for docker.io/library/docker:20.10-dind                                                                                       1.1s
 => [stevedore internal] load metadata for docker.io/library/golang:1.19-alpine                                                                                      1.1s
 => [stevedore golang 1/7] FROM docker.io/library/golang:1.19-alpine@sha256:470c8d0638c5b7007a6118baee531c30e0516a18e45b35bff1f8ab92cf8f896d                         0.0s
 => => resolve docker.io/library/golang:1.19-alpine@sha256:470c8d0638c5b7007a6118baee531c30e0516a18e45b35bff1f8ab92cf8f896d                                          0.0s
 => [stevedore internal] load build context                                                                                                                          0.1s
 => => transferring context: 353.83kB                                                                                                                                0.0s
 => CACHED [stevedore dind 1/1] FROM docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                     0.0s
 => => resolve docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                                           0.0s
 => [stevedore stage-2 1/8] FROM docker.io/library/debian:bookworm@sha256:d568e251e460295a8743e9d5ef7de673c5a8f9027db11f4e666e96fb5bed708e                           0.0s
 => => resolve docker.io/library/debian:bookworm@sha256:d568e251e460295a8743e9d5ef7de673c5a8f9027db11f4e666e96fb5bed708e                                             0.0s
 => CACHED [stevedore stage-2 2/8] RUN apt-get update   && apt-get install --no-install-recommends --yes     build-essential     curl     git     libffi-dev     li  0.0s
 => [stevedore stage-2 3/8] RUN pip3 install --break-system-packages   ansible==5.9.0   cryptography   docker   requests                                            53.5s
 => CACHED [stevedore golang 2/7] WORKDIR /usr/src/app                                                                                                               0.0s
 => CACHED [stevedore golang 3/7] RUN apk add --no-cache make build-base                                                                                             0.0s
 => CACHED [stevedore golang 4/7] COPY go.mod go.sum ./                                                                                                              0.0s
 => CACHED [stevedore golang 5/7] RUN go mod download                                                                                                                0.0s
 => CACHED [stevedore golang 6/7] COPY . ./                                                                                                                          0.0s
 => CACHED [stevedore golang 7/7] RUN CGO_ENABLED=0 go build -ldflags "-s -w -X 'github.com/gostevedore/stevedore/internal/core/domain/release.BuildDate=$(date +%c  0.0s
 => [stevedore stage-2 4/8] COPY --from=golang /usr/local/bin/stevedore /usr/local/bin/stevedore                                                                     0.3s
 => [stevedore stage-2 5/8] COPY --from=dind /usr/local/bin/dockerd-entrypoint.sh /usr/local/bin/dockerd-entrypoint.sh                                               0.0s
 => [stevedore stage-2 6/8] COPY examples/10-use-ansible-playbook-driver-example/stack/stevedore/entrypoint.sh /usr/local/bin/entrypoint.sh                          0.0s
 => [stevedore stage-2 7/8] COPY examples/10-use-ansible-playbook-driver-example/stack/stevedore/wait-for-dockerd.sh /usr/local/bin/wait-for-dockerd.sh              0.0s
 => [stevedore stage-2 8/8] WORKDIR /src                                                                                                                             0.0s
 => [stevedore] exporting to docker image format                                                                                                                    19.1s
 => => exporting layers                                                                                                                                             12.4s
 => => exporting manifest sha256:3577b14c00f6c62cac3f8858b7e912e6d0512ad0d7913a1de67f4bf789d51576                                                                    0.0s
 => => exporting config sha256:372782712c6529e0552431d4741654c3ef976bbd8f00d0da1e69c506edcf11c7                                                                      0.0s
 => => sending tarball                                                                                                                                               6.7s
 => [stevedore stevedore] importing to docker                                                                                                                        4.5s
[+] Running 4/4
 ✔ Network 10-use-ansible-playbook-driver-example_default         Created                                                                                            0.1s
 ✔ Container 10-use-ansible-playbook-driver-example-stevedore-1   Started                                                                                            0.9s
 ✔ Container 10-use-ansible-playbook-driver-example-dockerauth-1  Started                                                                                            0.8s
 ✔ Container 10-use-ansible-playbook-driver-example-registry-1    Started                                                                                            0.8s
```

### Waiting for Dockerd to be Ready

Before starting the execution of the Stevedore command, it is important to ensure that the Docker daemon (dockerd) is ready. The stevedore service Docker image includes a script, [wait-for-dockerd.sh](./stack/stevedore/wait-for-dockerd.sh), which can be used to ensure the readiness of the Docker daemon.

```sh
 Run example 10-use-ansible-playbook-driver-example

 [10-use-ansible-playbook-driver-example] Waiting for dockerd
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
```

### Building the base image

In this example, the command `stevedore build ubuntu --ansible-connection-local` creates the `Ubuntu` Docker image, resulting in the locally available image `registry.stevedore.test/ubuntu:focal`.

To create a Docker image using Ansible, you start a base container for provisioning. In this example, the building process is performed locally using the `--ansible-connection-local` flag.

```sh
 [10-use-ansible-playbook-driver-example] Building Ubuntu image
ubuntu:focal ──
ubuntu:focal ── PLAY [Create builder stage] ****************************************************
ubuntu:focal ──
ubuntu:focal ── TASK [Gathering Facts] *********************************************************
ubuntu:focal ── ok: [docker_image_builder]
ubuntu:focal ──
ubuntu:focal ── TASK [Creating container builder 'builder_ansible-playbook__ubuntu_focal'] *****
ubuntu:focal ── changed: [docker_image_builder]
ubuntu:focal ──
ubuntu:focal ── TASK [Creating container builder output message] *******************************
ubuntu:focal ── skipping: [docker_image_builder]
ubuntu:focal ──
ubuntu:focal ── TASK [Add 'builder_ansible-playbook__ubuntu_focal' to ansible inventory] *******
ubuntu:focal ── changed: [docker_image_builder]
```

When the base image is ready, you can start provisioning the container and perform various tasks. For instance, in this example, a user is created.

```sh
ubuntu:focal ── PLAY [Provision image container builder stage] *********************************
ubuntu:focal ──
ubuntu:focal ── TASK [Include project variables] ***********************************************
ubuntu:focal ── ok: [builder_ansible-playbook__ubuntu_focal]
ubuntu:focal ──
ubuntu:focal ── TASK [Install Python dependencies on builder_ansible-playbook__ubuntu_focal] ***
ubuntu:focal ── changed: [builder_ansible-playbook__ubuntu_focal] => (item=test -e /usr/bin/python3 || (apt-get update && apt-get install -y python3-minimal python3-apt))
ubuntu:focal ──
ubuntu:focal ── TASK [Upgrade base packages] ***************************************************
ubuntu:focal ── changed: [builder_ansible-playbook__ubuntu_focal]
ubuntu:focal ──
ubuntu:focal ── TASK [Create user apenella] ****************************************************
ubuntu:focal ── changed: [builder_ansible-playbook__ubuntu_focal]
ubuntu:focal ──
ubuntu:focal ── TASK [Uninstall Python package on container builder 'builder_ansible-playbook__ubuntu_focal'] ***
ubuntu:focal ── changed: [builder_ansible-playbook__ubuntu_focal] => (item=apt-get purge -y --auto-remove python3-minimal python3-apt)
ubuntu:focal ──
ubuntu:focal ── TASK [Clean packages and dependencies on container builder 'builder_ansible-playbook__ubuntu_focal'] ***
ubuntu:focal ── changed: [builder_ansible-playbook__ubuntu_focal] => (item=apt-get autoremove -y)
ubuntu:focal ── changed: [builder_ansible-playbook__ubuntu_focal] => (item=apt-get clean -y)
ubuntu:focal ── changed: [builder_ansible-playbook__ubuntu_focal] => (item=apt-get autoclean -y)
ubuntu:focal ── changed: [builder_ansible-playbook__ubuntu_focal] => (item=rm -rf /tmp/* /var/tmp/* /var/lib/apt/lists/*)
```

Finally, you can create the Docker image by committing the provisioned container as the base.

```sh
ubuntu:focal ── PLAY [Create image stage] ******************************************************
ubuntu:focal ──
ubuntu:focal ── TASK [Gathering Facts] *********************************************************
ubuntu:focal ── ok: [docker_image_builder]
ubuntu:focal ──
ubuntu:focal ── TASK [Include project variables] ***********************************************
ubuntu:focal ── ok: [docker_image_builder]
ubuntu:focal ──
ubuntu:focal ── TASK [Extend dockerfile instructions] ******************************************
ubuntu:focal ── ok: [docker_image_builder]
ubuntu:focal ──
ubuntu:focal ── TASK [Set image name] **********************************************************
ubuntu:focal ── ok: [docker_image_builder]
ubuntu:focal ──
ubuntu:focal ── TASK [Set image tag] ***********************************************************
ubuntu:focal ── ok: [docker_image_builder]
ubuntu:focal ──
ubuntu:focal ── TASK [Set image registry namespace] ********************************************
ubuntu:focal ── skipping: [docker_image_builder]
ubuntu:focal ──
ubuntu:focal ── TASK [Set image registry host] *************************************************
ubuntu:focal ── ok: [docker_image_builder]
ubuntu:focal ──
ubuntu:focal ── TASK [Commit the container builder 'builder_ansible-playbook__ubuntu_focal' to create registry.stevedore.test/ubuntu:focal] ***
ubuntu:focal ── changed: [docker_image_builder]
ubuntu:focal ──
ubuntu:focal ── TASK [Remove container builder 'builder_ansible-playbook__ubuntu_focal'] *******
ubuntu:focal ── changed: [docker_image_builder]
ubuntu:focal ──
ubuntu:focal ── PLAY RECAP *********************************************************************
ubuntu:focal ── builder_ansible-playbook__ubuntu_focal : ok=6    changed=5    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
ubuntu:focal ── docker_image_builder       : ok=11   changed=4    unreachable=0    failed=0    skipped=2    rescued=0    ignored=0
ubuntu:focal ──
```

### Building the app1 image

You can use the previously created base image as `app1` parent image. That `app1` uses the `docker` driver on its image definition.
When building the `app1` image, the previously created base image (`registry.stevedore.test/app1:v1-ubuntufocal`) is used as the parent image. The `app1` image definition utilizes the `docker` driver demostrating the possibility to use both drives amongst the image definitions.

```sh
 [10-use-ansible-playbook-driver-example] Building Ubuntu image
registry.stevedore.test/app1:v1-ubuntufocal Step 1/5 : ARG image_from_fully_qualified_name
registry.stevedore.test/app1:v1-ubuntufocal Step 2/5 : FROM ${image_from_fully_qualified_name}
registry.stevedore.test/app1:v1-ubuntufocal ---> 3c9858f9d212
registry.stevedore.test/app1:v1-ubuntufocal Step 3/5 : COPY ./${app_name}/app.sh /app.sh
registry.stevedore.test/app1:v1-ubuntufocal ---> bea48d3b28e7
registry.stevedore.test/app1:v1-ubuntufocal Step 4/5 : CMD ["/app.sh"]
registry.stevedore.test/app1:v1-ubuntufocal ---> Running in 595474ea7df7
registry.stevedore.test/app1:v1-ubuntufocal ---> 8e568f4a5a78
registry.stevedore.test/app1:v1-ubuntufocal Step 5/5 : LABEL created_at=2023-06-16T05:42:54.880400037Z
registry.stevedore.test/app1:v1-ubuntufocal ---> Running in b53223959b25
registry.stevedore.test/app1:v1-ubuntufocal ---> a6859d67e096
registry.stevedore.test/app1:v1-ubuntufocal  ‣ sha256:a6859d67e09686a407973d70b660500bc62f2eff08473f724685c65d66aa837c
registry.stevedore.test/app1:v1-ubuntufocal [Warning] One or more build-args [image_from_name image_from_registry_host image_from_tag] were not consumed
registry.stevedore.test/app1:v1-ubuntufocal Successfully built a6859d67e096
registry.stevedore.test/app1:v1-ubuntufocal Successfully tagged registry.stevedore.test/app1:v1-ubuntufocal
```

### Cleaning the stack

```sh
Stopping the stack to run 10-use-ansible-playbook-driver-example

[+] Running 4/4
 ✔ Container 10-use-ansible-playbook-driver-example-stevedore-1   Removed                                                                                            2.5s
 ✔ Container 10-use-ansible-playbook-driver-example-registry-1    Removed                                                                                            0.2s
 ✔ Container 10-use-ansible-playbook-driver-example-dockerauth-1  Removed                                                                                            0.2s
 ✔ Network 10-use-ansible-playbook-driver-example_default         Removed                                                                                            0.3s
```
