package logging

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
)

// Chi attaches a request-scoped *slog.Logger to r.Context() and logs one line per request when it completes.
func Chi(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()

		l := slog.Default().With(
			slog.String("http.request.method", r.Method),
			slog.String("http.request.path", r.URL.Path),
		)
		if id := chimw.GetReqID(r.Context()); id != "" {
			l = l.With(slog.String("http.request.id", id))
		}

		ctx := With(r.Context(), l)
		r = r.WithContext(ctx)

		next.ServeHTTP(ww, r)

		l.InfoContext(ctx, "http.request.complete",
			slog.Int("http.response.status_code", ww.Status()),
			slog.Duration("http.server.request.duration", time.Since(start)),
		)
	})
}

// Recover logs panics with slog (including stack) and returns HTTP 500.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if p := recover(); p != nil {
				From(r.Context()).ErrorContext(r.Context(), "panic recovered",
					slog.Any("panic", p),
					slog.String("stack", string(debug.Stack())),
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
