package models

import (
	"auri/ipaclient"
	"auri/logger"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/suite/v3"
)

type ModelSuite struct {
	*suite.Model
}

func Test_ModelSuite(t *testing.T) {
	ipaclient.New = ipaclient.NewMock
	require.NoError(t, logger.CreateFileLoggers())
	ConnectDB("test")

	model, err := suite.NewModelWithFixtures(packr.New("app:models:test:fixtures", "../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	as := &ModelSuite{
		Model: model,
	}
	suite.Run(t, as)
}
