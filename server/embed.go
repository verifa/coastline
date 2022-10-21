//go:build oapi

package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed oapi/index.html
var oapiIndex embed.FS
var oapiEnabled bool = true

var oapiSite http.FileSystem = oapiSubDir()

func oapiSubDir() http.FileSystem {
	subDir, err := fs.Sub(oapiIndex, "oapi")
	if err != nil {
		log.Fatal("directory does not exist in oapi filesystem: oapi")
	}
	return http.FS(subDir)
}
