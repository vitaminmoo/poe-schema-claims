package ctxutil

import (
	"context"
	"net/http"
)

type errorKey struct{}

func GetErrorPtr(ctx context.Context) *error {
	errorPtr := ctx.Value(errorKey{})
	return errorPtr.(*error)
}

type mid struct {
	mockedEmail string
}

func NewMiddleware() *mid {
	return &mid{}
}

func (l *mid) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	r = r.WithContext(l.Context(r.Context()))
	next(rw, r)
}

func (l *mid) Context(ctx context.Context) context.Context {
	var errorPtr error
	ctx = context.WithValue(ctx, errorKey{}, &errorPtr)
	return ctx
}
