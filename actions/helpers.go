package actions

import (
	"auri/config"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"

	"github.com/gobuffalo/buffalo/render"
)

//tokenExpiredPageRenderer returns the renderer of page with expired hint
func tokenExpiredPageRenderer() render.Renderer {
	return r.HTML("token-expired.plush.html", config.GetInstance().Layout)
}

//getDBConnection returns the initialized DB connection from our context
func getDBConnection(c buffalo.Context) (*pop.Connection, error) {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return nil, errors.New("no DB transaction found")
	}
	return tx, nil
}
