package models

import (
	"auri/ipaclient"
	"auri/logger"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gobuffalo/suite/v4"
)

type ModelSuite struct {
	*suite.Model
}

func Test_ModelSuite(t *testing.T) {
	ipaclient.New = ipaclient.NewMock
	require.NoError(t, logger.CreateFileLoggers())
	ConnectDB("test")

	model, err := suite.NewModelWithFixtures(os.DirFS("../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	as := &ModelSuite{
		Model: model,
	}
	suite.Run(t, as)
}
