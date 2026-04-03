package frontend

import (
	"embed"
	"io/fs"
)

//go:generate pnpm install
//go:generate pnpm run build

//go:embed dist
var dist embed.FS

func FS() (fs.FS, error) {
	return fs.Sub(dist, "dist")
}
