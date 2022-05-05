package middleware

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

var limiter *rate.Limiter

func init() {
	limiter = rate.NewLimiter(5, 10)
}

type LimiterMiddleware struct {
	Next http.Handler
}

func (lm *LimiterMiddleware) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	err := limiter.Wait(ctx)
	if err != nil {
		w.WriteHeader(http.StatusTooManyRequests)
	} else {
		lm.Next.ServeHTTP(w, r)
	}
}

func NewLimiterMiddleware(next http.Handler) LimiterMiddleware {
	if next == nil {
		next = http.DefaultServeMux
	}
	return LimiterMiddleware{
		Next: next,
	}
}
