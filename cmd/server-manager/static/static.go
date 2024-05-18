package static

import (
	"embed"
	"net/http"
)

// Pack frontend compiled assets into this package
//xxxgo:generate esc -o static_embed.go -pkg=static -ignore=static.go .

//go:embed *
var static embed.FS

func FS(_ bool) http.FileSystem {
	return http.FS(static)
}

func FSMustByte(name string) []byte {
	if b, err := static.ReadFile(name); err != nil {
		panic(err)
	} else {
		return b
	}
}
