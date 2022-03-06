package preset

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io/fs"
	"net/http"
)

type ContextKey string

const RequestKey = ContextKey("preset")

func DecoderMiddleware(dataFs fs.FS) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			data, err := base64.RawURLEncoding.DecodeString(
				chi.URLParam(req, "enc"),
			)
			if err != nil {
				http.Error(res, http.StatusText(400), 400)
				panic(err)
			}

			var entries PresetShorthand
			if err = json.Unmarshal(data, &entries); err != nil {
				http.Error(res, http.StatusText(400), 400)
				panic(err)
			}

			preset, err := entries.ToPreset(dataFs)
			if err != nil {
				http.Error(res, http.StatusText(400), 400)
				panic(err)
			}

			ctx := context.WithValue(req.Context(), RequestKey, preset)

			next.ServeHTTP(res, req.WithContext(ctx))
		})
	}
}
