// Package grifts implements buffalo grifts, jobs which might be executed via auri binary
package grifts

import (
	"auri/actions"

	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
