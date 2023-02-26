# RELEASE NOTES

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
- Bump github.com/apenella/go-docker-builder to v0.7.4
- Bump github.com/docker/docker to v23.0.1+incompatible
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
