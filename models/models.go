package models

import (
	"auri/config"

	"github.com/gobuffalo/pop/v6"
)

// DB is a connection to your database to be used
// throughout your application.
var DB *pop.Connection

// ConnectDB initialized the DB connection using given environment
func ConnectDB(env string) error {
	if config.IsProdBuild() {
		pop.AddLookupPaths(config.ConfigDir)
		pop.LoadConfigFile()
	}

	var err error
	DB, err = pop.Connect(env)
	if err != nil {
		return err
	}
	pop.Debug = env == "development"

	return nil
}
