//go:generate npm run build

package main

import (
	"github.com/gabe565/relax-sounds/internal/server"
	flag "github.com/spf13/pflag"
	"log"
	"net/http"
	"os"
)

func main() {
	var err error

	address := flag.String("address", ":3000", "Override listen address.")
	staticDir := flag.String("static", "", "Override static asset directory. Useful for development. If left empty, embedded assets are used.")
	flag.Parse()

	contentFs := os.DirFS(*staticDir)
	router := server.Setup(contentFs)

	log.Println("Listening on " + *address)
	err = http.ListenAndServe(*address, router)
	if err != nil {
		log.Fatalln(err)
	}
}
