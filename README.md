[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

![stevedore-logo](docs/logo/logo_4_stevedore.png "Stevedore logo")

Stevedore is a tool that allows you to build Docker images at scale, a factory to build Docker images. It is not an alternative to Dockerfile or Buildkit, but a way to use them to improve the experience you have building and promoting bunches of Docker images

## [Stevedore website](https://gostevedore.github.io/) - [Documentation](https://gostevedore.github.io/documentation/)

- [Stevedore website - Documentation](#stevedore-website---documentation)
- [Why stevedore?](#why-stevedore)
- [Getting started guide](#getting-started-guide)
  - [Installation](#installation)
  - [Initial setup](#initial-setup)
  - [Build the Docker images for an application](#build-the-docker-images-for-an-application)
  - [Promote the images to Docker registry](#promote-the-images-to-docker-registry)


## Why stevedore?
Stevedore is a helpful tool when you need to **build bunches of Docker images** and build them in a **standardized way**. It lets you define how to build your Docker images as well as the relationship amongst those images. Having that relationship, you can **build automatically the children's images** when the parent image is ready.

To build Docker images, you can use **multiple sources as build context**, for example, building images using a folder from your local host or from a git repository. And not only that, but you can also merge several sources to create the ultimate build context.

You can **generate automatically semver tags** following the version tree when you use the semantic version on your tags.

Finally, Stevedore provides a **credentials store** to authenticate on your behalf to a Docker registry when you want to push or **promote Docker images**. It also provides the capability to authenticate to AWS ECR.

## Getting started guide

### Installation

You can use the script provided on the Stevedore's repository to install it.
```sh
curl -fsSL https://raw.githubusercontent.com/gostevedore/stevedore/main/scripts/install.sh | sudo bash -
```
Visit the [documentation](https://gostevedore.github.io/documentation/getting-started/quick-start/#download-and-install-stevedore) to see other installation methods.

### Initial setup
- Create the folders structure to define images and builders.
```sh
/ $ mkdir -p /docker/images /docker/builders
```

- Initialize Stevedore. We create a configuration file for this guide, but it can be also defined on environment variables.
```sh
/ $ cd /docker
/docker $ stevedore initialize --builders-path builders --credentials-storage-type local --generate-credentials-encryption-key --images-path images --log-path-file /dev/null
2023-02-03 21:57:20     INFO    Executing command 'stevedore [COMMAND] [OPTIONS] initialize'
```

You can validate the configuration using the get configuration subcommand.
```sh
/docker $ stevedore get configuration

 builders_path: builders
 concurrency: 4
 semantic_version_tags_enabled: false
 images_path: images
 log_path: /dev/null
 push_images: false
 semantic_version_tags_templates:
   - {{ .Major }}.{{ .Minor }}.{{ .Patch }}
 credentials:
   storage_type: local
   format: json
   encryption_key: 1c591ac2d9c2664db265704052c17a67
```

- Create credentials to log in to the Docker registry. Since you create basic auth credentials, it means using a username and password, you are asked to introduce the password.
```sh
/docker $ stevedore create credentials registry.stevedore.test --username admin
Password:
```

Validate that the credentials has been stored.
```sh
/docker $ stevedore get credentials
ID                      TYPE              CRENDENTIALS
registry.stevedore.test username-password username=admin
```

Since Stevedore's configuration provides an encryption key, the credentials content is encrypted at rest.
```sh
/docker $ cat credentials/82e99d42ee1191bb42fbfb444920104d
2d067af765e8de39fd76b5cd1a430768a66464208bf8fe0c092bcb146a24feed377b916952494e6cea55187f397aee8c2d90a55ef9882cf85ee97ed5660700afa002767e028b4ea6bde274a524e7100f5729601f2d44caa08a1a102af7f79079723f35953133be56e31d0eaf44f52255a5c94512c74625dca3f00b77f9031d4f48f5cf38293f7a2a90f727c9a5eedf57e001ea8766a6d1e47d20354ad3ca6cf022ee70b97b3598e7377355a7b52f62fab8b6628b230c33cf0234ea1208c0d6ecef65e8e1206e7daf15acbbfb62d77650982c9f487129534b367a7fc2b519fd04538bffe87e57184adabc57e613be4b0e106480f9078c8c1c916b65de0039a3adc6cea70b962c0d93477e114e25a2a160db0218e9312d0df00d9b1044e0b4981982834094a36d8d9e4fd9ed766605c1a0b43ff11219d6e5bebb414e084102ce8cd57e86a658455de5fff927ba9be039be4afd393b7e8137cdefb268ed7c79ff37a60c1be3693372d7890b9f8f7c81fa559719ff83371e080efaee86b239e127a094136d056641a4a58245f71414b53dd96214149c6f54de055064163fa9dcb0a6b274d9269a49e1e4f76a9c0a89ec12d40c577ea1b5c1b3f9f8454f5473518c58e8c139a96335850a6415df7e2d41f5ab383f85d30281fd29493db0f0fb577aba35326cd5063b463a41161888514dbce09ccea5887494733256d74341e2623a5328c656/docker #
```

### Build the Docker images for an application
The idea is to create an application that has multiple versions, and must also be built from multiple parents.

First of all we create the application.
Take a look at the Dockerfile arguments `image_from_name` and `image_from_tag`. These arguments are automatically generated by Stevedore to provide information about the parent. However, you can also use any other argument set in the image definition through the `vars` or `persistent_vars` attributes.

```sh
/ $ mkdir -p /apps/my-app
/ $ cd /apps/my-app
/apps/my-app $ cat << EOF > Dockerfile
ARG image_from_name
ARG image_from_tag

FROM \${image_from_name}:\${image_from_tag}

CMD ["echo","Hey there!"]
EOF
```

- Define the builders. Since we want to define multiple image versions for the application, we define a single builder to standardize the Docker build.
```sh
/ $ cd /docker/builders
/docker/builders $ cat << EOF > apps.yaml
builders:
  my-app:
    driver: docker
    options:
      context:
        - path: /apps/my-app
EOF
```

The builder is already defined.
```sh
/ $ cd /docker
/docker $ stevedore get builders
NAME    DRIVER
my-apps docker
```

- Define the foundational images, the parent images.
On the parent images we define the `persistent labels`. Labels defined here are also set to all the images created from these parent images.
```sh
/ $ cd /docker/images
/docker/images $ cat << EOF > foundational.yaml
images:
  busybox:
    "1.36":
      persistent_labels:
        created_at: "{{ .DateRFC3339Nano }}"
    "1.35":
      persistent_labels:
        created_at: "{{ .DateRFC3339Nano }}"
EOF
```

- Define the images for our application.
```sh
/docker/images $ cat << EOF > applications.yaml
images:
  my-app:
    "2.1.0":
      name: "{{ .Name }}"
      version: "{{ .Version }}-{{ .Parent.Name }}{{ .Parent.Version }}"
      registry: registry.stevedore.test
      builder: my-app
      parents:
        busybox:
          - "1.35"
    "3.2.1":
      name: "{{ .Name }}"
      version: "{{ .Version }}-{{ .Parent.Name }}{{ .Parent.Version }}"
      registry: registry.stevedore.test
      builder: my-app
      parents:
        busybox:
          - "1.35"
          - "1.36"
EOF
```

We ensure that images are already defined.
```sh
/ $ cd /docker
/docker $ stevedore get images --tree
├─── busybox:1.35
│  ├─── registry.stevedore.test/my-app:3.2.1-busybox1.35
│  ├─── registry.stevedore.test/my-app:2.1.0-busybox1.35
├─── busybox:1.36
│  ├─── registry.stevedore.test/my-app:3.2.1-busybox1.36
```

- Start building the Docker images. It will create all `my-app` versions.
```sh
/docker $ stevedore build my-app
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 1/7 : ARG image_from_name
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 2/7 : ARG image_from_tag
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 3/7 : FROM ${image_from_name}:${image_from_tag}
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 1/7 : ARG image_from_name
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 2/7 : ARG image_from_tag
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 3/7 : FROM ${image_from_name}:${image_from_tag}
registry.stevedore.test/my-app:2.1.0-busybox1.35 Step 1/7 : ARG image_from_name
registry.stevedore.test/my-app:2.1.0-busybox1.35 Step 2/7 : ARG image_from_tag
registry.stevedore.test/my-app:2.1.0-busybox1.35 Step 3/7 : FROM ${image_from_name}:${image_from_tag}
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  1.36:  Pulling from library/busybox
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  205dae5015e7:  Pull complete
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  Digest: sha256:7b3ccabffc97de872a30dfd234fd972a66d247c8cfc69b0550f276481852627c
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  Status: Downloaded newer image for busybox:1.36
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> 66ba00ad3de8
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 4/7 : ARG message=my-app!
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> Running in 7abead2499e1
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> 5f0e093a1a94
registry.stevedore.test/my-app:2.1.0-busybox1.35 ‣  1.35:  Pulling from library/busybox
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  1.35:  Pulling from library/busybox
registry.stevedore.test/my-app:2.1.0-busybox1.35 ‣  2461e8255644:  Pull complete
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  2461e8255644:  Pull complete
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  Digest: sha256:f4ed5f2163110c26d42741fdc92bd1710e118aed4edb19212548e8ca4e5fca22
registry.stevedore.test/my-app:3.2.1-busybox1.35 ‣  Status: Downloaded newer image for busybox:1.35
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> f68fa78323e7
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 4/7 : ARG message=my-app!
registry.stevedore.test/my-app:2.1.0-busybox1.35 ‣  Digest: sha256:f4ed5f2163110c26d42741fdc92bd1710e118aed4edb19212548e8ca4e5fca22
registry.stevedore.test/my-app:2.1.0-busybox1.35 ‣  Status: Image is up to date for busybox:1.35
registry.stevedore.test/my-app:2.1.0-busybox1.35 ---> f68fa78323e7
registry.stevedore.test/my-app:2.1.0-busybox1.35 Step 4/7 : ARG message=my-app!
registry.stevedore.test/my-app:2.1.0-busybox1.35 ---> Running in 72d2a3a19bd5
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> Running in d317a4adc9a3
registry.stevedore.test/my-app:2.1.0-busybox1.35 ---> 0116092ba9da
registry.stevedore.test/my-app:2.1.0-busybox1.35 Step 5/7 : RUN echo "Hey there! Welcome to ${message}" > /message.txt
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> 03bf5a7e8561
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 5/7 : RUN echo "Hey there! Welcome to ${message}" > /message.txt
registry.stevedore.test/my-app:2.1.0-busybox1.35 ---> Running in a4c2d4b4ccb4
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> Running in 41aacd4c92a7
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> 5de239d250df
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 6/7 : CMD ["cat","/message.txt"]
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> Running in 0c32a6be979e
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> e50f09eee688
registry.stevedore.test/my-app:3.2.1-busybox1.36 Step 7/7 : LABEL created_at=2023-02-03T22:42:55.122127826Z
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> Running in 4068fc05135e
registry.stevedore.test/my-app:3.2.1-busybox1.36 ---> d56ea6314840
registry.stevedore.test/my-app:3.2.1-busybox1.36  ‣ sha256:d56ea631484054906c3daa8289b064b91725714ddcda4e3874d1e0a7b6561c49
registry.stevedore.test/my-app:3.2.1-busybox1.36 Successfully built d56ea6314840
registry.stevedore.test/my-app:3.2.1-busybox1.36 Successfully tagged registry.stevedore.test/my-app:3.2.1-busybox1.36
registry.stevedore.test/my-app:2.1.0-busybox1.35 ---> 75869c3433a5
registry.stevedore.test/my-app:2.1.0-busybox1.35 Step 6/7 : CMD ["cat","/message.txt"]
registry.stevedore.test/my-app:2.1.0-busybox1.35 ---> Running in 761056fa55b5
registry.stevedore.test/my-app:2.1.0-busybox1.35 ---> 3e4e261bed5a
registry.stevedore.test/my-app:2.1.0-busybox1.35 Step 7/7 : LABEL created_at=2023-02-03T22:42:55.121439497Z
registry.stevedore.test/my-app:2.1.0-busybox1.35 ---> Running in f03539b96af7
registry.stevedore.test/my-app:2.1.0-busybox1.35 ---> 2e362f2c8f26
registry.stevedore.test/my-app:2.1.0-busybox1.35  ‣ sha256:2e362f2c8f265633aec978d68ab5dc1e57efbf27cc2fe442dd8b88c4d82e6582
registry.stevedore.test/my-app:2.1.0-busybox1.35 Successfully built 2e362f2c8f26
registry.stevedore.test/my-app:2.1.0-busybox1.35 Successfully tagged registry.stevedore.test/my-app:2.1.0-busybox1.35
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> 0fa400361b6b
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 6/7 : CMD ["cat","/message.txt"]
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> Running in 1ff331631429
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> 778d69df86fa
registry.stevedore.test/my-app:3.2.1-busybox1.35 Step 7/7 : LABEL created_at=2023-02-03T22:42:55.121439497Z
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> Running in 07ed3cc21395
registry.stevedore.test/my-app:3.2.1-busybox1.35 ---> 612c3f025134
registry.stevedore.test/my-app:3.2.1-busybox1.35  ‣ sha256:612c3f02513437418e3bd10f784af44efe09b36583b3c1b9ba54138d3158f1f4
registry.stevedore.test/my-app:3.2.1-busybox1.35 Successfully built 612c3f025134
registry.stevedore.test/my-app:3.2.1-busybox1.35 Successfully tagged registry.stevedore.test/my-app:3.2.1-busybox1.35
```

Validate the recently created images.
```sh
/docker $ docker images
REPOSITORY                       TAG                  IMAGE ID       CREATED         SIZE
registry.stevedore.test/my-app   3.2.1-busybox1.35    612c3f025134   2 minutes ago   4.86MB
registry.stevedore.test/my-app   3.2.1-busybox1.36    d56ea6314840   2 minutes ago   4.87MB
registry.stevedore.test/my-app   2.1.0-busybox1.35    2e362f2c8f26   2 minutes ago   4.86MB
busybox                          1.36                 66ba00ad3de8   4 weeks ago     4.87MB
busybox                          1.35                 f68fa78323e7   6 weeks ago     4.86MB
```

### Promote the images to Docker registry
Now that we already have the images created, and since we did not set the push after build, we promote the image to Docker registry and we push it to the `stable` namespace.
```sh
/docker $ stevedore promote registry.stevedore.test/my-app:3.2.1-busybox1.36 --promote-image-registry-namespace stable
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  The push refers to repository [registry.stevedore.test/stable/my-app]

registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  b64792c17e4a:  Pushed
registry.stevedore.test/my-app:3.2.1-busybox1.36 ‣  3.2.1-busybox1.36: digest: sha256:b1aa5de2f4bf9c031a2047a87fb5c556d0d436123316cad81078462648e58d4b size: 528
```
