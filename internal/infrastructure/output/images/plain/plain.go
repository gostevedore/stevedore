package images

import (
	"fmt"

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
	writer        repository.ImagesPlainPrinter
	referenceName repository.ImageReferenceNamer
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

func WithReferenceName(ref repository.ImageReferenceNamer) PlainOutputOptionsFunc {
	return func(o *PlainOutput) {
		o.referenceName = ref
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
	return []string{"NAME", "VERSION", "BUILDER", "IMAGE FULL NAME", "PARENT"}
}

// Output writes into writer the images from list in plain format
func (o *PlainOutput) Output(list []*image.Image) error {
	errContext := "(output::images::PlainOutput::Output)"
	content := [][]string{}
	content = append(content, outputHeader())

	if o.writer == nil {
		return errors.New(errContext, "Images plain text output requires a writer")
	}

	for _, image := range list {
		imageSlice, err := o.imageToOutputSlice(image)
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

func (o *PlainOutput) imageToOutputSlice(i *image.Image) ([]string, error) {
	errContext := "(output::images::imageToOutputSlice)"

	if o.referenceName == nil {
		return nil, errors.New(errContext, "Images plain text output requires a reference name")
	}

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

	if _, isString := i.Builder.(string); isString {
		res = append(res, i.Builder.(string))
	} else {
		if i.Builder == nil {
			res = append(res, NA)
		} else {
			res = append(res, IN_LINE_BUILDER)
		}
	}

	if i.Name != "" && i.Version != "" {
		ref, err := o.referenceName.GenerateName(i)
		if err != nil {
			// instead of returned the error fmt is used as a fallback
			ref = fmt.Sprintf("%s:%s", i.Name, i.Version)
		}
		res = append(res, ref)
	} else {
		res = append(res, NA)
	}

	if i.Parent != nil {
		parent, err := o.referenceName.GenerateName(i.Parent)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}
		res = append(res, parent)
	} else {
		res = append(res, NA)
	}

	return res, nil
}
