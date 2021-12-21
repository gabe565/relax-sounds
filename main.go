//go:generate npm run build

package main

import (
	"embed"
	"flag"
	"github.com/gabe565/relax-sounds/internal/server"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed dist
var dist embed.FS

func main() {
	var err error

	address := flag.String("address", ":3000", "Override listen address.")
	staticDir := flag.String("static", "", "Override static asset directory. Useful for development. If left empty, embedded assets are used.")
	flag.Parse()

	var contentFs fs.FS
	if *staticDir != "" {
		contentFs = os.DirFS(*staticDir)
	} else {
		contentFs, err = fs.Sub(dist, "dist")
		if err != nil {
			panic(err)
		}
	}
	router := server.Setup(contentFs)

	log.Println("Listening on " + *address)
	err = http.ListenAndServe(*address, router)
	if err != nil {
		log.Fatalln(err)
	}
}
