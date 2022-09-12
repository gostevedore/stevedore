package dryrun

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/gostevedore/stevedore/internal/core/domain/image"
)

// DryRunDriver is a driver that just simulates the build process
type DryRunDriver struct {
	write io.Writer
}

// NewDryRunDriver creates a new DryRunDriver
func NewDryRunDriver(w io.Writer) *DryRunDriver {

	if w == nil {
		w = os.Stdout
	}

	return &DryRunDriver{
		write: w,
	}
}

// Build simulate a new image build
func (d *DryRunDriver) Build(ctx context.Context, i *image.Image, options *image.BuildDriverOptions) error {
	fmt.Fprintln(d.write)
	fmt.Fprintln(d.write, fmt.Sprintf(" builder:	%+v", i.Builder))
	if len(i.Children) > 0 {
		fmt.Fprintln(d.write, " children:")
		for _, child := range i.Children {
			fmt.Fprintln(d.write, fmt.Sprintf(" - %s:%s", child.Name, child.Version))
		}
	}
	fmt.Fprintln(d.write, fmt.Sprintf(" lables: %+v", i.Labels))
	fmt.Fprintln(d.write, fmt.Sprintf(" name: %s", i.Name))

	if i.Parent != nil {
		fmt.Fprintln(d.write, " parent:")
		fmt.Fprintln(d.write, fmt.Sprintf(" - %s:%s", i.Parent.Name, i.Parent.Version))
	}
	fmt.Fprintln(d.write, fmt.Sprintf(" presistent labels: %+v", i.PersistentLabels))
	fmt.Fprintln(d.write, fmt.Sprintf(" presistent vars: %+v", i.PersistentVars))
	fmt.Fprintln(d.write, fmt.Sprintf(" registry host: %s", i.RegistryHost))
	fmt.Fprintln(d.write, fmt.Sprintf(" registry namespace: %s", i.RegistryNamespace))
	fmt.Fprintln(d.write, fmt.Sprintf(" tags: %v", i.Tags))
	fmt.Fprintln(d.write, fmt.Sprintf(" vars: %v", i.Vars))
	fmt.Fprintln(d.write, fmt.Sprintf(" version: %v", i.Version))
	fmt.Fprintln(d.write, fmt.Sprintf(" options: %+v", *options))
	fmt.Fprintln(d.write)

	return nil
}
