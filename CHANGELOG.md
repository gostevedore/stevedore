# CHANGELOG

## [Unreleased]

### Bumped

- Bump up github.com/Masterminds/semver/v3 from v3.2.0 to v3.2.1
- Bump up github.com/apenella/go-docker-builder from v0.7.8 to v0.7.9. It fixes an error when creating the build context having symbolic links.
- Bump up github.com/aws/aws-sdk-go-v2 from v1.17.2 to v1.23.1
- Bump up github.com/aws/aws-sdk-go-v2/config from v1.18.4 to v1.25.5
- Bump up github.com/aws/aws-sdk-go-v2/credentials from v1.13.4 to v1.16.4
- Bump up github.com/aws/aws-sdk-go-v2/service/ecr from v1.17.24 to v1.23.1
- Bump up github.com/aws/aws-sdk-go-v2/service/sts from v1.17.6 to v1.25.4
- Bump up github.com/docker/distribution from v2.8.2+incompatible to v2.8.3+incompatible
- Bump up github.com/go-git/go-git/v5 from v5.6.1 to v5.10.0
- Bump up github.com/gruntwork-io/terratest from v0.41.9 to v0.46.7
- Bump up github.com/spf13/afero from v1.9.5 to v1.10.0
- Bump up github.com/spf13/cobra from v1.6.1 to v1.8.0
- Bump up github.com/spf13/viper from v1.14.0 to v1.17.0
- Bump up github.com/stretchr/testify from v1.8.2 to v1.8.4
- Bump up go.uber.org/zap from v1.24.0 to v1.26.0
- Bump up golang.org/x/term from v0.5.0 to v0.14.0

### Changed

- By default, install the Stevedore binary in $HOME/bin
- Do not use Golang Docker images as base images in the testing applications
- Implement retry mechanism in the functional tests
- The installation script uses `/bin/sh` instead `/bin/bash`
- The release process creates ARM binaries
- Use Docker 24.0 in the examples and tests

## [v0.11.1]

### Added

- Included examples
- New variable mapping named `image_from_fully_qualified_name` that provides the fully qualified name of the parent Docker image as a build argument
- New variable mapping named `image_fully_qualified_name` that provides the fully qualified name of current image as a build argument. It is used by the `ansible-plybook` driver

### Bumped

- Bump up github.com/apenella/go-docker-builder to v0.7.8
- Bump up github.com/spf13/afero to v1.9.5
- Bump up github.com/docker/distribution to v2.8.2+incompatible
- Bump up github.com/docker/docker to v20.10.24+incompatible
- Bump up github.com/go-git/go-git/v5 to v5.6.1

### Fixed

- Install script uses the artefact name updated on v0.11.x
- Use the default variables mapping definition when in the builder is defined as empty builder
- On promote command, mark as deprecated the flags --promote-image-namespace and --promote-image-registry
- On promote command, enable semver tag is aware of the source image tag
- On build command, pull-parent-image pulls the parent image (fixed on github.com/apenella/go-docker-builder to v0.7.8)
- On promote command, remove image after push (fixed on github.com/apenella/go-docker-builder to v0.7.8)

## [v0.11.0]

### Added

- Installation script
- Create credentials local store
- Create credentials environment variables store
- Allow to cypher credentials store
- Create a semver generator to parse image tags
- Create a new promote service
- Create an infrastructure repository for docker promotions
- Create a handler to promote cli subcommand
- Create a logger object
- `BuildDriverer` interface accepts `Build(context.Context,BuildDriverOptions) error`
- Docker build driver accepts the username-password and private key as authentication methods for the git docker build context
- Struct to render images accepts dates through `DateRFC3339` and `DateRFC3339Nano` attributes
- Subcomand get images accepts filters
- Subcomand get builders accept filters

### Changed

- Updated license to Apache 2.0
- The whole project has been refactored following the ports-and-adapter architectonical design. Inside the `ìnternal` folder there are 5 main subfolders:
  - `application`: Use cases implementation
  - `core`: The domain objects and the repositories to interact with them
  - `entrypoint`: The elements that initialize and execute each command subsystem
  - `handler`: Defines the handler for each subcommand
  - `infrastrucutre`: Implementation for the driven and driver actor
