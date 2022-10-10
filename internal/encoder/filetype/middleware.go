package filetype

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type ContextKey string

const RequestKey = ContextKey("filetype")

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := FileType(0)
		err := t.UnmarshalText([]byte(chi.URLParam(r, "filetype")))
		if err != nil {
			if errors.Is(err, ErrInvalidFileType) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			panic(err)
		}

		ctx := context.WithValue(r.Context(), RequestKey, t)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
