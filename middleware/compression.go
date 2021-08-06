package middleware

import (
	"net/http"
)

type CompressionMiddleware struct {
	Next http.Handler
}

func (middleware *CompressionMiddleware) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	middleware.Next.ServeHTTP(w, r)
}

func NewCompressionMiddleware(next http.Handler) CompressionMiddleware {
	if next == nil {
		next = http.DefaultServeMux
	}
	return CompressionMiddleware{
		Next: next,
	}
}
