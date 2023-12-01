package functional

import "github.com/gruntwork-io/terratest/modules/testing"

type quiteLogger struct{}

func (l *quiteLogger) Logf(t testing.TestingT, format string, args ...interface{}) {}
