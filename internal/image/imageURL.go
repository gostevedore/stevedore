package image

import (
	"strings"

	errors "github.com/apenella/go-common-utils/error"
)

const (
	tagSeparator = ":"
	urlSeparator = "/"
)

type ImageURL struct {
	Registry  string
	Namespace string
	Name      string
	Tag       string
}

func Parse(imageName string) (*ImageURL, error) {

	parsedImageName := &ImageURL{}

	imageNameBlock1 := strings.Split(imageName, tagSeparator)

	if len(imageNameBlock1) > 2 {
		return nil, errors.New("(ImageURL::Parse)", "Invalid image name")
	}

	if len(imageNameBlock1) == 2 {
		parsedImageName.Tag = sanitizeTag(imageNameBlock1[1])
	}

	imageNameBlock2 := strings.Split(imageNameBlock1[0], urlSeparator)
	// if len(imageNameBlock2) > 3 {
	// 	return nil, errors.New("(ImageURL::Parse)", "Invalid image name")
	// }

	parsedImageName.Name = imageNameBlock2[len(imageNameBlock2)-1]

	// remove name from imageNameBlock2
	if len(imageNameBlock2) > 1 {
		imageNameBlock2 = imageNameBlock2[:len(imageNameBlock2)-1]

		// get registry host
		if len(imageNameBlock2) > 1 {
			parsedImageName.Registry = imageNameBlock2[0]
			imageNameBlock2 = imageNameBlock2[1:]
		}

		// get image namespace
		parsedImageName.Namespace = strings.Join(imageNameBlock2, urlSeparator)

	}

	return parsedImageName, nil
}

func (i *ImageURL) URL() (string, error) {
	if i.Name == "" {
		return "", errors.New("(image::URL)", "Image name is not defined")
	}

	str := i.Name

	if i.Tag != "" {
		tag := sanitizeTag(i.Tag)
		str = strings.Join([]string{str, tag}, tagSeparator)
	}

	if i.Namespace != "" {
		str = strings.Join([]string{i.Namespace, str}, urlSeparator)
	}

	if i.Registry != "" {
		str = strings.Join([]string{i.Registry, str}, urlSeparator)
	}

	return str, nil
}

func sanitizeTag(input string) string {

	chars := map[string]string{
		"/": "_",
		":": "_",
	}

	for originalChar, newChar := range chars {
		input = strings.ReplaceAll(input, originalChar, newChar)
	}

	return input
}
