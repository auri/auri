// Package public contains embedded assets of our application
package public

import (
	"embed"
	"io/fs"

	"github.com/gobuffalo/buffalo"
)

//go:embed *
var files embed.FS

func FS() fs.FS { //revive:disable-line
	return buffalo.NewFS(files, "public")
}
