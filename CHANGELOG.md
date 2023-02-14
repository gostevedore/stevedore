# CHANGELOG

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
