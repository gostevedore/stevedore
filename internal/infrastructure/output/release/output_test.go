package release

// import (
// 	"bytes"
// 	"fmt"
// 	"io"
// 	"runtime"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestNewRelease(t *testing.T) {
// 	r := NewRelease(nil)
// 	r.Header = ""

// 	res := &Release{
// 		Header: "",
// 		OsArch: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
// 	}

// 	assert.Equal(t, res, r, "Unexpected release")
// }

// func TestPrintVersion(t *testing.T) {
// 	var w bytes.Buffer

// 	r := &Release{
// 		Header:    "HEADER",
// 		BuildDate: "Thu Mar  3 23:05:25 2005",
// 		Version:   "0.0",
// 		OsArch:    "linux/amd64",
// 		Commit:    "asdfqwer",
// 		Writer:    io.Writer(&w),
// 	}

// 	res := `Stevedore 0.0 Commit: asdfqwer linux/amd64 BuildDate: Thu Mar  3 23:05:25 2005
// `
// 	r.PrintVersion()
// 	assert.Equal(t, res, w.String(), "Unexpected release output")
// }
