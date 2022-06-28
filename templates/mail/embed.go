package mail

import (
	"auri/config"

	"embed"
	"io/fs"
	"os"

	"github.com/gobuffalo/buffalo"
)

//go:embed *.html
var files embed.FS

func FS() fs.FS { //revive:disable-line
	// use own templates if defined
	if config.GetInstance().EmailTemplatePath != "" {
		return os.DirFS(config.GetInstance().EmailTemplatePath)
	}

	return buffalo.NewFS(files, "mail")
}
