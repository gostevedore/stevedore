# RELEASE NOTES

## [0.10.0]

### Added
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
