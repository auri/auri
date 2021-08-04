package actions

import (
	"auri/mail"
	"auri/models"
)

func (as *ActionSuite) Test_EmailValidationInvalidToken() {
	as.LoadFixture("request-unverified-validToken-not-in-ipa")
	request := models.Request{}
	as.NoError(as.DB.Eager().Find(&request, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	req := as.JSON("/emailvalidation/%s", "37e16abd-b69e-4b65-a3cd-a9c7b250f722").Get()
	as.Equal(404, req.Code, req.Body.String())
}

func (as *ActionSuite) Test_EmailValidationValidToken() {
	defer mail.ClearMock()
	as.LoadFixture("request-unverified-validToken-not-in-ipa")
	request := models.Request{}
	as.NoError(as.DB.Eager().Find(&request, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))

	req := as.JSON("/emailvalidation/%s", request.Token).Get()
	as.Equal(200, req.Code, req.Body.String())
	as.NoError(as.DB.Eager().Where("email_verification = true").First(&request))
}

func (as *ActionSuite) Test_EmailValidationExpiredToken() {
	as.LoadFixture("request-unverified-expiredToken-not-in-ipa")
	request := models.Request{}
	as.NoError(as.DB.Eager().Find(&request, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))

	req := as.JSON("/emailvalidation/%s", request.Token).Get()
	as.Equal(400, req.Code, req.Body.String())
}

func (as *ActionSuite) Test_EmailValidationAlreadyValidated() {
	as.LoadFixture("request-verified-expiredToken-not-in-ipa")
	request := models.Request{}
	as.NoError(as.DB.Eager().Find(&request, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	as.True(request.MailVerification)

	req := as.JSON("/emailvalidation/%s", request.Token).Get()
	as.True(request.MailVerification)
	as.Equal(400, req.Code, req.Body.String())
}
