package images

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

const (
	NA              = "-"
	IN_LINE_BUILDER = "<in-line>"
)

// PlainOutputOptionsFunc is a function used to configure the service
type PlainOutputOptionsFunc func(*PlainOutput)

type PlainOutput struct {
	writer repository.ImagesPlainPrinter
}

func NewPlainOutput(options ...PlainOutputOptionsFunc) *PlainOutput {
	output := &PlainOutput{}
	output.Options(options...)
	return output
}

func WithWriter(w repository.ImagesPlainPrinter) PlainOutputOptionsFunc {
	return func(o *PlainOutput) {
		o.writer = w
	}
}

// Options configure the service
func (o *PlainOutput) Options(opts ...PlainOutputOptionsFunc) {
	for _, opt := range opts {
		opt(o)
	}
}

// outputHeader returns the header for the output
func outputHeader() []string {
	return []string{"NAME", "VERSION", "REGISTRY", "NAMESPACE", "BUILDER", "PARENT"}
}

func (o *PlainOutput) Output(list []*image.Image) error {
	errContext := "(output::images::PlainOutput::Output)"
	content := [][]string{}
	content = append(content, outputHeader())

	for _, image := range list {
		imageSlice, err := imageToOutputSlice(image)
		if err != nil {
			return errors.New(errContext, "", err)
		}
		content = append(content, imageSlice)
	}

	err := o.writer.PrintTable(content)
	if err != nil {
		return errors.New(errContext, "", err)
	}

	return nil
}

func imageToOutputSlice(i *image.Image) ([]string, error) {
	errContext := "(output::images::imageToOutputSlice)"

	res := []string{}

	if i.Name != "" {
		res = append(res, i.Name)
	} else {
		res = append(res, NA)
	}

	if i.Version != "" {
		res = append(res, i.Version)
	} else {
		res = append(res, NA)
	}

	if i.RegistryHost != "" {
		res = append(res, i.RegistryHost)
	} else {
		res = append(res, NA)
	}

	if i.RegistryNamespace != "" {
		res = append(res, i.RegistryNamespace)
	} else {
		res = append(res, NA)
	}

	if _, isString := i.Builder.(string); isString {
		res = append(res, i.Builder.(string))
	} else {
		if i.Builder == nil {
			res = append(res, NA)
		} else {
			res = append(res, IN_LINE_BUILDER)
		}
	}

	if i.Parent != nil {
		parent, err := i.Parent.DockerNormalizedNamed()
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}

		res = append(res, parent)
	}

	return res, nil
}
