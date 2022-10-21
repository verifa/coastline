//go:build !oapi

package server

import (
	"embed"
	"net/http"
)

//go:embed oapi/index.html
var oapiIndex embed.FS
var oapiEnabled bool = false

var oapiSite http.FileSystem
