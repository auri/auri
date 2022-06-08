package templates

import (
	"embed"
	"io/fs"

	"github.com/gobuffalo/buffalo"
)

//go:embed *.html
var files embed.FS

func FS() fs.FS { //revive:disable-line
	return buffalo.NewFS(files, "mail")
}
