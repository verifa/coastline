//go:build ui

package ui

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed all:build
var Build embed.FS
var Enabled bool = true

var Site http.FileSystem = site()

func site() http.FileSystem {
	subDir, err := fs.Sub(Build, "build")
	if err != nil {
		log.Fatal("directory does not exist in ui filesystem: build")
	}
	return http.FS(subDir)
}
