package release

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	Version, Commit, BuildDate, OsArch string
)

// Release
type Release struct {
	BuildDate string
	Commit    string
	OsArch    string
	Version   string
}

// NewRelease
func NewRelease() *Release {
	return &Release{
		BuildDate: strings.TrimSpace(BuildDate),
		Version:   strings.TrimSpace(Version),
		Commit:    strings.TrimSpace(Commit),
		OsArch:    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
