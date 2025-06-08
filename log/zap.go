package log

import (
	"context"
	"net/http"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey struct{}

var (
	once      sync.Once
	startTime = time.Now()
)

func Save(ctx context.Context, z *zap.Logger) context.Context {
	ctx = context.WithValue(ctx, contextKey{}, z)
	return ctx
}

func Load(ctx context.Context) *zap.Logger {
	logger := zap.L() // default to global logger
	if l, ok := ctx.Value(contextKey{}).(*zap.Logger); ok {
		logger = l
	}
	return logger
}

func With(ctx context.Context, fields ...zapcore.Field) (context.Context, *zap.Logger) {
	z := Load(ctx).With(fields...)
	return Save(ctx, z), z
}

func Init(level string) *zap.Logger {
	z := zap.L()
	once.Do(func() {
		// FIXME: Add production config
		atom := zap.NewAtomicLevelAt(zap.DebugLevel)
		consoleErrors := zapcore.Lock(os.Stderr)
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core := zapcore.NewCore(consoleEncoder, consoleErrors, atom)
		z = zap.New(core)
		zap.ReplaceGlobals(z)
		http.Handle("/loglevel", atom)
	})
	return z
}
