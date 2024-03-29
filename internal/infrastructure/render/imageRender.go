package render

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/infrastructure/now"
)

// ImageRender contains the information to render an image from template
type ImageRender struct {
	// Name    string
	// Version string
	// Parent  *image.Image
	// Image   ImageSerializer
	now Nower
}

// NewImageRender returns a new instance of the ImageRender
func NewImageRender(n Nower) *ImageRender {
	return &ImageRender{
		now: n,
	}
}

// Render renders template image into incomming object
func (r *ImageRender) Render(name, version string, i *image.Image) (*image.Image, error) {
	var renderBuffer bytes.Buffer
	var renderedImage *image.Image
	var err error
	errContext := "(render::Render)"

	if i == nil {
		return nil, errors.New(errContext, "An image is required to render")
	}

	if r.now == nil {
		r.now = now.NewNow()
	}

	renderedImage, err = i.Copy()
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	renderObj := struct {
		Name            string
		Version         string
		Parent          *image.Image
		Image           *image.Image
		DateRFC3339     string
		DateRFC3339Nano string
	}{
		Name:            name,
		Version:         version,
		Parent:          i.Parent,
		Image:           renderedImage,
		DateRFC3339:     r.now.NowFunc()(time.RFC3339),
		DateRFC3339Nano: r.now.NowFunc()(time.RFC3339Nano),
	}

	serialized, err := renderObj.Image.YAMLMarshal()
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	tmpl, err := template.New(renderObj.Name + ":" + renderObj.Version).Parse(string(serialized))
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	err = tmpl.Execute(&renderBuffer, renderObj)
	if err != nil {
		return nil, errors.New(errContext, fmt.Sprintf("Error rendering image %s:%s from the following image definition template:\n\n%s\nInput values:\n%+v", name, version, string(serialized), renderObj), err)
	}

	err = renderObj.Image.YAMLUnmarshal(renderBuffer.Bytes())
	if err != nil {
		return nil, errors.New(errContext, "", err)
	}

	return renderedImage, nil
}
