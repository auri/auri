package main

import (
	"auri/actions"
	"auri/config"

	"fmt"
	"log"
	"os"
)

const helpInfo = `NAME:
  auri - Automated User Registration IPA

USAGE:
  auri [global options] [command]

VERSION:
   %s

AUTHORS:
  Mohamad Asswad <Mohamad.Asswad@telekom.de>
  Daniel Ajbassow <Daniel.Ajbassow@telekom.de>
  Artem Sidorenko <artem.sidorenko@t-systems.com>
  Sergej Schischkowki <sergej.schischkowski@t-systems.com>

COMMANDS:
  task              Display or execute a task, see below for details
  migrate           Run database migrations
  version           Show the version information
  help              Show the help information

  when command is missing, auri is started as a daemon
                 
GLOBAL OPTIONS:  
  --help            Show the help information
  --configdir DIR   Use DIR as configuration directory
                 
TASKS:           
  task list         Show the available tasks
  
`

// main is the starting point of this Buffalo application.
func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	app := actions.App()

	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}

/*
  dirty hack to have own CLI flags
  see upstream issue:
  https://github.com/gobuffalo/buffalo/issues/2070#issuecomment-741806786
*/
func init() {
	if len(os.Args) > 1 {
		if os.Args[1] == "--help" || os.Args[1] == "help" {
			fmt.Printf(helpInfo, config.GetVersion())
			os.Exit(0)
		} else if os.Args[1] == "--configdir" {
			if len(os.Args) < 2 {
				log.Fatal("Syntax error: --configdir requires an argument")
			}
			config.ConfigDir = os.Args[2]

			// we have to get rid of config args from os.Args
			// as buffalo has no idea how to handle that
			os.Args = append([]string{os.Args[0]}, os.Args[3:]...)
		}
	}
}
