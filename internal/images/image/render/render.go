package render

import (
	"bytes"
	"html/template"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/images/image"
)

// ImageRender contains the information to render an image from template
type ImageRender struct {
	Name    string
	Version string
	Parent  *image.Image
	Image   ImageSerializer
}

// NewImageRender returns a new instance of the ImageRender
func NewImageRender(name, version string, parent *image.Image, image ImageSerializer) *ImageRender {
	return &ImageRender{
		Name:    name,
		Version: version,
		Parent:  parent,
		Image:   image,
	}
}

// Render renders template image into incomming object
func (r *ImageRender) Render() error {
	var renderBuffer bytes.Buffer
	errContext := "(render::Render)"

	serialized, err := r.Image.YAMLMarshal()
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	tmpl, err := template.New(r.Name + ":" + r.Version).Parse(string(serialized))
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	err = tmpl.Execute(&renderBuffer, r)
	if err != nil {
		return errors.New(errContext, err.Error())
	}

	err = r.Image.YAMLUnmarshal(renderBuffer.Bytes())

	return nil
}
