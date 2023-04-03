package main

import (
	"github.com/gabe565/relax-sounds/internal/server"
	flag "github.com/spf13/pflag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

const EnvPrefix = "RELAX_SOUNDS_"

func main() {
	address := flag.String("address", ":3000", "Override listen address.")
	frontendDir := flag.String("frontend", defaultFrontend, "Override frontend asset directory."+frontendHelpExt)
	dataDir := flag.String("data", "data", "Override data directory.")
	flag.Parse()

	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		optName := strings.ToUpper(f.Name)
		optName = strings.ReplaceAll(optName, "-", "_")
		varName := EnvPrefix + optName
		if val, ok := os.LookupEnv(varName); !f.Changed && ok {
			if err := f.Value.Set(val); err != nil {
				log.Fatalln(err)
			}
		}
	})

	var frontendFs fs.FS
	if *frontendDir != "" {
		frontendFs = os.DirFS(*frontendDir)
	} else {
		var err error
		frontendFs, err = fs.Sub(frontendEmbed, "frontend/dist")
		if err != nil {
			panic(err)
		}
	}

	server := &http.Server{
		Addr:    *address,
		Handler: server.Setup(frontendFs, os.DirFS(*dataDir)),
	}

	log.Println("Listening on " + *address)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
