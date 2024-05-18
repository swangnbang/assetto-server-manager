package views

import (
	"embed"
)

// Pack frontend compiled assets into this package
//
// xxxgo:generate esc -o static_embed.go -pkg=views .

//go:embed **/***
var static embed.FS
