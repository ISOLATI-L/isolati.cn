package middleware

import (
	"context"
	"net/http"
	"time"
)

type TimeoutMiddleware struct {
	Next http.Handler
}

func (tm *TimeoutMiddleware) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	r = r.Clone(ctx)
	ch := make(chan struct{})
	go func() {
		// time.Sleep(6 * time.Second)
		tm.Next.ServeHTTP(w, r)
		ch <- struct{}{}
	}()
	select {
	case <-ch:
		return
	case <-ctx.Done():
		w.WriteHeader(http.StatusRequestTimeout)
	}
	ctx.Done()
}

func NewTimeoutMiddleware(next http.Handler) TimeoutMiddleware {
	if next == nil {
		next = http.DefaultServeMux
	}
	return TimeoutMiddleware{
		Next: next,
	}
}
