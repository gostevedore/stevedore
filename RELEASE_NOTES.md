# RELEASE NOTES

## Undefined

### Added
- Create credentials local store
- Create credentials environment variables store
- Allow to cypher credentials store
- Create a semver generator
- Create a new promote service
- Create an infrastructure repository for docker promotions
- Create a handler to promote cli subcommand
- Create a logger object
- `BuildDriverer` interface accepts `Build(context.Context,BuildDriverOptions) error`
- Docker build driver accepts the username-password and private key as authentication methods for the git docker build context
- Struct to render images accepts dates through `DateRFC3339` and `DateRFC3339Nano` attributes
- Subcomand get images accepts filters
- Subcomand get builders accepts filters

### Changed
- The whole project has been refactored following the ports-and-adapter architectonical design. Inside the `ìnternal` folder there are 5 main subfolders:
  - `application`: Use cases implementation
  - `core`: The domain objects and the repositories to interact with them
  - `entrypoint`: The elements that initialize and execute each command subsystem
  - `handler`: Defines the handler for each subcommand
  - `infrastrucutre`: Implementation for the driven and driver actor
- Bump to go1.19 
- Bump to yaml.v3
- Bump to `go-ansible` v1.1.7 version
- Bump to `go-docker-builder` v0.7.2 version

- 
- [DEPRECATED] On promote subcommand, use `remove-local-images-after-push` instead of `remove-promote-tags`



### Removed
- Image tags are not sanetized any more
