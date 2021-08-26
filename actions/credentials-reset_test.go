package actions

import (
	"auri/helpers"
	"auri/ipaclient"
	"auri/mail"
	"auri/models"
	"time"

	"github.com/tehwalris/go-freeipa/freeipa"
)

func (as *ActionSuite) Test_CredentialResetRequestPage() {
	res := as.HTML("/credentialreset").Get()
	as.Equal(200, res.Code, res.Body.String())
}

func (as *ActionSuite) Test_CredentialResetSubmitRequestUnknownEmail() {
	defer mail.ClearMock()
	reset := models.Reset{Email: "user-not-present-in-IPA@example.com"}
	req := as.JSON("/credentialreset/")
	reg := req.Post(&reset)
	as.Equal(200, reg.Code, reg.Body.String())
	as.False(mail.MessageExistsInMock())
}

func (as *ActionSuite) Test_CredentialResetSubmitValidRequestDuplicateValid() {
	defer mail.ClearMock()
	as.LoadFixture("resetToken-valid-user-present-in-IPA")
	reset := models.Reset{}
	as.NoError(as.DB.Eager().Find(&reset, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	as.True(reset.ExpiryDate.Truncate(time.Second).After(helpers.TimeNow().Truncate(time.Second)))
	count, err := as.DB.Eager().Where("email = ? ", reset.Email).Count(&models.Reset{})
	as.NoError(err)
	as.Equal(1, count)

	req := as.JSON("/credentialreset/")
	reg := req.Post(&models.Reset{Email: "user-present-in-IPA@example.com"})
	as.Equal(200, reg.Code, reg.Body.String())
	as.True(mail.MessageExistsInMock())
	m := mail.LastMessageFromMock()
	as.Equal("user-present-in-IPA@example.com", m.To[0])
	as.Equal("Password reset for your account", m.Subject)

	as.NoError(as.DB.Eager().Find(&reset, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	as.True(reset.ExpiryDate.Truncate(time.Second).After(helpers.TimeNow().Truncate(time.Second)))
	count, err = as.DB.Eager().Where("email = ? ", reset.Email).Count(&models.Reset{})
	as.NoError(err)
	as.Equal(1, count)
}

func (as *ActionSuite) Test_CredentialResetSubmitValidRequestDuplicateExpired() {
	defer mail.ClearMock()
	as.LoadFixture("resetToken-expired-user-present-in-IPA")
	reset := models.Reset{}

	as.NoError(as.DB.Eager().Find(&reset, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	as.True(reset.ExpiryDate.Truncate(time.Second).Before(helpers.TimeNow().Truncate(time.Second)))
	count, err := as.DB.Eager().Where("email = ? ", reset.Email).Count(&models.Reset{})
	as.NoError(err)
	as.Equal(1, count)

	req := as.JSON("/credentialreset/")
	reg := req.Post(&models.Reset{Email: "user-present-in-IPA@example.com"})
	as.Equal(200, reg.Code, reg.Body.String())
	as.True(mail.MessageExistsInMock())
	m := mail.LastMessageFromMock()
	as.Equal("user-present-in-IPA@example.com", m.To[0])
	as.Equal("Password reset for your account", m.Subject)

	as.NoError(as.DB.Eager().Find(&reset, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	as.True(reset.ExpiryDate.Truncate(time.Second).After(helpers.TimeNow().Truncate(time.Second)))
	count, err = as.DB.Eager().Where("email = ? ", reset.Email).Count(&models.Reset{})
	as.NoError(err)
	as.Equal(1, count)
}

func (as *ActionSuite) Test_CredentialResetSubmitRequestProtectedGroup() {
	defer mail.ClearMock()
	req := as.JSON("/credentialreset/")
	reg := req.Post(&models.Reset{Email: "admin-reset-present-in-IPA@example.com"})
	as.Equal(200, reg.Code, reg.Body.String())
	as.False(mail.MessageExistsInMock())
}

func (as *ActionSuite) Test_CredentialResetSubmitRequestKnownEmail() {
	defer mail.ClearMock()
	req := as.JSON("/credentialreset/")
	reg := req.Post(&models.Reset{Email: "user-present-in-IPA@example.com"})
	as.True(mail.MessageExistsInMock())
	m := mail.LastMessageFromMock()
	as.Equal(200, reg.Code, reg.Body.String())
	as.Equal("user-present-in-IPA@example.com", m.To[0])
	as.Equal("Password reset for your account", m.Subject)
}

func (as *ActionSuite) Test_CredentialResetSubmitRequestWrongEmail() {
	defer mail.ClearMock()

	req := as.JSON("/credentialreset/").Post(&models.Reset{Email: ""})
	as.Equal(400, req.Code, req.Body.String())

	req = as.JSON("/credentialreset/").Post(&models.Reset{Email: "somewrongformat"})
	as.Equal(400, req.Code, req.Body.String())

	req = as.JSON("/credentialreset/").Post(&models.Reset{Email: "somewrongformat@blablums"})
	as.Equal(400, req.Code, req.Body.String())

	req = as.JSON("/credentialreset/").Post(&models.Reset{Email: "wrongdomain@somewrongdomain.com"})
	as.Equal(400, req.Code, req.Body.String())

	as.False(mail.MessageExistsInMock())
}

func (as *ActionSuite) Test_CredentialResetTokenValid() {
	as.LoadFixture("resetToken-valid")
	reset := models.Reset{}
	as.NoError(as.DB.Eager().Find(&reset, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	req := as.HTML("/credentialreset/process/%s", reset.Token).Get()
	as.Equal(200, req.Code, req.Body.String())
}

func (as *ActionSuite) Test_CredentialResetTokenWrong() {
	reset := models.Reset{}
	as.Error(as.DB.Eager().Find(&reset, "ABCDEFb3-e913-4c08-b4dd-627e26ccbc0b"))
	req := as.HTML("/credentialreset/process/%s", "ABCDEFb3-e913-4c08-b4dd-627e26ccbc0b").Get()
	as.Equal(404, req.Code, req.Body.String())
}

func (as *ActionSuite) Test_CredentialResetTokenNotValid() {
	as.LoadFixture("resetToken-expired")
	reset := models.Reset{}
	as.NoError(as.DB.Eager().Find(&reset, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	req := as.HTML("/credentialreset/process/%s", reset.Token).Get()
	as.Equal(400, req.Code, req.Body.String())
}

func (as *ActionSuite) Test_CredentialResetProcessingTokenWrong() {
	reset := models.Reset{}
	request := models.Request{PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIK0wmN/Cr3JXqmLW7u+g9pTh+wyqDHpSQEIQczXkVx9q test@test.de"}
	as.Error(as.DB.Eager().Find(&reset, "ABCDEFb3-e913-4c08-b4dd-627e26ccbc0b"))
	res := as.JSON("/credentialreset/process/ABCDEFb3-e913-4c08-b4dd-627e26ccbc0b/create").Post(&request)
	as.Equal(404, res.Code, res.Body.String())
}

func (as *ActionSuite) Test_CredentialResetProcessingTokenNotValid() {
	as.LoadFixture("resetToken-expired")
	reset := models.Reset{}
	request := models.Request{PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIK0wmN/Cr3JXqmLW7u+g9pTh+wyqDHpSQEIQczXkVx9q test@test.de"}
	as.NoError(as.DB.Eager().Find(&reset, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	res := as.JSON("/credentialreset/process/" + reset.Token + "/create").Post(&request)
	as.Equal(400, res.Code, res.Body.String())
	as.NoError(as.DB.Eager().Find(&reset, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
}

func (as *ActionSuite) Test_CredentialResetProcessingProtectedGroup() {
	as.LoadFixture("resetToken-valid-admin-user-present-in-IPA")
	reset := models.Reset{}
	request := models.Request{PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIK0wmN/Cr3JXqmLW7u+g9pTh+wyqDHpSQEIQczXkVx9q test@test.de"}
	as.NoError(as.DB.Eager().Find(&reset, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	res := as.JSON("/credentialreset/process/" + reset.Token + "/create").Post(&request)
	as.Equal(500, res.Code, res.Body.String())
	as.NoError(as.DB.Eager().Find(&reset, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
}

func (as *ActionSuite) Test_CredentialResetProcessingAddSSHKey() {
	as.LoadFixture("resetToken-valid-user-present-in-IPA")
	reset := models.Reset{}
	request := models.Request{PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIK0wmN/Cr3JXqmLW7u+g9pTh+wyqDHpSQEIQczXkVx9q test@test.de"}
	as.NoError(as.DB.Eager().Find(&reset, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	res := as.JSON("/credentialreset/process/" + reset.Token + "/create").Post(&request)
	as.Equal(200, res.Code, res.Body.String())
	as.Error(as.DB.Eager().Find(&reset, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
}

func (as *ActionSuite) Test_CredentialResetProcessingAddPuttySSHKey() {
	as.LoadFixture("resetToken-valid-user-present-in-IPA")
	reset := models.Reset{}
	request := models.Request{PublicKey: `---- BEGIN SSH2 PUBLIC KEY ----
Comment: "ed25519-key-20210826"
AAAAC3NzaC1lZDI1NTE5AAAAIALoF39Bi+IqrjGnRdXSZRA8ih/FcB3NXWamTLLu
o4uJ
---- END SSH2 PUBLIC KEY ----`}
	as.NoError(as.DB.Eager().Find(&reset, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	res := as.JSON("/credentialreset/process/" + reset.Token + "/create").Post(&request)
	as.Equal(200, res.Code, res.Body.String())
	as.Error(as.DB.Eager().Find(&reset, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))

	as.Equal(
		&freeipa.User{
			UID:          "user-present-in-IPA",
			Ipasshpubkey: &[]string{"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIALoF39Bi+IqrjGnRdXSZRA8ih/FcB3NXWamTLLuo4uJ ed25519-key-20210826"},
		},
		ipaclient.UserModifiedFromMock())
}

func (as *ActionSuite) Test_CredentialResetProcessingAddPassword() {
	as.LoadFixture("resetToken-valid-user-present-in-IPA")
	reset := models.Reset{}
	request := models.Request{Password: "Testtest12345678", PasswordConfirmation: "Testtest12345678"}
	as.NoError(as.DB.Eager().Find(&reset, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	res := as.JSON("/credentialreset/process/" + reset.Token + "/create").Post(&request)
	as.Equal(200, res.Code, res.Body.String())
	as.Error(as.DB.Eager().Find(&reset, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
}

func (as *ActionSuite) Test_CredentialResetProcessingAddPasswordAndSSHKey() {
	as.LoadFixture("resetToken-valid-user-present-in-IPA")
	reset := models.Reset{}
	request := models.Request{PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIK0wmN/Cr3JXqmLW7u+g9pTh+wyqDHpSQEIQczXkVx9q test@test.de", Password: "Testtest12345678", PasswordConfirmation: "Testtest12345678"}
	as.NoError(as.DB.Eager().Find(&reset, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	res := as.JSON("/credentialreset/process/" + reset.Token + "/create").Post(&request)
	as.Equal(200, res.Code, res.Body.String())
	as.Error(as.DB.Eager().Find(&reset, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
}
