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
		sv.PreRelease = fmt.Sprintf("%s", v.Prerelease())
	}

	if len(v.Metadata()) > 0 {
		sv.Build = fmt.Sprintf("%s", v.Metadata())
	}

	return sv, nil
}

// Validate
func Validate(version string) bool {
	_, err := semver.NewVersion(version)
	if err != nil {
		return false
	}
	return true
}

// VersionTree return semver versions from v and based on list templates. example: from 1.2.3 --> [1, 1.2, 1.2.3]. It returns an error when a template could not be parsed
func (v *SemVer) VersionTree(listTmpl []string) ([]string, error) {

	res := []string{}

	tmpl := template.New("listVersion")
	for _, t := range listTmpl {
		buff := &bytes.Buffer{}
		tx, err := tmpl.Parse(t)
		if err != nil {
			return nil, errors.New("(semver::ListVersion)", fmt.Sprintf("Error parsing '%s'", t), err)
		}
		tx.Execute(io.Writer(buff), v)
		res = append(res, buff.String())
	}

	return res, nil
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
