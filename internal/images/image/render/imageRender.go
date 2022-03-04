package render

import (
	"bytes"
	"fmt"
	"html/template"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/images/image"
)

// ImageRender contains the information to render an image from template
type ImageRender struct {
	// Name    string
	// Version string
	// Parent  *image.Image
	// Image   ImageSerializer
}

// NewImageRender returns a new instance of the ImageRender
func NewImageRender() *ImageRender {
	return &ImageRender{}
}

// Render renders template image into incomming object
func (r *ImageRender) Render(name, version string, i *image.Image) (*image.Image, error) {
	var renderBuffer bytes.Buffer
	var renderedImage *image.Image
	var err error
	errContext := "(render::Render)"

	renderedImage, err = i.Copy()
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	renderObj := struct {
		Name    string
		Version string
		Parent  *image.Image
		Image   *image.Image
	}{
		Name:    name,
		Version: version,
		Parent:  i.Parent,
		Image:   renderedImage,
	}

	fmt.Printf(">>> %+v\n", renderObj)

	serialized, err := renderObj.Image.YAMLMarshal()
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	tmpl, err := template.New(renderObj.Name + ":" + renderObj.Version).Parse(string(serialized))
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	err = tmpl.Execute(&renderBuffer, renderObj)
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	err = renderObj.Image.YAMLUnmarshal(renderBuffer.Bytes())

	return renderedImage, nil
}
