package plan

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/gostevedore/stevedore/internal/core/domain/image"
	"github.com/gostevedore/stevedore/internal/core/ports/repository"
)

// BasePlan is the base plan which all plans should extend
type BasePlan struct {
	images repository.ImagesStorerReader
}

func (p *BasePlan) findImages(name string, versions []string) ([]*image.Image, error) {
	var images []*image.Image
	var imageAux []*image.Image
	var err error

	errContext := "(plan::BasePlan::images)"

	if p.images == nil {
		return nil, errors.New(errContext, "Images storer is nil")
	}

	if versions == nil || len(versions) < 1 {
		images, err = p.images.FindByName(name)
		if err != nil {
			return nil, errors.New(errContext, "", err)
		}
	} else {
		for _, version := range versions {
			imageAux, err = p.images.FindGuaranteed(name, version)
			if err != nil {
				return nil, errors.New(errContext, "", err)
			}
			images = append(images, imageAux...)
		}
	}

	if len(images) < 1 {
		msg := ""
		if len(versions) < 1 {
			msg = fmt.Sprintf("The image '%s' seems not to be defined", name)
		} else {
			msg = fmt.Sprintf("The version(s) %s for the image '%s' seems not to be defined", versions, name)
		}

		return nil, errors.New(errContext, msg)
	}

	return images, nil
}
