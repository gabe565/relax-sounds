package server

import (
	"github.com/gabe565/relax-sounds/internal/mixer"
	"github.com/gabe565/relax-sounds/internal/playlist"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

func Setup(staticFs, dataFs fs.FS) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Heartbeat("/api/health"))
	router.Use(middleware.Logger)
	router.Use(middleware.Compress(5, "text/html", "text/css", "application/javascript", "application/json", "font/woff2"))
	router.Use(middleware.Recoverer)

	// Static Files
	staticserv := http.FileServer(http.FS(staticFs))
	router.Get("/*", fsPwaHandler(router, staticFs, staticserv))

	// Serve index as 404
	router.NotFound(func(res http.ResponseWriter, req *http.Request) {
		req.URL.Path = "/"
		staticserv.ServeHTTP(res, req)
	})

	// Data
	dataserv := http.FileServer(http.FS(dataFs))
	router.With(StripPrefix("/data")).Get("/data/*", fsPwaHandler(router, dataFs, dataserv))

	// Mixer
	router.With(playlist.DecoderMiddleware(dataFs)).Get("/api/mix/{enc}", mixer.Mix)

	return router
}

func fsPwaHandler(router *chi.Mux, filesystem fs.FS, h http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		stripped := strings.TrimLeft(req.URL.Path, string(os.PathSeparator))
		if _, err := fs.Stat(filesystem, stripped); !os.IsNotExist(err) {
			h.ServeHTTP(res, req)
		} else {
			router.NotFoundHandler().ServeHTTP(res, req)
		}
	}
}
