package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"runtime"
	"time"

	"github.com/x-dvr/go-service-template/internal/logging"
)

func NewRequestLogger(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ctx := withRequestID(r.Context(), r.Header.Get("x-request-id"))
			r = r.WithContext(ctx)

			rle := requestLogEntry{
				method: r.Method,
				path:   r.URL.EscapedPath(),
			}
			logger.LogAttrs(ctx, slog.LevelInfo, "-> HTTP", slog.Any("request", rle))

			writer := &logResponseWriter{w: w}
			next.ServeHTTP(writer, r)
			rle.elapsed = int(time.Since(start).Microseconds())
			rle.status = writer.status

			logger.LogAttrs(ctx, slog.LevelInfo, "<- HTTP", slog.Any("request", rle))
		})
	}
}

func Recover(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					ctx := r.Context()
					stack := make([]byte, 4096)
					stack = stack[:runtime.Stack(stack, false)]
					rle := requestLogEntry{
						method: r.Method,
						path:   r.URL.EscapedPath(),
						status: http.StatusInternalServerError,
					}
					logger.LogAttrs(ctx, slog.LevelError, "HTTP Panic",
						slog.Any("request", rle),
						slog.String("stack", string(stack)),
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}

	return ""
}

func withRequestID(ctx context.Context, requestID string) context.Context {
	if requestID == "" {
		return ctx
	}

	ctx = context.WithValue(ctx, requestIDKey, requestID)
	ctx = logging.AppendCtx(ctx, slog.String("req_id", requestID))

	return ctx
}

type ctxKey string

const (
	requestIDKey ctxKey = "reqID"
)

type requestLogEntry struct {
	method  string
	path    string
	status  int
	elapsed int
}

func (r requestLogEntry) LogValue() slog.Value {
	attrs := make([]slog.Attr, 0, 10)
	attrs = append(attrs,
		slog.String("method", r.method),
		slog.String("path", r.path),
	)

	if r.status != 0 {
		attrs = append(attrs, slog.Int("status", r.status))
	}

	if r.elapsed != 0 {
		attrs = append(attrs, slog.Int("time_taken_mu", r.elapsed))
	}

	return slog.GroupValue(attrs...)
}

type logResponseWriter struct {
	w      http.ResponseWriter
	status int
}

func (lrw *logResponseWriter) Header() http.Header {
	return lrw.w.Header()
}

func (lrw *logResponseWriter) Write(b []byte) (int, error) {
	return lrw.w.Write(b)
}

func (lrw *logResponseWriter) WriteHeader(h int) {
	lrw.status = h
	lrw.w.WriteHeader(h)
}
