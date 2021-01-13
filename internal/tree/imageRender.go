package tree

import (
	"bytes"
	"fmt"
	"html/template"
	"stevedore/internal/image"

	errors "github.com/apenella/go-common-utils/error"
	"gopkg.in/yaml.v2"
)

// ImageRender
type ImageRender struct {
	Name    string
	Version string
	Parent  *image.Image
	Image   *image.Image
}

// RenderizeImage
func RenderizeImage(r *ImageRender) error {
	var renderBuffer bytes.Buffer

	mItem, err := yaml.Marshal(&r.Image)
	if err != nil {
		return errors.New("(tree::RenderizeImage)", fmt.Sprintf("Error marshalling image '%s'", r.Name), err)
	}

	tmpl, err := template.New(r.Name + ":" + r.Version).Parse(string(mItem))
	if err != nil {
		return errors.New("(tree::RenderizeImage)", fmt.Sprintf("Error parsing template to renderize '%s'", r.Name), err)
	}

	err = tmpl.Execute(&renderBuffer, r)
	if err != nil {
		return errors.New("(tree::RenderizeImage)", fmt.Sprintf("Error renderizing image '%s'", r.Name), err)

	}
	err = yaml.Unmarshal(renderBuffer.Bytes(), &r.Image)
	if err != nil {
		return errors.New("(tree::RenderizeImage)", fmt.Sprintf("Error unmarshalling image '%s'", r.Name), err)
	}

	return nil
}
