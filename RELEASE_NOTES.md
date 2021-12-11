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
- On docker build driver, are accepted user/password and private key as auth method for git docker build context.

### Changed
- Use go-ansible master version
- Use go-docker-builder master version
- Drivers has been adapted to use go-ansible and go-docker-builder version
- Promote uses copy package for go-docker-builder
- Promote Cobra subcommand initializes all required repositories, services and handlers on its prerun function
- [DEPRECATED] On promote subcommand, use `remove-local-images-after-push` instead of `remove-promote-tags`
- Rename `Driverer` interface to `BuildDriverer`
- refactor ansible-playbook build driver to allow testing
- refactor docker build driver to allow testing
