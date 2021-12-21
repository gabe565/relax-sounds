//go:generate npm run build

package main

import (
	"github.com/gabe565/relax-sounds/internal/server"
	"log"
	"net/http"
)

type StreamMixer struct {
	files []uint8
}

func main() {
	router := server.Setup()

	addr := "0.0.0.0:3000"
	log.Println("Listening on " + addr)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatalln(err)
	}
}
