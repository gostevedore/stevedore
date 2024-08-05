package name

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/distribution/reference"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
)

// DockerNormalizedReferenceName generates docker normalized names
type DockerNormalizedReferenceName struct{}

// NewDockerNormalizedReferenceName return a DockerNormalizedReferenceName instance
func NewDockerNormalizedReferenceName() *DockerNormalizedReferenceName {
	return &DockerNormalizedReferenceName{}
}

// GenerateName return a name for a given image. In case of error an error is returned
func (n *DockerNormalizedReferenceName) GenerateName(i *image.Image) (string, error) {
	var err error
	var named reference.Named
	errContext := "(name::images::DockerNormalizedReferenceName::GenerateName)"

	if i.Name == "" {
		return "", errors.New(errContext, "Image reference name can not be generated because image name is undefined")
	}

	if i.Version == "" {
		return "", errors.New(errContext, "Image reference name can not be generated because image version is undefined")
	}

	name := fmt.Sprintf("%s:%s", i.Name, i.Version)

	if i.RegistryNamespace != "" {
		name = fmt.Sprintf("%s/%s", i.RegistryNamespace, name)
	}

	if i.RegistryHost != "" {
		name = fmt.Sprintf("%s/%s", i.RegistryHost, name)
	}

	named, err = reference.ParseNormalizedNamed(name)
	if err != nil {
		return "", errors.New(errContext, fmt.Sprintf("Image name '%s' could not be normalized", name), err)
	}

	return named.String(), nil

}
