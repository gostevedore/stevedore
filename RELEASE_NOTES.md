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

### Changed
- Use go-ansible master version
- Use go-docker-builder master version
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

### Removed
- Image tags are not sanetized any more