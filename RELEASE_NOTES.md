# RELEASE NOTES

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
- Bump up github.com/docker/docker from v20.10.24+incompatible to v24.0.7+incompatible
- Bump up github.com/go-git/go-git/v5 from v5.6.1 to v5.10.0
- Bump up github.com/gruntwork-io/terratest from v0.41.9 to v0.46.7
- Bump up github.com/spf13/afero from v1.9.5 to v1.10.0
- Bump up github.com/spf13/cobra from v1.6.1 to v1.8.0
- Bump up github.com/spf13/viper from v1.14.0 to v1.17.0
- Bump up github.com/stretchr/testify from v1.8.2 to v1.8.4
- Bump up go.uber.org/zap from v1.24.0 to v1.26.0
- Bump up golang.org/x/term from v0.5.0 to v0.14.0

### Changed

- The installation script uses `/bin/sh` instead `/bin/bash`
- The release process creates ARM binaries
- By default, install the Stevedore binary in $HOME/bin
