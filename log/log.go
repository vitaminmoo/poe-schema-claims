package log

import (
	"net/http"
	"time"

	"github.com/urfave/negroni"
	"github.com/vitaminmoo/poe-schema-claims/ctxutil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct{}

func New() *logger {
	return &logger{}
}

func (l *logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := r.Context()
	start := time.Now()

	next(rw, r)

	res := rw.(negroni.ResponseWriter)

	responseSize := int64(res.Size())
	for hKey, hValue := range res.Header() {
		responseSize += int64(len([]byte(hKey)))
		for _, v := range hValue {
			responseSize += int64(len([]byte(v)))
		}
	}

	zapFields := []zap.Field{
		zap.String("host", r.URL.Host),
		zap.String("pattern", r.Pattern),
		zap.Int("code", res.Status()),
		zap.Duration("duration", time.Since(start)),
		zap.Int64("response_size", responseSize),
	}

	errPtr := ctxutil.GetErrorPtr(ctx)
	if errPtr != nil {
		zapFields = append(zapFields, zap.NamedError("error", *errPtr))
	}

	if ctx.Err() != nil {
		zapFields = append(zapFields, zap.NamedError("context_err", ctx.Err()))
	}

	z := Load(ctx)

	level := zapcore.InfoLevel
	if res.Status() >= 500 {
		level = zapcore.ErrorLevel
	} else if res.Status() >= 400 {
		level = zapcore.WarnLevel
	} else if r.URL.Path == "/healthz" {
		return
	}

	z.Log(level, "http response", zapFields...)
}
