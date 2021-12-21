//go:generate npm run build

package main

import (
	"flag"
	"github.com/gabe565/relax-sounds/internal/server"
	"log"
	"net/http"
)

type StreamMixer struct {
	files []uint8
}

func main() {
	var err error

	address := flag.String("address", ":3000", "Override listen address.")
	flag.Parse()

	router := server.Setup()

	log.Println("Listening on " + *address)
	err = http.ListenAndServe(*address, router)
	if err != nil {
		log.Fatalln(err)
	}
}
