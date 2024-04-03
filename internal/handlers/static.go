package handlers

import (
	"os"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/spf13/cobra"
)

func StaticFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("public", "frontend/dist", "Public directory")
}

func StaticHandler(app *pocketbase.PocketBase) echo.HandlerFunc {
	publicDir, err := app.RootCmd.PersistentFlags().GetString("public")
	if err != nil {
		panic(err)
	}

	return apis.StaticDirectoryHandler(os.DirFS(publicDir), true)
}
