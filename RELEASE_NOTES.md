# RELEASE NOTES

## [Unreleased]

### Added
- Included examples

### Bumped
- Bump up github.com/apenella/go-docker-builder to v0.7.7
- Bump up github.com/spf13/afero to v1.9.5
- Bump up github.com/docker/distribution to v2.8.2+incompatible
- Bump up github.com/docker/docker to v20.10.24+incompatible
- Bump up github.com/go-git/go-git/v5 to v5.6.1

### Fixed
- Install script uses the artefact name updated on v0.11.0
- Use the default variables mapping definition when in the builder is defined as empty
- On promote command, mark as deprecated the flags --promote-image-namespace and --promote-image-registry
- On promote command, enable semver tag is aware of the source image tag
