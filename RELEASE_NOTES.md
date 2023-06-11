# RELEASE NOTES

## [Unreleased]

### Added
- Included examples
- New variable mapping named `image_from_fully_qualified_name` that provides the fully qualified name of the parent Docker image as a build argument.

### Bumped
- Bump up github.com/apenella/go-docker-builder to v0.7.8
- Bump up github.com/spf13/afero to v1.9.5
- Bump up github.com/docker/distribution to v2.8.2+incompatible
- Bump up github.com/docker/docker to v20.10.24+incompatible
- Bump up github.com/go-git/go-git/v5 to v5.6.1

### Fixed
- Install script uses the artefact name updated on v0.11.x
- Use the default variables mapping definition when in the builder is defined as empty builder
- On promote command, mark as deprecated the flags --promote-image-namespace and --promote-image-registry
- On promote command, enable semver tag is aware of the source image tag
- On build command, pull-parent-image pulls the parent image (fixed on github.com/apenella/go-docker-builder to v0.7.8)
- On promote command, remove image after push (fixed on github.com/apenella/go-docker-builder to v0.7.8)
