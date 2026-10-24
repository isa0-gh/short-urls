package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// LoggingMiddleware logs one structured record for each completed request.
func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			response := &loggedResponseWriter{ResponseWriter: w}

			defer func() {
				panicked := recover()
				status := response.status
				if status == 0 && panicked == nil {
					status = http.StatusOK
				}

				level := slog.LevelInfo
				if panicked != nil || status >= http.StatusInternalServerError {
					level = slog.LevelError
				} else if status >= http.StatusBadRequest {
					level = slog.LevelWarn
				}

				attrs := []slog.Attr{
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Int("status", status),
					slog.Int("bytes", response.bytesWritten),
					slog.Duration("duration", time.Since(started)),
				}
				if panicked != nil {
					attrs = append(attrs, slog.Bool("panicked", true))
				}
				logger.LogAttrs(r.Context(), level, "http request", attrs...)
				if panicked != nil {
					panic(panicked)
				}
			}()

			next.ServeHTTP(response, r)
		})
	}
}

type loggedResponseWriter struct {
	http.ResponseWriter
	status       int
	bytesWritten int
}

func (w *loggedResponseWriter) WriteHeader(status int) {
	if status >= 100 && status < 200 && status != http.StatusSwitchingProtocols {
		w.ResponseWriter.WriteHeader(status)
		return
	}
	if w.status != 0 {
		return
	}
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *loggedResponseWriter) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(body)
	w.bytesWritten += n
	return n, err
}

func (w *loggedResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
