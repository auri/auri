package ipaclient

import (
	"auri/logger"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type IpaClientTestSuite struct {
	suite.Suite
}

func (suite *IpaClientTestSuite) SetupTest() {
	New = NewMock
	require.NoError(suite.T(), logger.CreateFileLoggers())
}

func TestIpaClientTestSuite(t *testing.T) {
	suite.Run(t, new(IpaClientTestSuite))
}

func (suite *IpaClientTestSuite) Test_AddUser() {
	err := AddUser("ttestfoo", "test", "testfoo", "test@example.com", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAEjkPQpaaFHDTIdUJVx8jA4esihdfioaalmWXUlY0QK test@test.DE")
	require.NoError(suite.T(), err)
}

func (suite *IpaClientTestSuite) Test_LoginExists() {
	// user is present
	res, err := LoginExists("user-present-in-IPA")
	require.NoError(suite.T(), err)
	require.True(suite.T(), res)

	// user isn't present
	res, err = LoginExists("user-not-present-in-IPA")
	require.NoError(suite.T(), err)
	require.False(suite.T(), res)
}

func (suite *IpaClientTestSuite) Test_EMailExists() {
	// mail is present
	res, err := EMailExists("user-present-in-IPA@example.com")
	require.NoError(suite.T(), err)
	require.True(suite.T(), res)

	// mail doesn't exist
	res, err = EMailExists("user-not-present-in-IPA@example.com")
	require.NoError(suite.T(), err)
	require.False(suite.T(), res)
}

func (suite *IpaClientTestSuite) Test_UserFind() {
	// user is present
	login, err := UserFind("user-present-in-IPA@example.com")
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), "user-present-in-IPA", login)

	// user is present, but email is in uppercase
	login, err = UserFind("USER-present-in-IPA@example.com")
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), "user-present-in-IPA", login)

	// user doesn't exist
	login, err = UserFind("user-not-present-in-IPA@example.com")
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), "", login)
}

func (suite *IpaClientTestSuite) Test_GetGroups() {
	// user is present
	groups, err := GetGroups("user-present-in-IPA")
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), []string{"group1", "group2", "indirectgroup1", "indirectgroup2"}, groups)

	// user doesn't exist
	_, err = GetGroups("user-not-present-in-IPA")
	require.Error(suite.T(), err)
}

func (suite *IpaClientTestSuite) Test_UserModify() {
	err := UserModify("user-present-in-IPA@example.com", "*", "")
	require.NoError(suite.T(), err)
}
