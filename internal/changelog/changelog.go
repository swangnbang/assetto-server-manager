package changelog

import (
	"html/template"

	"github.com/russross/blackfriday"
)

// Pack changelog into this package
// xxxgo:generate esc -o changelog_embed.go -pkg=changelog ../../CHANGELOG.md

// go:embed CHANGELOG.md
var changelog []byte

func LoadChangelog() (template.HTML, error) {
	return template.HTML(blackfriday.Run(changelog)), nil
}
