package logging

import (
	"context"
	"log/slog"
)

type loggerKey struct{}

// With returns ctx carrying logger l. Handlers should use From(ctx) for request-scoped fields.
func With(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

// From returns the logger stored on ctx, or slog.Default() if none was attached.
func From(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok && l != nil {
		return l
	}
	return slog.Default()
}
