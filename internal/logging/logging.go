package logging

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

type contextKey struct{}

var requestIDKey = contextKey{}

type handler struct {
	inner slog.Handler
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	var reqID string
	if v, ok := ctx.Value(requestIDKey).(string); ok && v != "" {
		reqID = v
	} else if lc, ok := lambdacontext.FromContext(ctx); ok && lc != nil {
		reqID = lc.AwsRequestID
	}
	if reqID != "" {
		r.AddAttrs(slog.String("request_id", reqID))
	}
	return h.inner.Handle(ctx, r)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{inner: h.inner.WithAttrs(attrs)}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{inner: h.inner.WithGroup(name)}
}

func Init() {
	SetOutput(os.Stdout)
}

func SetOutput(w io.Writer) {
	jsonHandler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(&handler{inner: jsonHandler}))
}

func WithRequestID(ctx context.Context, reqID string) context.Context {
	if reqID == "" {
		return ctx
	}
	return context.WithValue(ctx, requestIDKey, reqID)
}
