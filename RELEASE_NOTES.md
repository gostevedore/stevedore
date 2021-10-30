# RELEASE NOTES

## [Undefined]

### Added
- Create credentials store
- Create a semver generator
- Create a new promote service
- Create a infrastructure repository for docker promotions
- Create a handler for promote cli subcommand
- Create a logger object

### Changed
- Use go-ansible master version
- Use go-docker-builder v0.6.0 version
- Update `Driverer` interface to `Run(context.Context) error`
- Drivers has been adapted to use go-ansible and go-docker-builder version
- Promote uses copy package for go-docker-builder
- Promote Cobra subcommand initializes all required repositories, services and handlers on its prerun function
- [DEPRECATED] On promote subcommand, use `remove-local-images-after-push` instead of `remove-promote-tags`
