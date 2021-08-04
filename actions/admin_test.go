package actions

import (
	"auri/ipaclient"
	"auri/mail"
	"auri/models"
	"fmt"

	"github.com/tehwalris/go-freeipa/freeipa"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/logger"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type mockContext struct {
	buffalo.Context
	flash *buffalo.Flash
}

func (mockContext) Value(interface{}) interface{} {
	return nil
}

func (mc mockContext) Logger() buffalo.Logger {
	return logger.New(logrus.InfoLevel)
}

func (mc *mockContext) Flash() *buffalo.Flash {
	mc.flash = new(buffalo.Flash)
	mc.flash.Clear()
	return mc.flash
}

func (mc mockContext) Error(code int, err error) error {
	return buffalo.HTTPError{Status: code, Cause: errors.WithStack(err)}
}

func (as *ActionSuite) Test_GetRequestsListOk() {
	as.LoadFixture("handle requests")
	res := as.HTML("/admin").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "domain@example.com")
	as.Contains(body, "No")
}

func (as *ActionSuite) Test_GetAllRequestsNOkNoTX() {
	mc := &mockContext{}
	err := RequestResource{}.List(mc)
	as.Error(err)

	bError, ok := err.(buffalo.HTTPError)
	as.True(ok)
	as.Equal(500, bError.Status)
	as.Equal("An error has occurred. Please reload the page and try again", bError.Error())
}

func (as *ActionSuite) Test_GetRequestsListNotOk() {
	res := as.HTML("/admin").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.NotContains(body, "domain@example.de")
	as.NotContains(body, "Yes")
}

func (as *ActionSuite) Test_FindUserDestroyOk() {
	defer mail.ClearMock()
	as.LoadFixture("handle requests")
	request := models.Request{}
	as.NoError(as.DB.Eager().Find(&request, "f19d5ab3-e913-4c08-b4dd-627e26ccbb0c"))
	req := as.JSON(fmt.Sprintf("/admin/%s", request.ID))
	resp := req.Delete()
	as.Equal(302, resp.Code)
	as.Equal("/admin/", resp.Location())
	as.Error(as.DB.Eager().Find(&request, "f19d5ab3-e913-4c08-b4dd-627e26ccbb0c"))
}

func (as *ActionSuite) Test_FindUserDestroyNotOkNoUser() {
	req := as.JSON(fmt.Sprintf("/admin/%s", "f29d5ab3-e913-4c08-b4dd-627e26ccbb0c"))
	resp := req.Delete()
	as.Equal(400, resp.Code, resp.Body.String())
}

func (as *ActionSuite) Test_FindUserDestroyNOkNoTX() {
	mc := &mockContext{}
	err := RequestResource{}.Destroy(mc)
	as.Error(err)
	bError, ok := err.(buffalo.HTTPError)
	as.True(ok)
	as.Equal(500, bError.Status)
	as.Equal("An error has occurred. Please reload the page and try again", bError.Error())
}

func (as *ActionSuite) Test_FindUserUpdateOk() {
	defer mail.ClearMock()
	as.LoadFixture("handle requests")
	request := models.Request{}
	as.NoError(as.DB.Eager().Find(&request, "f19d5ab3-e913-4c08-b4dd-627e26ccbb0c"))
	req := as.JSON(fmt.Sprintf("/admin/%s", request.ID))
	resp := req.Put(&request)
	as.Equal(302, resp.Code)
	as.Equal("/admin/", resp.Location())
	as.Error(as.DB.Eager().Find(&request, "f19d5ab3-e913-4c08-b4dd-627e26ccbb0c"))
}

func (as *ActionSuite) Test_FindUserUpdateNotOkNoUser() {
	req := as.JSON(fmt.Sprintf("/admin/%s", "f29d5ab3-e913-4c08-b4dd-627e26ccbb0c"))
	resp := req.Put(models.Request{})
	as.Equal(400, resp.Code, resp.Body.String())
}

func (as *ActionSuite) Test_FindUserUpdateNotOkNoTX() {
	mc := &mockContext{}
	err := RequestResource{}.Update(mc)
	as.Error(err)
	bError, ok := err.(buffalo.HTTPError)
	as.True(ok)
	as.Equal(500, bError.Status)
	as.Equal("An error has occurred. Please reload the page and try again", bError.Error())
}

func (as *ActionSuite) Test_UserApproved() {
	defer mail.ClearMock()
	as.LoadFixture("handle requests")
	request := models.Request{}
	reset := models.Reset{}
	as.NoError(as.DB.Eager().Find(&request, "f19d5ab3-e913-4c08-b4dd-627e26ccbb0c"))
	req := as.JSON(fmt.Sprintf("/admin/%s", request.ID))
	resp := req.Put(&request)

	as.NoError(as.DB.Where("email = ? ", request.Email).First(&reset))
	m := mail.LastMessageFromMock()
	as.Equal("domain@example.com", m.To[0])
	as.Equal("noreply@ipa.example.com", m.From)
	as.Equal("Account request was approved", m.Subject)
	as.Equal(302, resp.Code)
	as.Equal("/admin/", resp.Location())
	as.Error(as.DB.Eager().Find(&request, "f19d5ab3-e913-4c08-b4dd-627e26ccbb0c"))

	as.Equal(
		&freeipa.User{
			UID:       "ttest",
			Givenname: &request.Name,
			Sn:        request.LastName,
			Mail:      &[]string{request.Email}},
		ipaclient.UserAddedFromMock())
}

func (as *ActionSuite) Test_UserDisapproved() {
	defer mail.ClearMock()
	as.LoadFixture("handle requests")
	request := models.Request{}
	as.NoError(as.DB.Eager().Find(&request, "f19d5ab3-e913-4c08-b4dd-627e26ccbb0c"))
	req := as.JSON(fmt.Sprintf("/admin/%s", request.ID))
	resp := req.Delete()

	m := mail.LastMessageFromMock()
	as.Equal("domain@example.com", m.To[0])
	as.Equal("noreply@ipa.example.com", m.From)
	as.Equal("Account request was declined", m.Subject)
	as.Equal(302, resp.Code)
	as.Equal("/admin/", resp.Location())
	as.Error(as.DB.Eager().Find(&request, "f19d5ab3-e913-4c08-b4dd-627e26ccbb0c"))
}
