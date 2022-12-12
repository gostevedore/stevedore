# RELEASE NOTES

## [Undefined]

### Added
- Create credentials store
- Create a semver generator
- Create a new promote service
- Create a infrastructure repository for docker promotions
- Create a handler for promote cli subcommand
- Create a logger object
- `BuildDriverer` interface accepts `Build(context.Context,BuildDriverOptions) error`
- On docker build driver, are accepted user/password and private key as auth method for git docker build context
- Struct to render images accepts dates through `DateRFC3339` and `DateRFC3339Nano` attributes
- Included filters on get images
- Included filters on get builders

### Changed
- Bump to go1.19 
- Bump to yaml.v3
- Use go-ansible v1.1.7 version
- Use go-docker-builder v0.7.0 version
- Drivers has been adapted to use go-ansible and go-docker-builder version
- Promote uses copy package for go-docker-builder
- Promote Cobra subcommand initializes all required repositories, services and handlers on its prerun function
- [DEPRECATED] On promote subcommand, use `remove-local-images-after-push` instead of `remove-promote-tags`
- Rename `Driverer` interface to `BuildDriverer`
- Refactor ansible-playbook build driver to allow testing
- Refactor docker build driver to allow testing
- **BREAKING-CHANGES** package `"github.com/gostevedore/stevedore/internal/build"` has been replaced by `"github.com/gostevedore/stevedore/internal/builders"`.
- **BREAKING-CHANGES** package `"github.com/gostevedore/stevedore/internal/builders/builder"` has an strict definition of builder options
- `BuilderOptions` accepts multiple context that are unified before perform an image build
- Drivers receives `BuildDriverOptions` instead of `BuildOptions`
- Image uses docekr reference normalized names
- Images has been splitted to image configration and image as dto
- build engine has been replaced by application service build
- Image could be defined without specifing a builder, in that case `default` builder is used instead

### Removed
- Image tags are not sanetized any more
