# Envvars Credentials Store Example

The example demonstrates the utilization of environment variables as the [credentials store](https://gostevedore.github.io/docs/reference-guide/credentials/credentials-store/) in Stevedore. 
To accomplish this, the [Docker Compose file](./docker-compose.yml) includes the service `stevedore` with specific environment variables defined. The `STEVEDORE_CREDENTIALS_ENCRYPTION_KEY` environment variable is set with the encryption key, and the `STEVEDORE_ENVVARS_CREDENTIALS_82E99D42EE1191BB42FBFB444920104D` environment variable contains the credentials for the `registry.stevedore.test`.

- [Envvars Credentials Store Example](#envvars-credentials-store-example)
  - [Requirements](#requirements)
  - [Stack](#stack)
  - [Usage](#usage)
  - [Expected Output](#expected-output)
    - [Starting the stack](#starting-the-stack)
    - [Getting Credentials](#getting-credentials)
    - [Building images](#building-images)
    - [Cleaning the stack](#cleaning-the-stack)
  - [Additional information](#additional-information)
    - [Builders](#builders)
    - [Credentials](#credentials)


## Requirements
- Docker. _Tested on Docker server 20.10.21 and Docker API 1.41_
- Docker's Compose plugin or `docker-compose`. _Tested on Docker Compose version v2.17.3_
- `make` utility. _Tested on version 4.3-4.1build1_

## Stack
The stack required to run this example is defined in that [Docker Compose file](./docker-compose.yml). The stack consists of three services: a Docker Registry, a Docker Registry authorization and a Stevedore service. The Docker registry is used to store the Docker images built by Stevedore during the example execution. The Stevedore service is where the example is executed.

The Stevedore service is built from a container which is defined in that [Dockerfile](https://github.com/gostevedore/stevedore/blob/main/examples/03-envvars-credentials-store-example/stack/stevedore/Dockerfile).

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
Starting the stack to run 03-envvars-credentials-store-example

[+] Building 10.7s (21/21) FINISHED                                                                                                                              
 => [internal] load .dockerignore                                                                                                                           0.0s
 => => transferring context: 2B                                                                                                                             0.0s
 => [internal] load build definition from Dockerfile                                                                                                        0.0s
 => => transferring dockerfile: 989B                                                                                                                        0.0s
 => [internal] load metadata for docker.io/library/docker:20.10-dind                                                                                        1.3s
 => [internal] load metadata for docker.io/library/golang:1.19-alpine                                                                                       1.8s
 => [internal] load build context                                                                                                                           0.0s
 => => transferring context: 180.74kB                                                                                                                       0.0s
 => [stage-1 1/6] FROM docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                          0.0s
 => => resolve docker.io/library/docker:20.10-dind@sha256:af96c680a7e1f853ebdd50c1e0577e5df4089b033102546dd6417419564df3b5                                  0.0s
 => [golang 1/8] FROM docker.io/library/golang:1.19-alpine@sha256:4147c2ad24a347e53a8a361b663dcf82fca4157a3f9a5136d696a4e53dd22b65                          0.0s
 => => resolve docker.io/library/golang:1.19-alpine@sha256:4147c2ad24a347e53a8a361b663dcf82fca4157a3f9a5136d696a4e53dd22b65                                 0.0s
 => CACHED [golang 2/8] WORKDIR /usr/src/app                                                                                                                0.0s
 => CACHED [golang 3/8] RUN apk add --no-cache make build-base                                                                                              0.0s
 => CACHED [golang 4/8] COPY go.mod ./                                                                                                                      0.0s
 => CACHED [golang 5/8] COPY go.sum ./                                                                                                                      0.0s
 => CACHED [golang 6/8] RUN go mod download && go mod verify                                                                                                0.0s
 => [golang 7/8] COPY . ./                                                                                                                                  0.4s
 => [golang 8/8] RUN go build -ldflags "-s -w" -v -o /usr/local/bin/stevedore ./cmd/stevedore.go                                                            5.2s
 => CACHED [stage-1 2/6] COPY --from=golang /usr/local/go /usr/local/go                                                                                     0.0s
 => [stage-1 3/6] COPY --from=golang /usr/local/bin/stevedore /usr/local/bin/stevedore                                                                      0.1s
 => [stage-1 4/6] WORKDIR /go                                                                                                                               0.1s
 => [stage-1 5/6] RUN mkdir -p "/go/src" "/go/bin" && chmod -R 777 "/go"                                                                                    0.1s
 => [stage-1 6/6] COPY test/stack/client/entrypoint.sh /usr/local/bin/entrypoint.sh                                                                         0.0s
 => exporting to docker image format                                                                                                                        1.7s
 => => exporting layers                                                                                                                                     0.6s
 => => exporting manifest sha256:ac9e8575a110dc597bcdc64f39ec632c06177d57125a5f0bc7b0a0c40a5c7572                                                           0.0s
 => => exporting config sha256:5b5ffbbecd86d7924e315763076f9cf85bc4b91681adbc14586877a3c2ae4e8a                                                             0.0s
 => => sending tarball                                                                                                                                      1.1s
 => importing to docker                                                                                                                                     0.2s
[+] Running 4/4
 ✔ Network 03-envvars-credentials-store-example_default         Created                                                                                     0.1s 
 ✔ Container 03-envvars-credentials-store-example-stevedore-1   Started                                                                                     0.5s 
 ✔ Container 03-envvars-credentials-store-example-dockerauth-1  Started                                                                                     0.5s 
 ✔ Container 03-envvars-credentials-store-example-registry-1    Started                                                                                     0.8s 
```

### Getting Credentials
To obtain the credentials, use the command: `stevedore get credentials --show-secrets`.

```sh
 Run example 03-envvars-credentials-store-example

 [03-envvars-credentials-store-example] Get credentials
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
/certs/server/cert.pem: OK
 Waiting for dockerd to be ready...
/certs/client/cert.pem: OK
ID                      TYPE              CREDENTIALS
registry.stevedore.test username-password username=admin, password=admin
```

### Building images
The example uses the command `stevedore build my-app --image-version 3.2.1 --push-after-build` to build and automatically promote the images to the Docker registry.

```sh
 [03-envvars-credentials-store-example] Build my-app and push images after build
 Waiting for dockerd to be ready...
 Waiting for dockerd to be ready...
/certs/server/cert.pem: OK
 Waiting for dockerd to be ready...
/certs/client/cert.pem: OK
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 1/7 : ARG image_from_name
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 2/7 : ARG image_from_tag
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 1/7 : ARG image_from_name
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 3/7 : FROM ${image_from_name}:${image_from_tag}
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 2/7 : ARG image_from_tag
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 3/7 : FROM ${image_from_name}:${image_from_tag}
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  1.35:  Pulling from library/busybox 
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  c15cbdab5f8e:  Pull complete 
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  Digest: sha256:b4e4a06de46acc0958cd93e2eeb769077d255f06a7c3a91196509c16b7bc989e 
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  Status: Downloaded newer image for busybox:1.35 
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> dddc7578369a
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 4/7 : ARG whoami=unknown
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> Running in 6508c0d9f02e
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> 5f23aaccbaba
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 5/7 : RUN echo "Hey there, I'm ${whoami}!" > /whoami.txt
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> Running in 7ed0e6569811
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> 33a9c74f3276
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 6/7 : CMD ["cat","/whoami.txt"]
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> Running in 2fb0eb6f5d7e
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> d7db3e3ccaa5
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 7/7 : LABEL created_at=2023-05-17T16:44:17.242832601Z
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> Running in 98af2bff465e
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> d0a3380b6bb7
registry.stevedore.test/my-app:3.2.1-busybox1.35  ‣ sha256:d0a3380b6bb7a960abc933678722ddd9358dc7c59a431c23ce0188fa30d7605f
registry.stevedore.test/my-app:3.2.1-busybox1.35 Successfully built d0a3380b6bb7
registry.stevedore.test/my-app:3.2.1-busybox1.35 Successfully tagged registry.stevedore.test/my-app:3.2.1-busybox1.35
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  The push refers to repository [registry.stevedore.test/my-app] 
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  1.36:  Pulling from library/busybox 
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  a58ecd4f0c86:  Pull complete 
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  Digest: sha256:9e2bbca079387d7965c3a9cee6d0c53f4f4e63ff7637877a83c4c05f2a666112 
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  Status: Downloaded newer image for busybox:1.36 
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> af2c3e96bcf1
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  cab22fd6c10e:  Pushed 
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  42ef21f45b9a:  Pushing [========================>                          ]  2.406MB/4.855MB
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> Running in d4fc0cafba00
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> 59c14dad7f27
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 5/7 : RUN echo "Hey there, I'm ${whoami}!" > /whoami.txt
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  cab22fd6c10e:  Pushed 
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  42ef21f45b9a:  Pushed 
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  3.2.1-busybox1.35: digest: sha256:584ec9691bc9cc733b802ae0146f9f69d7f7f4cd56e1c50819d81d9dda10a113 size: 735 
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> a80b636783f3
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 6/7 : CMD ["cat","/whoami.txt"]
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> Running in 5fe148b33fd9
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> 48ab053ea782
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 7/7 : LABEL created_at=2023-05-17T16:44:17.242594057Z
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> Running in 2a4621194d75
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> c50b9be6bcdf
registry.stevedore.test/my-app:3.2.1-busybox1.36  ‣ sha256:c50b9be6bcdf23ea4444e6f6c61d2780c893800cfe3259853692d07edbd0aa98
registry.stevedore.test/my-app:3.2.1-busybox1.36 Successfully built c50b9be6bcdf
registry.stevedore.test/my-app:3.2.1-busybox1.36 Successfully tagged registry.stevedore.test/my-app:3.2.1-busybox1.36
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  The push refers to repository [registry.stevedore.test/my-app] 
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  82818cd94998:  Pushed 
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  1f1d08b81bbe:  Pushed 
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  3.2.1-busybox1.36: digest: sha256:b2e3fd5dd2ba05b3b4c395051cd2828c14bd8346c312f4a19604e2ee03e6b18b size: 735 
```

### Cleaning the stack
```sh
Stopping the stack to run 03-envvars-credentials-store-example

[+] Running 6/6
 ✔ Container 03-envvars-credentials-store-example-stevedore-1                 Removed                                                                      10.3s 
 ✔ Container 03-envvars-credentials-store-example-stevedore-run-7eaf9c4dc1d2  Removed                                                                       0.0s 
 ✔ Container 03-envvars-credentials-store-example-registry-1                  Removed                                                                       0.3s 
 ✔ Container 03-envvars-credentials-store-example-stevedore-run-90197b4f5fb0  Removed                                                                       0.0s 
 ✔ Container 03-envvars-credentials-store-example-dockerauth-1                Removed                                                                       0.2s 
 ✔ Network 03-envvars-credentials-store-example_default                       Removed                                                                       0.5s 
```

## Additional information
In addition to the core steps outlined in the example, the following section provides additional information and insights to further enhance your understanding of how this example uses Stevedore.

### Builders
In this example, the image definitions utilize an [inline builder](https://gostevedore.github.io/docs/reference-guide/builder/#in-line-builder) instead of a [global builder](https://gostevedore.github.io/docs/reference-guide/builder/#global-builder). This is achieved by removing the `builders_path` parameter from the Stevedore configuration, within the [./stevedore.yaml](stevedore.yaml) file, as well as removing the `./builders` folder.

### Credentials
This example utilizes environment variables as the [credentials store](https://gostevedore.github.io/docs/reference-guide/credentials/credentials-store/#envvars-storage). Stevedore recognizes this as the `envvars` storage type. To use this store, the credentials section in the [./stevedore.yaml](stevedore.yaml) file should be configured as follows. However, you can also configure it by setting the `STEVEDORE_CREDENTIALS_STORAGE_TYPE` environment variable.
```yaml
credentials:
  storage_type: envvars
```

The `envvars` store requires an encryption key, which in this example is set through the environment variable `STEVEDORE_CREDENTIALS_ENCRYPTION_KEY` instead of defining it in the [./stevedore.yaml](stevedore.yaml) file. It's important to note that any [Stevedore's configuration parameter can be set using environment variables](https://gostevedore.github.io/docs/getting-started/configuration/#configuration-from-environment-variables).

The `stevedore create credentials` subcommand is used to create the credentials. When using the `envvars` storage, it does not create the environment variable that contains the credentials. Instead, it returns the environment variable name and its value, which you need to set in your system.
```
/app/examples/03-envvars-credentials-store-example # stevedore create credentials registry.stevedore.test --username admin
Password: 
 You must create the following environment variable to use the recently created credentials: 
  STEVEDORE_ENVVARS_CREDENTIALS_82E99D42EE1191BB42FBFB444920104D=adee358c1be79793bae7328c750f29ec4fe77dcaeefc1807b969bdfd086f1db2051bab4225b5bfcc2c8f6dd88ad8d88abe0f04d70959c71a27c6e40e701307f9bfc4ba9120697c9d4162b94620c3ad5c7ec22a0314f6aecbd768cce71ba4ef1f558b29f9ec11d8dce96a30004647155e454807f730c0abb0d17399025c20ca20bb589071d7806a3879153a2c48e72ff01cc26f50dabc2855350fd27fd483ecf31356c4a28ec3e6c869d3b92a554f02e9dd510744af8319f428df95a34835fa0aa8ac2ccef2d4572d4099f20f7c2a640c68f2e0ff8edda84303caab1c7456c482510b0d2766cdd57351e3619083c75007c2f64d9c27455c7198803c397d188c7d6e8ce1a96451e39b60c25dd8d77647c0ffeb94cf074c6d1d9815669c35a3f2e2ebd227b316b843a07b3ad4114816bd4a5d43f101d7aa98e085d7046572404701ce8d4f0ab4d01177d27dac8dd99feeb3692fa2588d8a120e089f561eaeb732dd0d26c408ba2c11a4c42f2e2e171a74dc8815fd90c3619ff7e58c602b2c077c1231acffe542e2a576b3e9e43138fd68a8cfd2ddc96b7beba371ef4c5fc615bbdfebda1d01d2b9df32c354d73c55aa6764257e9d480188127407015f6bd62f768f8af6531723113cc71504ede4175f69c2e63e6080f216a109ff58bf7b09ed45263caabf30fd4857d19bf8c678de7c5b1664b7686e675702b5c46b9329ff53f8f2abf2b0 
```
