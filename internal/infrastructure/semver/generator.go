package semver

import (
	errors "github.com/apenella/go-common-utils/error"
)

type SemVerGenerator struct {
	semver *SemVer
}

func NewSemVerGenerator() *SemVerGenerator {
	return &SemVerGenerator{}
}

func (g *SemVerGenerator) GenerateSemvVer(version string) error {

	errContext := "(semver::GenerateSemvVer)"
	sv, err := NewSemVer(version)
	if err != nil {
		return errors.New(errContext, "Semantic version could not be generated", err)
	}
	g.semver = sv
	return nil
}

func (g *SemVerGenerator) GenerateVersionTree(tmpl []string) ([]string, error) {
	return g.semver.VersionTree(tmpl)
}

func (g *SemVerGenerator) GenerateSemverList(versions []string, tmpls []string) ([]string, error) {
	var err error
	list := []string{}
	versionsMap := make(map[string]struct{})

	errContext := "(semver::GenerateSemverList)"

	for _, version := range versions {

		err = g.GenerateSemvVer(version)
		if err == nil {
			svtree, err := g.GenerateVersionTree(tmpls)
			if err != nil {
				return nil, errors.New(errContext, "Version tree could not be generated in order to create semver list", err)
			} else {
				for _, v := range svtree {
					versionsMap[v] = struct{}{}
				}
			}
		}
	}

	for version := range versionsMap {
		list = append(list, version)
	}

	return list, nil
}
