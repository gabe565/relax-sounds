package server

import (
	"net/http"
	"strings"
)

func StripPrefix(prefix string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
