package name

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
)

// DefaultReferenceName generates images names
type DefaultReferenceName struct{}

// NewDefaultReferenceName return a DefaultReferenceName instance
func NewDefaultReferenceName() *DefaultReferenceName {
	return &DefaultReferenceName{}
}

// GenerateName return a name for a given image. In case of error an error is returned
func (n *DefaultReferenceName) GenerateName(i *image.Image) (string, error) {

	errContext := "(name::images::DefaultReferenceName::GenerateName)"

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

	return name, nil
}
