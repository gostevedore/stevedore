# RELEASE NOTES

## [v0.11.2]

### Fixed

- Solved the error caused when creating a build context containing symbolic links, by bumping up github.com/apenella/go-docker-builder to v0.8.1

### Added

- Include [go vet](https://pkg.go.dev/cmd/vet) analysis in the CI workflow
- Include the [errcheck](https://github.com/kisielk/errcheck) tool in the CI workflow

### Bumped

- Bump github.com/Masterminds/semver/v3 from v3.2.0 to v3.2.1
- Bump github.com/apenella/go-docker-builder from v0.7.8 to v0.8.1. It fixes an error when creating the build context having symbolic links.
- Bump github.com/aws/aws-sdk-go-v2 from v1.17.2 to v1.23.1
- Bump github.com/aws/aws-sdk-go-v2/config from v1.18.4 to v1.25.5
- Bump github.com/aws/aws-sdk-go-v2/credentials from v1.13.4 to v1.16.4
- Bump github.com/aws/aws-sdk-go-v2/service/ecr from v1.17.24 to v1.23.1
- Bump github.com/aws/aws-sdk-go-v2/service/sts from v1.17.6 to v1.25.4
- Bump github.com/docker/distribution from v2.8.2+incompatible to v2.8.3+incompatible
- Bump github.com/go-git/go-git/v5 from v5.6.1 to v5.10.0
- Bump github.com/gruntwork-io/terratest from v0.41.9 to v0.46.7
- Bump github.com/spf13/afero from v1.9.5 to v1.10.0
- Bump github.com/spf13/cobra from v1.6.1 to v1.8.0
- Bump github.com/spf13/viper from v1.14.0 to v1.17.0
- Bump github.com/stretchr/testify from v1.8.2 to v1.8.4
- Bump go.uber.org/zap from v1.24.0 to v1.26.0
- Bump golang.org/x/term from v0.5.0 to v0.14.0

### Changed

- By default, install the Stevedore binary in $HOME/bin
- Do not build Golang applications on functional tests
- Implement a retry mechanism in the functional tests
- Isolate test execution from the guest host by running functional and unit tests inside a Docker container
- The installation script uses `/bin/sh` instead `/bin/bash`
- The release process creates ARM binaries
- Use Docker 24.0 in the examples and tests
