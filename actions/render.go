package actions

import (
	"auri/config"
	"auri/public"
	"auri/templates"

	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/buffalo/render"
)

var r *render.Engine

var customAssetsBox *packr.Box

// DoInit for manual initialization
func DoInit() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.plush.html",

		// Box containing all of the templates:
		TemplatesFS: templates.FS(),
		AssetsFS:    public.FS(),

		// Add template helpers here:
		Helpers: render.Helpers{
			// for non-bootstrap form helpers uncomment the lines
			// below and import "github.com/gobuffalo/helpers/forms"
			// forms.FormKey:     forms.Form,
			// forms.FormForKey:  forms.FormFor,
		},
	})
	// custom assets if configured
	if config.GetInstance().WebAssetPath != "" {
		customAssetsBox = packr.New("custom:assets", config.GetInstance().WebAssetPath)
	}
}