- Bump up go to 1.19.6
- Bump up yaml to yaml.v3
- Bump up github.com/apenella/go-ansible to v1.1.7
- Bump github.com/apenella/go-docker-builder to v0.7.5
- Bump github.com/docker/docker to v20.10.23+incompatible
- Bump github.com/go-git/go-git/v5 to v5.5.2
- Bump github.com/spf13/afero to v1.9.4
- Bump github.com/stretchr/testify to v1.8.2
- Bump golang.org/x/term to v0.5.0
- [DEPRECATED FLAG] On the promote subcommand, use `remove-local-images-after-push` instead of `remove-promote-tags`
- [DEPRECATED FLAG] On the build subcommand, use `ansible-connection-local` instead of `connection-local`
- [DEPRECATED FLAG] On the build subcommand, use `ansible-intermediate-container-name` instead of `builder-name`
- [DEPRECATED FLAG] On the build subcommand, use `ansible-inventory-path` instead of `inventory`
- [DEPRECATED FLAG] On the build subcommand, use `ansible-limit` instead of `limit`
- [DEPRECATED FLAG] On the build subcommand, use `image-from-name` instead of `image-from`
- [DEPRECATED FLAG] On the build subcommand, use `image-registry-host` instead of `registry`
- [DEPRECATED FLAG] On the build subcommand, use `image-registry-namespace` instead of `namespace`
- [DEPRECATED FLAG] On the build subcommand, use `persistent-variable` instead of `set-persistent`
- [DEPRECATED FLAG] On the build subcommand, use `variable` instead of `set`
- [DEPRECATED FLAG] On the build subcommand, use `build-on-cascade` instead of `cascade`
- [DEPRECATED FLAG] On the build subcommand, use `concurrency` instead of `num-workers`
- [DEPRECATED FLAG] On the build subcommand, `no-push` is the stevedore default behaviour, use `push-after-build` flag to push an image
- [DEPRECATED FLAG] On the create credentials subcommand, the credentials id must be passed as a command argument instead of using `registry-host` flag
- [DEPRECATED FLAG] On the create credentials subcommand, `credentials-dir` is deprecated and will be ignored. Credentials parameters are set through the `credentials` section of the configuration file or using the flag `local-storage-path`

### Removed

- Image tags are not sanitized any more

## [0.10.3]

### Changed

- The command descriptions has been updated

### Fixed

- fix promote image with multiple tags and remove promoted tags enable. The image is deleted after first push.
- fix promote command was hidden

## [0.10.2]

### Fixed

- fix load semver templates from file on promote command

## [0.10.1]

### Fixed

- fix stevedore init creates configuration file with execution permissions
- fix promote does not uses -S flag to generate semantic version tags
- fix tags defined on the image are ignored

## [0.10.0]

### Added

- Make stevedore as public project
- New command to initizalize stevedore configuration
- Image builder could be defined as an string that identifies a global builder defined on `builders` configuration structure or could contain an in-line builder definition.
- Semver tags feature that autogenerate tags based on semantic version tags templates when an image tag is semver 2.0.0 compliance
- When is performed a build on cascade mode, can be defined the image's tree depth level of images to built
- New package logger that wraps zap logger and defines a global logger on a singleton mode
- New package console to use a global console variable on a singleton mode

### Removed

- Remove github.com/gostevedore/stevedore/internal/context package

### Changed

- On image definition struct, rename Childs to Children
- Default configuration location. Stevedore looks for default configuration file in: ./stevedore.yaml, ~/.config/stevedore/stevedore.yaml or ~/stevedore.yaml and loads the first configuration found.
- Default credentials locations is ~/.config/stevedore/credentials
- On configuration skip_push_images become push_images
- Update all packages to use the logger package to log messages
- Update all packages to use the console package to write missages to console
- Use golang context package instead of github.com/gostevedore/stevedore/internal/context
- Create an specific package for drivers
- Accept nil builders definition
- Include OS/arch and build date on version detail
- Update to github.com/apenella/go-docker-builder to v0.3.3
- Update to github.com/apenella/go-ansible to v0.6.1
- Update to github.com/docker/docker v20.10.0+incompatible
- Update to github.com/apenella/go-data-structures v0.2.0 which support cycles detection

### Fixed

- Fix raise an error when building a wildcarded image with an undefined version
- Fix set image from details to build images
- Fix negotiate docker API version on docker builder
- Fix sanitize tags to do not use '/' or ':'
- Fix to raise an error on generateTemplateGraph

---

**NOTE:** All version before v0.10.0 are not longer available but its changelog history is kept.

## [0.9.1]

