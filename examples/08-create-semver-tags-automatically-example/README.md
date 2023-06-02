# Create SemVer Tags Automatically

This example showcases how to use the SemVer specification to create automatically new image tags for a Docker image.

- [Create SemVer Tags Automatically](#create-semver-tags-automatically)
  - [Requirements](#requirements)
  - [Stack](#stack)
  - [Usage](#usage)
  - [Example Execution Insights](#example-execution-insights)
    - [Starting the Stack](#starting-the-stack)
    - [Waiting for Dockerd to be Ready](#waiting-for-dockerd-to-be-ready)
    - [Getting Images](#getting-images)
    - [Building images](#building-images)
    - [Cleaning the stack](#cleaning-the-stack)
  - [Additional Information](#additional-information)
    - [Stevedore Configuration](#stevedore-configuration)

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
 [08-create-semver-tags-automatically-example] Create SSH keys
[+] Building 0.0s (0/0)
[+] Creating 2/1
 ✔ Network 08-create-semver-tags-automatically-example_default  Created                                                                                              0.1s
 ✔ Volume "08-create-semver-tags-automatically-example_ssh"     Created                                                                                              0.0s
[+] Building 0.0s (0/0)
```

Once the SSH key pair is generated, the services are started.
```sh
Starting the stack to run 08-create-semver-tags-automatically-example

[+] Building 7.1s (37/37) FINISHED
 => [stevedore internal] load .dockerignore                                                                                                                          0.0s
 => => transferring context: 2B                                                                                                                                      0.0s
 => [stevedore internal] load build definition from Dockerfile                                                                                                       0.0s
 => => transferring dockerfile: 912B                                                                                                                                 0.0s
 => [stevedore internal] load metadata for docker.io/library/docker:20.10-dind                                                                                       0.4s
 => [stevedore internal] load metadata for docker.io/library/golang:1.19-alpine                                                                                      0.5s
 => [ssh-keygen internal] load .dockerignore                                                                                                                         0.0s
 => => transferring context: 2B                                                                                                                                      0.0s
 => [ssh-keygen internal] load build definition from Dockerfile                                                                                                      0.0s
 => => transferring dockerfile: 129B                                                                                                                                 0.0s
 => [gitserver internal] load .dockerignore                                                                                                                          0.0s
 => => transferring context: 2B                                                                                                                                      0.0s
 => [gitserver internal] load build definition from Dockerfile                                                                                                       0.0s
 => => transferring dockerfile: 1.41kB                                                                                                                               0.0s
 => [ssh-keygen internal] load metadata for docker.io/library/ubuntu:22.04                                                                                           0.4s
 => [gitserver internal] load metadata for docker.io/library/alpine:3.16                                                                                             0.4s
 => [ssh-keygen 1/2] FROM docker.io/library/ubuntu:22.04@sha256:dfd64a3b4296d8c9b62aa3309984f8620b98d87e47492599ee20739e8eb54fbf                                     0.0s
 => => resolve docker.io/library/ubuntu:22.04@sha256:dfd64a3b4296d8c9b62aa3309984f8620b98d87e47492599ee20739e8eb54fbf                                                0.0s
 => CACHED [ssh-keygen 2/2] RUN apt-get update     && apt-get install -y         openssh-client                                                                      0.0s
 => [ssh-keygen] exporting to docker image format                                                                                                                    0.3s
 => => exporting layers                                                                                                                                              0.0s
 => => exporting manifest sha256:b9dd944e51f68f2770d0f057a19a0fd3e50225c4cdf763e8dba0c8db364f5bae                                                                    0.0s
 => => exporting config sha256:b958b17ee04489b3cf462e68c3060e302055ab5dffac7d2469a46e5c519c91c2                                                                      0.0s
 => => sending tarball                                                                                                                                               0.3s
 => [gitserver 1/6] FROM docker.io/library/alpine:3.16@sha256:c2b622f6e510a0d25bccaffa9e67b75a6860cb09b74bb58cfc36a9ef4331109f                                       0.0s
 => => resolve docker.io/library/alpine:3.16@sha256:c2b622f6e510a0d25bccaffa9e67b75a6860cb09b74bb58cfc36a9ef4331109f                                                 0.0s
 => [gitserver internal] load build context                                                                                                                          0.0s
 => => transferring context: 96B                                                                                                                                     0.0s
 => CACHED [gitserver 2/6] WORKDIR /git                                                                                                                              0.0s
 => CACHED [gitserver 3/6] RUN apk add --no-cache         openssh         git     && rm -rf /var/cache/apk/*     && ssh-keygen -A     && adduser -D --home /home/gi  0.0s
 => CACHED [gitserver 4/6] RUN apk add  --no-cache         nginx         git-daemon         fcgiwrap         spawn-fcgi     && rm -rf /var/cache/apk/*               0.0s
 => CACHED [gitserver 5/6] COPY files/nginx.conf /etc/nginx/nginx.conf                                                                                               0.0s
 => CACHED [gitserver 6/6] COPY entrypoint.sh /entrypoint.sh                                                                                                         0.0s
 => [gitserver] exporting to docker image format                                                                                                                     0.1s
 => => exporting layers                                                                                                                                              0.0s
 => => exporting manifest sha256:f80426a553ce2ec2c26275b1bd00213a22e87f77b60c31643d6b7d43b0c6d215                                                                    0.0s
 => => exporting config sha256:5d432c89f104f66d94d76cd2eed4618b6236f3fdd318186acafdfb5f739e3c5d                                                                      0.0s
 => => sending tarball                                                                                                                                               0.1s
 => [stevedore golang 1/7] FROM docker.io/library/golang:1.19-alpine@sha256:4147c2ad24a347e53a8a361b663dcf82fca4157a3f9a5136d696a4e53dd22b65                         0.0s
 => => resolve docker.io/library/golang:1.19-alpine@sha256:4147c2ad24a347e53a8a361b663dcf82fca4157a3f9a5136d696a4e53dd22b65                                          0.0s
 => [stevedore stage-1 1/4] FROM docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                         0.0s
 => => resolve docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                                           0.0s
 => [stevedore internal] load build context                                                                                                                          0.1s
 => => transferring context: 310.17kB                                                                                                                                0.1s
 => [gitserver gitserver] importing to docker                                                                                                                        0.0s
 => CACHED [stevedore golang 2/7] WORKDIR /usr/src/app                                                                                                               0.0s
 => CACHED [stevedore golang 3/7] RUN apk add --no-cache make build-base                                                                                             0.0s
 => CACHED [stevedore golang 4/7] COPY go.mod go.sum ./                                                                                                              0.0s
 => CACHED [stevedore golang 5/7] RUN go mod download                                                                                                                0.0s
 => [stevedore golang 6/7] COPY . ./                                                                                                                                 0.5s
 => [ssh-keygen ssh-keygen] importing to docker                                                                                                                      0.0s
 => [stevedore golang 7/7] RUN go build -ldflags "-s -w" -v -o /usr/local/bin/stevedore ./cmd/stevedore.go                                                           5.3s
 => CACHED [stevedore stage-1 2/4] COPY --from=golang /usr/local/bin/stevedore /usr/local/bin/stevedore                                                              0.0s
 => CACHED [stevedore stage-1 3/4] COPY examples/07-inject-dockerfile-example/stack/stevedore/entrypoint.sh /usr/local/bin/entrypoint.sh                             0.0s
 => CACHED [stevedore stage-1 4/4] COPY examples/07-inject-dockerfile-example/stack/stevedore/wait-for-dockerd.sh /usr/local/bin/wait-for-dockerd.sh                 0.0s
 => [stevedore] exporting to docker image format                                                                                                                     0.5s
 => => exporting layers                                                                                                                                              0.0s
 => => exporting manifest sha256:be551e1fef67432b7afc40d7c73ce737046d04bcee90001919c0ace477f8d0ef                                                                    0.0s
 => => exporting config sha256:94301df09c29522e38473549fb0eabc05f11004164cf475d57f858fca10e4da9                                                                      0.0s
 => => sending tarball                                                                                                                                               0.5s
 => [stevedore stevedore] importing to docker                                                                                                                        0.0s
[+] Running 6/6
 ✔ Container 08-create-semver-tags-automatically-example-ssh-keygen-1  Started                                                                                       1.0s
 ✔ Container 08-create-semver-tags-automatically-example-stevedore-1   Started                                                                                       1.1s
 ✔ Container 08-create-semver-tags-automatically-example-worker-1      Started                                                                                       0.4s
 ✔ Container 08-create-semver-tags-automatically-example-dockerauth-1  Started                                                                                       1.1s
 ✔ Container 08-create-semver-tags-automatically-example-gitserver-1   Started                                                                                       1.0s
 ✔ Container 08-create-semver-tags-automatically-example-registry-1    Started                                                                                       1.3s
```

And finally, it generates a `known_hosts` file based on the SSH keys.
```sh
 [08-create-semver-tags-automatically-example] Create known_hosts file
# gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
# gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
# gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
# gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
# gitserver.stevedore.test:22 SSH-2.0-OpenSSH_9.0
```

### Waiting for Dockerd to be Ready
Before starting the execution of the Stevedore command, it is important to ensure that the Docker daemon (dockerd) is ready. The stevedore service Docker image includes a script, [wait-for-dockerd.sh]((./stack/stevedore/wait-for-dockerd.sh)), which can be used to ensure the readiness of the Docker daemon.
```sh
 Run example 08-create-semver-tags-automatically-example

 [08-create-semver-tags-automatically-example] Waiting for dockerd
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
```

### Getting Images
To view the images in tree format, run `stevedore get images --tree`.

```sh
 [08-create-semver-tags-automatically-example] Get images
├─── busybox:1.35
│  ├─── registry.stevedore.test/base:busybox-1.35
│  │  ├─── registry.stevedore.test/app1:0.1.2
│  │  ├─── registry.stevedore.test/app1:{{ .Version }}
├─── busybox:1.36
│  ├─── registry.stevedore.test/base:busybox-1.36
```

### Building images
Before diving into the process of automatically generating tags based on semantic versioning, it is important to first create the `base` image that the `app1` Docker image will be built from.
```sh
 [08-create-semver-tags-automatically-example] Build the base image, and push the images after build
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
registry.stevedore.test/base:busybox-1.35 ‣  1.35:  Pulling from library/busybox
registry.stevedore.test/base:busybox-1.35 ‣  c15cbdab5f8e:  Pull complete
registry.stevedore.test/base:busybox-1.35 ‣  Digest: sha256:b4e4a06de46acc0958cd93e2eeb769077d255f06a7c3a91196509c16b7bc989e
registry.stevedore.test/base:busybox-1.35 ‣  Status: Downloaded newer image for busybox:1.35
registry.stevedore.test/base:busybox-1.35 ---> dddc7578369a
registry.stevedore.test/base:busybox-1.35 Step 4/7 : RUN echo "anonymous:x:10001:10001:,,,:/app:/bin/sh" >> /etc/passwd &&     echo "anonymous:x:10001:" >> /etc/group &&     mkdir -p /app &&     chown 10001:10001 /app
registry.stevedore.test/base:busybox-1.35 ---> Running in 60216c71c17d
registry.stevedore.test/base:busybox-1.36 ---> 48307f5cc0d4
registry.stevedore.test/base:busybox-1.36 Step 5/7 : USER anonymous
registry.stevedore.test/base:busybox-1.36 ---> Running in 74c192cba8bd
registry.stevedore.test/base:busybox-1.36 ---> 12e6b2a99bb8
registry.stevedore.test/base:busybox-1.36 Step 6/7 : WORKDIR /app
registry.stevedore.test/base:busybox-1.36 ---> Running in 556cce08f5b6
registry.stevedore.test/base:busybox-1.36 ---> 895c2c66f5c7
registry.stevedore.test/base:busybox-1.36 Step 7/7 : LABEL created_at=2023-05-31T17:52:11.628612838Z
registry.stevedore.test/base:busybox-1.36 ---> Running in f9d74b6afc55
registry.stevedore.test/base:busybox-1.35 ---> 51181b15abed
registry.stevedore.test/base:busybox-1.35 Step 5/7 : USER anonymous
registry.stevedore.test/base:busybox-1.35 ---> Running in 7c85e0284374
registry.stevedore.test/base:busybox-1.36 ---> b58ce6aa1d33
registry.stevedore.test/base:busybox-1.36  ‣ sha256:b58ce6aa1d3325c983bc62b039bfea89c955cd9dd87f9411a929fa2731bfcb98
registry.stevedore.test/base:busybox-1.36 Successfully built b58ce6aa1d33
registry.stevedore.test/base:busybox-1.36 Successfully tagged registry.stevedore.test/base:busybox-1.36
registry.stevedore.test/base:busybox-1.36 ‣  The push refers to repository [registry.stevedore.test/base]
registry.stevedore.test/base:busybox-1.36 ‣  be4564931709:  Preparing
registry.stevedore.test/base:busybox-1.36 ‣  9547b4c33213:  Preparing
registry.stevedore.test/base:busybox-1.35 ---> 9939341d5b92
registry.stevedore.test/base:busybox-1.35 Step 6/7 : WORKDIR /app
registry.stevedore.test/base:busybox-1.36 ‣  be4564931709:  Pushed
registry.stevedore.test/base:busybox-1.36 ‣  9547b4c33213:  Pushing [>                                                  ]  66.56kB/4.863MB
registry.stevedore.test/base:busybox-1.35 ---> 050c5096cd69
registry.stevedore.test/base:busybox-1.35 Step 7/7 : LABEL created_at=2023-05-31T17:52:11.629677004Z
registry.stevedore.test/base:busybox-1.35 ---> Running in f982b5e64265
registry.stevedore.test/base:busybox-1.35 ---> e358f7a01612
registry.stevedore.test/base:busybox-1.35  ‣ sha256:e358f7a016128310289f69996d2f61937ab273da02c1ad06b61e91c34fccba6d
registry.stevedore.test/base:busybox-1.35 Successfully built e358f7a01612
registry.stevedore.test/base:busybox-1.35 Successfully tagged registry.stevedore.test/base:busybox-1.35
registry.stevedore.test/base:busybox-1.35 ‣  The push refers to repository [registry.stevedore.test/base]
registry.stevedore.test/base:busybox-1.36 ‣  be4564931709:  Pushed
registry.stevedore.test/base:busybox-1.35 ‣  be4564931709:  Layer already exists
registry.stevedore.test/base:busybox-1.35 ‣  42ef21f45b9a:  Pushed
registry.stevedore.test/base:busybox-1.35 ‣  busybox-1.35: digest: sha256:c33a4f1d388f16422733b768ce2f2782e2f8821f65093e2e68e40a484a229ab6 size: 735
```

Once the base image is prepared, the process of building the `app1` application can begin. The following command is used for building the application: `stevedore  build app1 --enable-semver-tags --image-version 0.1.2-rc1+$(date -u +"%a%d%m%Y%H%M")`.

It is important to note that the version specified in the command includes the use of the `date` command, making it dynamic. This allows for the inclusion of the current date and time in the version value. In the [image definition](./images/applications.yaml), a wildcard version is used to accommodate dynamic versioning.
```sh
 [08-create-semver-tags-automatically-example] Build the app1 image, and push the images after build
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 Step 1/8 : ARG image_from_name
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 Step 2/8 : ARG image_from_tag
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 Step 3/8 : ARG image_from_registry_host
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 Step 4/8 : FROM ${image_from_registry_host}/${image_from_name}:${image_from_tag}
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ---> e358f7a01612
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 Step 5/8 : ARG app_name
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ---> Running in 5b53a0d2c598
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ---> 4c97919950ce
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 Step 6/8 : COPY ./${app_name}/app.sh /app.sh
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ---> 261deb2b97dc
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 Step 7/8 : CMD ["/app.sh"]
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ---> Running in ed72352b60e1
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ---> f6c2ce400f23
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 Step 8/8 : LABEL created_at=2023-05-31T17:52:16.305237989Z
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ---> Running in e8eae7c1c5a3
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ---> 154edb892bde
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752  ‣ sha256:154edb892bde17db916f54f7f427c739da824599b7120107fb621d68723dd455
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 Successfully built 154edb892bde
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 Successfully tagged registry.stevedore.test/app1:0.1.2
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 Successfully tagged registry.stevedore.test/app1:0.1
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 Successfully tagged registry.stevedore.test/app1:0
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 Successfully tagged registry.stevedore.test/app1:0-PRERELEASE-rc1
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 Successfully tagged registry.stevedore.test/app1:0.1-PRERELEASE-rc1-BUILD-Wed310520231752
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 Successfully tagged registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  The push refers to repository [registry.stevedore.test/app1]
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  857d0ed80b00:  Pushed
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  be4564931709:  Mounted from base
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  42ef21f45b9a:  Mounted from base
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  0.1.2: digest: sha256:9b4d6dabfe6897be9270916fd893fd7e1fd731a07dd59c3edc3cd868d8b80856 size: 942
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  The push refers to repository [registry.stevedore.test/app1]
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  857d0ed80b00:  Layer already exists
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  be4564931709:  Layer already exists
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  42ef21f45b9a:  Layer already exists
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  0.1: digest: sha256:9b4d6dabfe6897be9270916fd893fd7e1fd731a07dd59c3edc3cd868d8b80856 size: 942
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  The push refers to repository [registry.stevedore.test/app1]
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  857d0ed80b00:  Layer already exists
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  be4564931709:  Layer already exists
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  42ef21f45b9a:  Layer already exists
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  0: digest: sha256:9b4d6dabfe6897be9270916fd893fd7e1fd731a07dd59c3edc3cd868d8b80856 size: 942
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  The push refers to repository [registry.stevedore.test/app1]
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  857d0ed80b00:  Layer already exists
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  be4564931709:  Layer already exists
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  42ef21f45b9a:  Layer already exists
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  0-PRERELEASE-rc1: digest: sha256:9b4d6dabfe6897be9270916fd893fd7e1fd731a07dd59c3edc3cd868d8b80856 size: 942
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  The push refers to repository [registry.stevedore.test/app1]
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  857d0ed80b00:  Layer already exists
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  be4564931709:  Layer already exists
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  42ef21f45b9a:  Layer already exists
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  0.1-PRERELEASE-rc1-BUILD-Wed310520231752: digest: sha256:9b4d6dabfe6897be9270916fd893fd7e1fd731a07dd59c3edc3cd868d8b80856 size: 942
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  The push refers to repository [registry.stevedore.test/app1]
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  857d0ed80b00:  Layer already exists
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  be4564931709:  Layer already exists
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  42ef21f45b9a:  Layer already exists
registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752 ‣  0.1.2-rc1_Wed310520231752: digest: sha256:9b4d6dabfe6897be9270916fd893fd7e1fd731a07dd59c3edc3cd868d8b80856 size: 942
```

### Cleaning the stack
```sh
Stopping the stack to run 08-create-semver-tags-automatically-example

[+] Running 8/8
 ✔ Container 08-create-semver-tags-automatically-example-worker-1      Removed                                                                                       0.0s
 ✔ Container 08-create-semver-tags-automatically-example-registry-1    Removed                                                                                       0.2s
 ✔ Container 08-create-semver-tags-automatically-example-stevedore-1   Removed                                                                                       3.4s
 ✔ Container 08-create-semver-tags-automatically-example-gitserver-1   Removed                                                                                       3.5s
 ✔ Container 08-create-semver-tags-automatically-example-ssh-keygen-1  Removed                                                                                       3.3s
 ✔ Container 08-create-semver-tags-automatically-example-dockerauth-1  Removed                                                                                       0.2s
 ✔ Volume 08-create-semver-tags-automatically-example_ssh              Removed                                                                                       0.0s
 ✔ Network 08-create-semver-tags-automatically-example_default         Removed                                                                                       0.5s
```

## Additional Information

### Stevedore Configuration
When you tag a Docker image with a semantic version compliant value, Stevedore offers the capability to automatically generate additional tags based on semantic versioning. This can be achieved by using the `--enable-semver-tags` flag in either the [build](https://gostevedore.github.io/docs/reference-guide/cli/#build) or [promote](https://gostevedore.github.io/docs/reference-guide/cli/#promote) subcommands.

To define the additional tags to be generated, you can configure a list of templates in the [semantic_version_tags_templates](https://gostevedore.github.io/docs/getting-started/configuration/#semantic_version_tags_templates) parameter of the Stevedore configuration.

Stevedore uses Golang's [text/template](https://pkg.go.dev/text/template) to define the additional tags.
In this example, the following templates are used:
```yaml
semantic_version_tags_templates:
- "{{ .Major }}.{{ .Minor }}.{{ .Patch }}"
- "{{ .Major }}.{{ .Minor }}"
- "{{ .Major }}"
- "{{ .Major }}{{ with .PreRelease }}-PRERELEASE-{{ . }}{{ end }}"
- "{{ .Major }}.{{ .Minor }}{{ with .PreRelease }}-PRERELEASE-{{ . }}{{ end }}{{ with .Build }}-BUILD-{{ . }}{{ end }}"
```

These templates result in the generation of the following tags for the `app1` Docker image:
- `registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752` 
- `registry.stevedore.test/app1:0.1.2` 
- `registry.stevedore.test/app1:0.1` 
- `registry.stevedore.test/app1:0` 
- `registry.stevedore.test/app1:0-PRERELEASE-rc1` 
- `registry.stevedore.test/app1:0.1-PRERELEASE-rc1-BUILD-Wed310520231752` 
- `registry.stevedore.test/app1:0.1.2-rc1_Wed310520231752` 

