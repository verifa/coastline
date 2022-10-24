//go:build !oapi

package server

import (
	"embed"
	"net/http"
)

var oapiIndex embed.FS
var oapiEnabled bool = false

var oapiSite http.FileSystem
