package middleware

import (
	"fmt"
	"net/http"
)

type RedirectMiddleware struct {
	newURL string
}

func (rm *RedirectMiddleware) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	newURL := fmt.Sprintf("%s%s", rm.newURL, r.URL.Path)
	http.Redirect(w, r, newURL, http.StatusFound)
}

func NewRedirectMiddleware(newURL string) RedirectMiddleware {
	return RedirectMiddleware{
		newURL: newURL,
	}
}
