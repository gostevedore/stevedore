package helpers

import "github.com/gruntwork-io/terratest/modules/testing"

type QuiteLogger struct{}

func (l *QuiteLogger) Logf(t testing.TestingT, format string, args ...interface{}) {}
