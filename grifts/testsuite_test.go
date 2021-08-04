package grifts

import (
	"auri/ipaclient"
	"auri/mail"
	"testing"

	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/suite/v3"
)

type TaskSuite struct {
	*suite.Model
}

func Test_TaskSuite(t *testing.T) {
	mail.New = mail.NewMock
	ipaclient.New = ipaclient.NewMock

	model, err := suite.NewModelWithFixtures(packr.New("app:models:test:fixtures", "../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	ts := &TaskSuite{
		Model: model,
	}

	suite.Run(t, ts)
}
