//go:build !ui

package ui

import (
	"embed"
	"net/http"
)

var Build embed.FS
var Enabled bool = false

var Site http.FileSystem
