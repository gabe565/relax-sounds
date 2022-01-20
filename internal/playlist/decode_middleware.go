package playlist

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io/fs"
	"net/http"
)

type ContextKey string

const RequestKey = ContextKey("playlist")

func DecoderMiddleware(dataFs fs.FS) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			data, err := base64.RawURLEncoding.DecodeString(
				chi.URLParam(req, "enc"),
			)
			if err != nil {
				http.Error(res, http.StatusText(400), 400)
				return
			}

			var entries PlaylistShorthand
			if err = json.Unmarshal(data, &entries); err != nil {
				http.Error(res, http.StatusText(400), 400)
				return
			}

			playlist, err := entries.ToPlaylist(dataFs)
			if err != nil {
				http.Error(res, http.StatusText(400), 400)
				return
			}

			ctx := context.WithValue(req.Context(), RequestKey, playlist)

			next.ServeHTTP(res, req.WithContext(ctx))
		})
	}
}
