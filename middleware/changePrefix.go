package middleware

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type ChangePrefixMiddleware struct {
	Next                http.Handler
	changePrefixPattern *regexp.Regexp
	domain              string
	newPrefix           string
}

func (cpm *ChangePrefixMiddleware) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	host := r.Host
	matches := cpm.changePrefixPattern.FindStringSubmatch(host)
	if len(matches) > 0 {
		newURL := fmt.Sprintf("https://%s%s%s", cpm.newPrefix, cpm.domain, r.URL)
		http.Redirect(w, r, newURL, http.StatusFound)
	} else {
		cpm.Next.ServeHTTP(w, r)
	}
}

func NewChangePrefixMiddleware(
	next http.Handler,
	domain string,
	prefix string,
	newPrefix string,
) ChangePrefixMiddleware {
	if next == nil {
		next = http.DefaultServeMux
	}
	patternStr := fmt.Sprintf(`^(%s)%s`, prefix, domain)
	changePrefixPattern, err := regexp.Compile(patternStr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return ChangePrefixMiddleware{
		Next:                next,
		changePrefixPattern: changePrefixPattern,
		domain:              domain,
		newPrefix:           newPrefix,
	}
}
