package middleware

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	var log = DefaultLogger(context.Background())
	assert.NotNil(t, log)
	assert.Implements(t, (*Logger)(nil), log)
	assert.IsType(t, new(slog.Logger), log)
	log.Info("print info log message")
}

func TestTestLogger(t *testing.T) {
	var log = new(TestLogger)
	assert.Implements(t, (*Logger)(nil), log)
	log.Info("print info message")
	log.Error("print error message")
	output := log.String()
	assert.NotEmpty(t, output)
	assert.Contains(t, output, "print info message")
	assert.Contains(t, output, "ERROR")
}
