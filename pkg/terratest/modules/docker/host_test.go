package docker

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDockerHostFromEnv(t *testing.T) {

	t.Parallel()

	tests := []struct {
		Input    string
		Expected string
	}{
		{
			"unix:///var/run/docker.sock",
			"localhost",
		},
		{
			"npipe:////./pipe/docker_engine",
			"localhost",
		},
		{
			"tcp://1.2.3.4:1234",
			"1.2.3.4",
		},
		{
			"tcp://1.2.3.4",
			"1.2.3.4",
		},
		{
			"ssh://1.2.3.4:22",
			"1.2.3.4",
		},
		{
			"fd://1.2.3.4:1234",
			"1.2.3.4",
		},
		{
			"",
			"localhost",
		},
		{
			"invalidValue",
			"localhost",
		},
		{
			"invalid::value::with::semicolons",
			"localhost",
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("GetDockerHostFromEnv: %s", test.Input), func(t *testing.T) {
			t.Parallel()

			testEnv := []string{
				"FOO=bar",
				fmt.Sprintf("DOCKER_HOST=%s", test.Input),
				"BAR=baz",
			}

			host := getDockerHostFromEnv(testEnv)
			assert.Equal(t, test.Expected, host)
		})
	}

	t.Run("GetDockerHostFromEnv: DOCKER_HOST unset", func(t *testing.T) {
		t.Parallel()

		testEnv := []string{
			"FOO=bar",
			"BAR=baz",
		}

		host := getDockerHostFromEnv(testEnv)
		assert.Equal(t, "localhost", host)
	})
}
