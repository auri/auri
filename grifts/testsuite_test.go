package grifts

import (
	"auri/ipaclient"
	"auri/mail"
	"os"
	"testing"

	"github.com/gobuffalo/suite/v4"
)

type TaskSuite struct {
	*suite.Model
}

func Test_TaskSuite(t *testing.T) {
	mail.New = mail.NewMock
	ipaclient.New = ipaclient.NewMock

	model, err := suite.NewModelWithFixtures(os.DirFS("../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	ts := &TaskSuite{
		Model: model,
	}

	suite.Run(t, ts)
}
