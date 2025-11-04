package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func MiddleLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("Component", "middleware/Logger"),
		)
		log.Info("Starting HTTP request logger")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("Method", r.Method),
				slog.String("Path", r.URL.Path),
				slog.String("Remote_addr", r.RemoteAddr),
				slog.String("User_agent", r.UserAgent()),
				slog.String("Request_id", middleware.GetReqID(r.Context())),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				entry.Info("Request completed",
					slog.Int("Status", ww.Status()),
					slog.Int("Bytes", ww.BytesWritten()),
					slog.String("Duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
