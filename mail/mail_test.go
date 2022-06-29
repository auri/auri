package mail

import (
	"auri/models"
	"testing"

	"github.com/gobuffalo/buffalo/render"

	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/suite/v3"
)

type MailSuite struct {
	*suite.Model
}

func Test_MailSuite(t *testing.T) {
	model := suite.NewModel()

	ms := &MailSuite{
		Model: model,
	}
	suite.Run(t, ms)
}

func (ms *MailSuite) Test_SendMailToUserOK() {
	New = NewMock
	defer ClearMock()
	request := &models.Request{Email: "test@example.com", PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAEjkPQpaaFHDTIdUJVx8jA4esihdfioaalmWXUlY0QK test@test.DE"}
	err := SendTokenMail(TokenMailOpts{Email: request.Email, Token: request.Token}, EMailValidation)
	ms.NoError(err)
	m := LastMessageFromMock()
	ms.Equal("test@example.com", m.To[0])
	ms.Equal("noreply@ipa.example.com", m.From)
	ms.Equal("Validation of email address required", m.Subject)
	body := []mail.Body{{Content: "<a  href=https://<%= BASE_URL %>/credentialreset/process/<%= userToken %>/new>Validate</a>"}}
	ms.NotContains(m.Bodies, body)

}
func (ms *MailSuite) Test_SendMailToUserForPwrOK() {
	New = NewMock
	defer ClearMock()
	request := &models.Request{Email: "test@example.com"}
	err := SendTokenMail(TokenMailOpts{Email: request.Email, Token: request.Token}, EMailValidation)
	ms.NoError(err)
	m := LastMessageFromMock()
	ms.Equal("test@example.com", m.To[0])
	ms.Equal("noreply@ipa.example.com", m.From)
	ms.Equal("Validation of email address required", m.Subject)
	body := []mail.Body{{Content: "<a  href=https://<%= BASE_URL %>/emailvalidation/<%= userToken %>Validate</a>"}}
	ms.NotContains(m.Bodies, body)

}

func (ms *MailSuite) Test_SendNewReqestMailToAdminOK() {
	New = NewMock
	defer ClearMock()
	request := &models.Request{}
	err := SendAdminNewRequestNotification(request.Email, request.CommentField.String)
	ms.NoError(err)
	m := LastMessageFromMock()
	ms.Equal("admin1@example.com", m.To[0])
	ms.Equal("admin2@example.com", m.To[1])
	ms.Equal("noreply@ipa.example.com", m.From)
	ms.Equal("New account request", m.Subject)

}

func (ms *MailSuite) Test_SendMail() {
	tests := []struct {
		name     string
		subject  string
		to       []string
		template string
		data     render.Data
		wantErr  string // empty for no erro
	}{
		{
			"valid approval mail",
			"Valid approval mail",
			[]string{"test@example.com"},
			"request_approved.plush.html",
			map[string]interface{}{
				"Email": "test@example.com",
				"Token": "TESTTOKEN",
				"Login": "testlogin",
			},
			"",
		},
		{
			"wrong template",
			"wrong template",
			[]string{"test@example.com"},
			"wrong-template",
			map[string]interface{}{
				"Email": "test@example.com",
				"Token": "TESTTOKEN",
			},
			"Can't load the email template: wrong-template: could not find template wrong-template.html",
		},
	}

	New = NewMock

	for _, tt := range tests {
		err := Send(tt.subject, tt.to, tt.template, tt.data)

		if (err != nil) != (tt.wantErr != "") {
			ms.Error(err)
		}

		if err != nil {
			ms.EqualError(err, tt.wantErr)
		} else {
			ms.NoError(err)
			m := LastMessageFromMock()
			ms.Equal(tt.to, m.To)
			ms.Equal("noreply@ipa.example.com", m.From)
			ms.Equal(tt.subject, m.Subject)
		}

		ClearMock()
	}
}
