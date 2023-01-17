package main

import (
	"context"
	"errors"
	"github.com/gabe565/relax-sounds/internal/server"
	flag "github.com/spf13/pflag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const EnvPrefix = "RELAX_SOUNDS_"

func main() {
	var err error

	address := flag.String("address", ":3000", "Override listen address.")
	frontendDir := flag.String("frontend", defaultFrontend, "Override frontend asset directory."+frontendHelpExt)
	dataDir := flag.String("data", "data", "Override data directory.")
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

	var frontendFs fs.FS
	if *frontendDir != "" {
		frontendFs = os.DirFS(*frontendDir)
	} else {
		frontendFs, err = fs.Sub(frontendEmbed, "frontend/dist")
		if err != nil {
			panic(err)
		}
	}

	server := &http.Server{
		Addr:    *address,
		Handler: server.Setup(frontendFs, os.DirFS(*dataDir)),
	}

	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 60 seconds
		ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
		defer func() {
			cancelTimeout()
		}()

		// Trigger graceful shutdown
		log.Println("Performing graceful shutdown...")
		if err := server.Shutdown(ctx); err != nil {
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				log.Println("Graceful shutdown timed out")
			} else {
				log.Println(err)
			}
		}
		cancel()
	}()

	log.Println("Listening on " + *address)
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-ctx.Done()
	log.Println("Exiting")
}