### Added

- On promote command multiple tags can be set at once
- On build command is it possible to override the image name defined on the image tree by --image-name flag
- Start using https://keepachangelog.com/en/1.0.0/ for formating release notes file
- Return a non-zero exit code when there is a failure on the application

### Changed

- Cobra 1.0.0
- Viper 1.7.1
- Go-ansible 0.5.1
- Go-docker-builder 0.2.3
- Update dependencies
- Use pacakge github.com/apenella/go-common-utils/error to manage errors

### Fixed

- Fix use of slash on docker tags when docker driver is used
- Fix promote tag always use the source image tag
- Fix log_path from default configuration file was ignored

## [0.9.0]

- Include command to promote images
- Include a value object for docker image url
- Fix create complete path when creates a new credentials
- Fix on ansible-playbook builder do not set persistent_vars and vars as string when add an extravar

## [0.8.1]

- Set a default driver name when it is not defined
- Use var mappings on ansible-playbook builder
- New options types on get builders options
- Include wildcard index to images tree to identify wilcarded images
- Update Docker's driver to accept git repositories build
- Use go-docker-builder v0.3.1 which supports docker builds from a git context
- Fix create complete path when creates a new credentials
- Fix builders without options panics with a nil pointer exception
- Fix remove wildcarded version images when all images are listed/build
- Fix tree's node name on generateTemplateGraph

## [0.8.0]

- BuildOptions defined on types package
- Separate each driver on its own package
- New driver DockerBuilder
- New package to manage docker registry credentials
- Fix copy persistent vars from parent to child nodes

## [0.7.1]

- Wildcard version on cascade
- Copy tags and children when an image is copied

## [0.7.0]

- Use copy image on renderizeGraphRec
- Skip wildcard version on FindByName
- GenerateWilcardVersionNode return a node instead of a node list
- Use Image copy method on GenerateWilcardVersionNode
- Defined a method to copy an Image
- Render wildcarded images
- Apply persistent vars during image building process
- Skip wildcard version when finding an image on image tree index
- Define persistent vars on image tree
- Accept wildcard versions on image tree definition
- Define ansible playbook output prefix on ansible builder
- Cobra middleware to manage interruptions
- Fix command line vars are not keep
- Fix stevedore gets hung when is interrupted

## [0.6.2]

- Cobra middleware to manage panics
- Set writer on ansible playbook and dummy builders
- Define default number of workers when num workers flag is defined over 0
- Set default config folder as last folder to lookup the config file

## [0.6.1]

- Fix reloadConfiguration from file to overwrite configuration
- Update viper and cobra packages

## [0.6.0]

- Include go vet on Makefile
- Set a configuration parameter on stevedore commands
- New command to get configuration
- Define configuration instances into stevedore main
- Update go-common-utils to v0.1.1
- Generate configuration package

## [0.5.1]

- Use images graph to explore the images relationship during a build on cascade
- Store nodes instead of images on the index

## [0.5.0]

- Include index to search alternative name-version for each image
- Include on search index image tree definition version and rendered version
- Any image could be referenced and build from multiple parents on the same definition
- Included image index to lookup images on the image tree
- Options vars has precedence over default ones
- Make BuildOptions' copy to avoid overwrite values
- Define dummy driver
- Render images when is generated the images tree
- Define image tree on tree package
- Reorganize layout structure

## [0.4.2]

- Fix typo on imageFromNamespace value (source_target_namespace)

## [0.4.1]

- Write from-namespace flag value to ImageFromNamespace variable

## [0.4.0]

- Include Cobra completion command

## [0.3.1]

- Wait for childs build when build on cascade

## [0.3.0]

- Define namespace and image to builder
- Define the container builders name to allow concurrent buildings

## [0.2.2]

- Tagging on ci

## [0.2.1]

- Get verbs visible on Cobra

## [0.2.0]

- Change command version writer to stdout
- Change flag names. From skip-push-images to dry-run, from tags to versions and from extra-vars to set

## [0.1.1]

- New image tree definition
- Use node name on renderImagesDetailRec name/version split
- Fix rendering image content

## [0.1.0]

**features:**

- Get builders command
- Get images command
  - Return as table format which is the default output
  - Return as tree

## [0.0]

### Added

- Build command:
  - Execute ansible-playbooks to build docker images
  - Build isoleted images
  - Build images on cascade
  - Build for specific version
  - Use go templating on image tree definition
