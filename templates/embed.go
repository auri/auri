package templates

import (
	"auri/config"

	"embed"
	"io/fs"
	"os"

	"github.com/gobuffalo/buffalo"
)

//go:embed *.html mail/*.html
var files embed.FS

func FS() fs.FS { //revive:disable-line
	// use own templates if defined
	if config.GetInstance().WebTemplatePath != "" {
		return os.DirFS(config.GetInstance().WebTemplatePath)
	}

	return buffalo.NewFS(files, "templates")
}
