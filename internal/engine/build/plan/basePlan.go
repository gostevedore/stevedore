package plan

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/images/image"
)

// BasePlan is the base plan which all plans should extend
type BasePlan struct {
	images ImagesStorer
}

func (p *BasePlan) findImages(name string, versions []string) ([]*image.Image, error) {
	var images []*image.Image
	var imageAux *image.Image
	var err error

	errContext := "(plan::BasePlan::images)"

	if p.images == nil {
		return nil, errors.New(errContext, "Images storer is nil")
	}

	if versions == nil || len(versions) < 1 {
		images, err = p.images.FindByName(name)
		if err != nil {
			return nil, errors.New(errContext, err.Error())
		}
	} else {
		for _, version := range versions {
			imageAux, err = p.images.Find(name, version)
			if err != nil {
				return nil, errors.New(errContext, err.Error())
			}
			images = append(images, imageAux)
		}
	}

	return images, nil
}
