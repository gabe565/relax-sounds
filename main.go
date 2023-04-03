package main

import (
	"context"
	"errors"
	"github.com/gabe565/relax-sounds/internal/server"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
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

	var group errgroup.Group
	group.Go(func() error {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sig

		// Shutdown signal with grace period of 60 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// Trigger graceful shutdown
		log.Info("Performing graceful shutdown...")
		return server.Shutdown(ctx)
	})

	group.Go(func() error {
		log.WithField("address", *address).Info("Listening")
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	})

	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	log.Info("Exiting")
}
