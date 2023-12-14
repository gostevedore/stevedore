package semver

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"text/template"

	"github.com/Masterminds/semver/v3"
	errors "github.com/apenella/go-common-utils/error"
)

// SemVer is a sematinc version representation
type SemVer struct {
	Major      string
	Minor      string
	Patch      string
	PreRelease string
	Build      string
}

// NewSemVer return an struct containing the version input parsed. It returns an error when version does not match to semver
func NewSemVer(version string) (*SemVer, error) {

	sv := &SemVer{}

	v, err := semver.NewVersion(version)
	if err != nil {
		return nil, errors.New("(semver::NewSemVer)", fmt.Sprintf("Error creating new version '%s'", version), err)
	}

	sv.Major = strconv.FormatInt(int64(v.Major()), 10)
	sv.Minor = strconv.FormatInt(int64(v.Minor()), 10)
	sv.Patch = strconv.FormatInt(int64(v.Patch()), 10)

	if len(v.Prerelease()) > 0 {
		sv.PreRelease = v.Prerelease()
	}

	if len(v.Metadata()) > 0 {
		sv.Build = v.Metadata()
	}

	return sv, nil
}

// Validate
func Validate(version string) bool {
	_, err := semver.NewVersion(version)
	return err == nil
}

// VersionTree return semver versions from v and based on list templates. example: from 1.2.3 --> [1, 1.2, 1.2.3]. It returns an error when a template could not be parsed
func (v *SemVer) VersionTree(listTmpl []string) ([]string, error) {

	versionList := []string{}

	templateListVersions := template.New("listVersion")
	for _, templateItem := range listTmpl {
		buff := &bytes.Buffer{}
		template, err := templateListVersions.Parse(templateItem)
		if err != nil {
			return nil, errors.New("(semver::ListVersion)", fmt.Sprintf("Error parsing '%s'", templateItem), err)
		}

		// error is not being managed because is preferred to skip it and continue with the remaining templates
		_ = template.Execute(io.Writer(buff), v)
		versionList = append(versionList, buff.String())
	}

	return versionList, nil
}

// String return *SemVer in string format
func (v *SemVer) String() string {
	str := fmt.Sprintf("%s.%s.%s", v.Major, v.Minor, v.Patch)

	if v.PreRelease != "" {
		str = fmt.Sprintf("%s-%s", str, v.PreRelease)
	}

	if v.Build != "" {
		str = fmt.Sprintf("%s+%s", str, v.Build)
	}

	return str
}
