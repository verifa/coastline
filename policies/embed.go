package policies

import (
	"embed"
)

//go:embed all:*.rego
var Policies embed.FS
