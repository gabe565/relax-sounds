package playlist

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func DecoderMiddleware(next http.Handler) http.Handler {
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

		playlist, err := entries.ToPlaylist()
		if err != nil {
			http.Error(res, http.StatusText(400), 400)
			return
		}
		playlist.Name = chi.URLParam(req, "name")

		ctx := context.WithValue(req.Context(), "playlist", playlist)

		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
