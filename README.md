
![stevedore-logo](docs/logo/logo_4_stevedore.png "Stevedore logo")

---

**Stevedore manages and governs the Docker image's building process**


<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->
# Contents
- [Features](#features)
- [Concepts](#concepts)
- [Quick start](#quick-start)

<!-- /code_chunk_output -->

# Features

- Define docker images relationship and orchestarte its building process.
> **_use case:_**
 You have an image shipped with a custom setup or user  and that image is configured as image `FROM` on all your microservices. Then every time you rebuild that image to include new security patches, all images that are depends to it are automatically built.

- Standarize and reuse Docker's images definition
> **_use case:_**
You define all your microservices based on the same skeleton, with the same way of testing or building them.
Then, you could define one Dockerfile and reuse it to generate an image for any of those microservices.

- Generate automatically semver tags when main tag is semver 2.0.0 compliance
> **_use case:_**
Supose that you are tagging your images using semantic versioning. In that case you could generate automatically new tags based on semver tree. 
>
>_example:_
_You have the tag `2.1.0`, then tags `2` and `2.1` are also generated._

- Centralized Docker registry's credentials store
> **_use case:_**
Supose that your Docker registry needs to be secured and should only be accessed by continuous delivery procedures. Then, `stevedore` could be authenticated on your behalf to push any created image to a Docker registry.

- Promote images to another Docker registry or to another registry namespace
> **_use case:_**
Your production environment only accept to pull images from specific Docker registry, and that docker registry is only used by your production environment. Supose that you have built and pushed an image to your staging Docker registry and you have passed all your end to end test. Then, you can promote the image from your staging Docker registry to your production one.

# Concepts

- Driver
    - Dockerfile
    - Ansible
- Builder
- Image tree
- Image

# Quick start

1. Initialize the stevedore project
```
stevedore init
```
**Usage:**
```
Create stevedore configuration file

Usage:
  stevedore init [flags]

Aliases:
  init

Flags:
  -C, --build-on-cascade                             On build, start children images building once an image build is finished. Its default value is 'false'
  -b, --builder-path-file string                     Builders location path. Its default value is 'stevedore.yaml'
  -d, --credentials-dir string                       Location path to store docker registry credentials. Its default value is 'credentials'
  -p, --credentials-password credentials-regristry   Docker registry password. It is ignored unless credentials-regristry value is defined
  -r, --credentials-registry-host string             Docker registry host to register credentials
  -u, --credentials-username credentials-regristry   Docker registry username. It is ignored unless credentials-regristry value is defined
  -s, --enable-semver-tags                           Generate extra tags when the main image tags is semver 2.0.0 compliance. Its default value is 'false'
      --force                                        Force to create configuration file when the file already exists
  -h, --help                                         help for init
  -l, --log-path-file string                         Log file location path. Its default value is '/dev/null'
  -P, --no-push-images                               On build, push images automatically after it finishes. Its default value is 'true'
  -w, --num-workers int                              It defines the number of workers to build images which corresponds to the number of images that can be build concurrently. Its default value is '4' (default -1)
  -T, --semver-tags-template strings                 List of templates which define those extra tags to generate when 'semantic_version_tags_enabled' is enabled. Its default value is '{{ .Major }}.{{ .Minor }}.{{ .Patch }}'
  -t, --tree-path-file string                        Images tree location path. Its default value is 'stevedore.yaml'

Global Flags:
  -c, --config string   Configuration file location path
```

2. Load credentials, in case you need
```
stevedore create credentials
```
**Usage:**
```
Create stevedore docker registry credentials

Usage:
  stevedore create credentials [flags]

Aliases:
  credentials, auth

Flags:
  -d, --credentials-dir string   Location path to store docker registry credentials (default "/home/aleix/.config/stevedore/credentials")
  -h, --help                     help for credentials
  -p, --password string          Docker registry password
  -r, --registry-host string     Docker registry host to register credentials
  -u, --username string          Docker registry username

Global Flags:
  -c, --config string   Configuration file location path
```

3. Define builders
4. Define the image tree
5. Start building
```
stevedore build
```
**Usage:**
```
Build images

Usage:
  stevedore build <image> [flags]

Flags:
  -b, --builder-name string            Intermediate builder's container name [only applies to ansible-playbook builders]
  -C, --cascade                        Build images on cascade. Children's image build is started once the image build finishes
  -d, --cascade-depth int              Number images levels to build when build on cascade is executed (default -1)
  -L, --connection-local               Use local connection for ansible [only applies to ansible-playbook builders]
      --debug                          Enable debug mode to show build options
  -D, --dry-run                        Run a dry-run build
  -S, --enable-semver-tags             Generate a set of tags for the image based on the semantic version tree when main version is semver 2.0.0 compliance
  -h, --help                           help for build
  -I, --image-from string              Image (FROM) parent's name
  -N, --image-from-namespace string    Image (FROM) parent's registry namespace
  -R, --image-from-registry string     Image (FROM) parent's registry host
  -V, --image-from-version string      Image (FROM) parent's version
  -i, --image-name string              Image name- It overrides image tree image name
  -v, --image-version strings          Image versions to be built. One or more image versions could be built
  -H, --inventory string               Specify inventory hosts' path or comma separated list of hosts [only applies to Ansible builders]
  -l, --limit string                   Further limit selected hosts to an additional pattern [only applies to Ansible builders]
  -n, --namespace string               Image's registry namespace where image will be stored
  -P, --no-push                        Do not push the image to registry once it is built
  -w, --num-workers int                Number of workers to execute builds
  -r, --registry string                Image's registry host where image will be stored
  -T, --semver-tags-template strings   List templates to generate tags following semantic version expression
  -s, --set strings                    Set variables to use during the build. The format of each variable must be <key>=<value>
  -p, --set-persistent strings         Set persistent variables to use during the build. A persistent variable will be available on child image during its build and could not be overwrite. The format of each variable must be <key>=<value>
  -t, --tag strings                    Give an extra tag for the docker image

Global Flags:
  -c, --config string   Configuration file location path

```