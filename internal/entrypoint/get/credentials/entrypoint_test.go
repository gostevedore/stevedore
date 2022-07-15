package credentials

import (
	"context"
	"testing"

	"github.com/gostevedore/stevedore/internal/infrastructure/configuration"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		desc       string
		entrypoint *Entrypoint
		args       []string
		conf       *configuration.Configuration
		err        error
	}{}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			assert.Equal(t, test.err, test.entrypoint.Execute(context.TODO(), test.args, test.conf))
		})
	}
}
