package render

import (
	"bytes"
	"html/template"
	"time"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/images/image"
	"github.com/gostevedore/stevedore/internal/images/image/render/now"
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
		return nil, errors.New(errContext, err.Error())
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
	if err != nil {
		return nil, errors.New(errContext, err.Error())
	}

	return renderedImage, nil
}
