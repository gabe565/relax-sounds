package server

import (
	"github.com/gabe565/relax-sounds/internal/mixer"
	"github.com/gabe565/relax-sounds/internal/playlist"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	flag "github.com/spf13/pflag"
	"net/http"
	"os"
	"path/filepath"
)

type Server struct {
	mux *chi.Mux
}

const Public = "dist"

var staticDir string

func init() {
	flag.StringVar(&staticDir, "static", "dist", "Override static asset directory. Useful for development.")
}

func Setup() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Heartbeat("/health"))
	router.Use(middleware.Logger)
	router.Use(middleware.Compress(5, "text/html", "text/css", "application/javascript", "application/json", "font/woff2"))
	router.Use(middleware.Recoverer)

	// Serve index as 404
	router.NotFound(func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, Public+"/index.html")
	})

	// Static Files
	fileserver := http.FileServer(http.Dir(staticDir))
	router.Get("/*", func(res http.ResponseWriter, req *http.Request) {
		requestPath := filepath.Join(Public, filepath.Clean("/"+req.URL.Path))
		if _, err := os.Stat(requestPath); !os.IsNotExist(err) {
			fileserver.ServeHTTP(res, req)
		} else {
			router.NotFoundHandler().ServeHTTP(res, req)
		}
	})

	// Mixer
	router.With(playlist.DecoderMiddleware).Get("/mix/{enc}", mixer.Mix)

	return router
}
