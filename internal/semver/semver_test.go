package semver

import (
	goerrors "errors"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/stretchr/testify/assert"
)

// TestParser tests Parse
func TestNewSemVer(t *testing.T) {

	tests := []struct {
		desc     string
		version  string
		err      error
		expected *SemVer
	}{
		{
			desc:    "Testing version major.minor.patch-prerelease+build",
			version: "1.2.3-rc0+build123",
			expected: &SemVer{
				Major:      "1",
				Minor:      "2",
				Patch:      "3",
				PreRelease: "rc0",
				Build:      "build123",
			},
			err: &errors.Error{},
		},
		{
			desc:    "Testing version major.minor.patch-prerelease",
			version: "1.2.3-rc0",
			expected: &SemVer{
				Major:      "1",
				Minor:      "2",
				Patch:      "3",
				PreRelease: "rc0",
			},
			err: &errors.Error{},
		},
		{
			desc:    "Testing version major.minor.patch+build",
			version: "1.2.3+build123",
			expected: &SemVer{
				Major: "1",
				Minor: "2",
				Patch: "3",
				Build: "build123",
			},
			err: &errors.Error{},
		},
		{
			desc:    "Testing version major.minor.patch",
			version: "1.2.3",
			expected: &SemVer{
				Major: "1",
				Minor: "2",
				Patch: "3",
			},
			err: &errors.Error{},
		},
		{
			desc:    "Testing version version major.minor",
			version: "1.2",
			expected: &SemVer{
				Major: "1",
				Minor: "2",
				Patch: "0",
			},
			err: &errors.Error{},
		},
		{
			desc:    "Testing version version major.minor-prerelease-hyphen",
			version: "1.2-rc1-0.0",
			expected: &SemVer{
				Major:      "1",
				Minor:      "2",
				Patch:      "0",
				PreRelease: "rc1-0.0",
			},
			err: errors.New("(semver::NewSemVer)", "Error creating new version '1.2-rc1-0.0'",
				goerrors.New("Invalid character(s) found in minor number \"2-rc1-0\"")),
		},
		{
			desc:     "Testing version invalid matching version [major.minor-prerelease+build+build]",
			version:  "1.2-rc1+build+build",
			expected: nil,
			err: errors.New("(semver::NewSemVer)", "Error creating new version '1.2-rc1+build+build'",
				goerrors.New("Invalid Semantic Version")),
		},
		{
			desc:     "Testing version invalid matching version [major.minor-prerelease+build+build]",
			version:  "1.2-rc1.1",
			expected: &SemVer{},
			err:      &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			sv, err := NewSemVer(test.version)

			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, sv)
			}

		})
	}
}

// TestListVersion
func TestVersionTree(t *testing.T) {

	tests := []struct {
		desc     string
		version  *SemVer
		tmpl     []string
		expected []string
		err      error
	}{
		{
			desc:     "Testing list version major, major.minor, major.minor.patch from major.minor.patch-prerelease",
			tmpl:     []string{"{{ .Major }}", "{{ .Major}}.{{ .Minor }}", "{{ .Major}}.{{ .Minor }}.{{ .Patch }}"},
			expected: []string{"1", "1.2", "1.2.3"},
			version: &SemVer{
				Major:      "1",
				Minor:      "2",
				Patch:      "3",
				PreRelease: "rc0",
				Build:      "build123",
			},
		},
		{
			desc:     "Testing list version major, major.minor, major.minor.patch from major.minor.patch-prerelease",
			tmpl:     []string{"{{ .Major }}", "{{ .Major}}.{{ .Minor }}", "{{ .Major}}.{{ .Minor }}.{{ .Patch }}"},
			expected: []string{"1", "1.2", "1.2.3"},
			version: &SemVer{
				Major:      "1",
				Minor:      "2",
				Patch:      "3",
				PreRelease: "rc0",
			},
		},
		{
			desc:     "Testing list version major, major.minor, major.minor.patch from major.minor.patch+build",
			tmpl:     []string{"{{ .Major }}", "{{ .Major}}.{{ .Minor }}", "{{ .Major}}.{{ .Minor }}.{{ .Patch }}"},
			expected: []string{"1", "1.2", "1.2.3"},
			version: &SemVer{
				Major: "1",
				Minor: "2",
				Patch: "3",
				Build: "build123",
			},
		},
		{
			desc:     "Testing list version major+build from major.minor.patch+build",
			tmpl:     []string{"{{ .Major }}{{ with .Build }}+{{ . }}{{ end }}"},
			expected: []string{"1+build"},
			version: &SemVer{
				Major: "1",
				Minor: "2",
				Patch: "3",
				Build: "build",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			l, err := test.version.VersionTree(test.tmpl)
			if err != nil && assert.Error(t, err) {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, l)
			}
		})
	}
}

// TestString tests String
func TestString(t *testing.T) {

	tests := []struct {
		desc     string
		version  *SemVer
		expected string
	}{
		{
			desc:     "Testing version major.minor.patch-prerelease+build",
			expected: "1.2.3-rc0+build123",
			version: &SemVer{
				Major:      "1",
				Minor:      "2",
				Patch:      "3",
				PreRelease: "rc0",
				Build:      "build123",
			},
		},
		{
			desc:     "Testing version major.minor.patch-prerelease",
			expected: "1.2.3-rc0",
			version: &SemVer{
				Major:      "1",
				Minor:      "2",
				Patch:      "3",
				PreRelease: "rc0",
			},
		},
		{
			desc:     "Testing version major.minor.patch+build",
			expected: "1.2.3+build123",
			version: &SemVer{
				Major: "1",
				Minor: "2",
				Patch: "3",
				Build: "build123",
			},
		},
		{
			desc:     "Testing version major.minor.patch",
			expected: "1.2.3",
			version: &SemVer{
				Major: "1",
				Minor: "2",
				Patch: "3",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			sv := test.version.String()
			assert.Equal(t, test.expected, sv)
		})
	}
}
