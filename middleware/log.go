package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

type (
	Logger interface {
		Info(string, ...any)
		Error(string, ...any)
	}

	UseLogger func(ctx context.Context) Logger
)

func DefaultLogger(_ context.Context) Logger {
	return slog.Default()
}

type TestLogger struct {
	strings.Builder
}

func (tl *TestLogger) Info(msg string, args ...any) {
	tl.WriteString(fmt.Sprintf("\nINFO %s: %+v", msg, args))
}

func (tl *TestLogger) Error(msg string, args ...any) {
	tl.WriteString(fmt.Sprintf("\nERROR %s: %+v", msg, args))
}
