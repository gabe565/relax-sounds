//go:generate -command npm npm --prefix frontend
//go:generate npm install
//go:generate npm run build

package main

import (
	"embed"
	"github.com/gabe565/relax-sounds/internal/server"
	flag "github.com/spf13/pflag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

const EnvPrefix = "RELAX_SOUNDS_"

//go:embed frontend/dist
var dist embed.FS

func main() {
	var err error

	address := flag.String("address", ":3000", "Override listen address.")
	staticDir := flag.String("static", "dist", "Override static asset directory. Useful for development.")
	dataDir := flag.String("data", "data", "Override data directory. Useful for development.")
	flag.Parse()

	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		optName := strings.ToUpper(f.Name)
		optName = strings.ReplaceAll(optName, "-", "_")
		varName := EnvPrefix + optName
		if val, ok := os.LookupEnv(varName); !f.Changed && ok {
			err = f.Value.Set(val)
			if err != nil {
				log.Fatalln(err)
			}
		}
	})

	var contentFs fs.FS
	if *staticDir != "" {
		contentFs = os.DirFS(*staticDir)
	} else {
		contentFs, err = fs.Sub(dist, "frontend/dist")
		if err != nil {
			panic(err)
		}
	}
	router := server.Setup(contentFs, os.DirFS(*dataDir))

	log.Println("Listening on " + *address)
	err = http.ListenAndServe(*address, router)
	if err != nil {
		log.Fatalln(err)
	}
}
