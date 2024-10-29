package handlers

import (
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

func StaticFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("public", "frontend/dist", "Public directory")
}

func StaticHandler(app *pocketbase.PocketBase) func(*core.RequestEvent) error {
	publicDir, err := app.RootCmd.PersistentFlags().GetString("public")
	if err != nil {
		panic(err)
	}

	return apis.Static(os.DirFS(publicDir), true)
}
