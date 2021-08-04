MockSMTP
========

MockSMTP is an in-memory implementation for buffalo `Sender` interface. It's intended to catch sent messages for test purposes.

Usage
-----

mailers/mailers_test.go:
```go
package mailers_test

import (
	"github.com/stanislas-m/mocksmtp"
	"gitlab.com/stanislas-m/mybuffaloproject/mailers"
)

func init() {
	mailers.SMTP = mocksmtp.New()
}

```

mailers/registration_mail.go:
```go
package mailers

import (
	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/buffalo/render"
	"github.com/pkg/errors"
)

func SendRegistrationMail(email string) error {
	m := mail.NewMessage()

	m.Subject = "Registration succeed"
	m.From = "Buffalo <noreply@gobuffalo.io>"
	m.To = []string{email}
	err := m.AddBody(r.HTML("registration_mail.html"), render.Data{})
	if err != nil {
		return errors.WithStack(err)
	}
	if err := SMTP.Send(m); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

```

mailers/registration_mail_test.go:
```go
package mailers_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/stanislas-m/mybuffaloproject/mailers"
)

func Test_SendRegistrationMail(t *testing.T) {
    r := require.New(t)
    ms, v := mailers.SMTP.(*mocksmtp.MockSMTP)
    r.True(v)
    // Clear SMTP queue after test
    defer ms.Clear()
    err := mailers.SendRegistrationMail("test@gobuffalo.io")
    r.NoError(err)
    // Test email
    r.Equal(1, ms.Count())
    m, err := ms.LastMessage()
    r.NoError(err)
    r.Equal("Buffalo <noreply@gobuffalo.io>", m.From)
    r.Equal("Registration succeed", m.Subject)
    r.Equal(1, len(m.To))
    r.Equal("test@gobuffalo.io", m.To[0])
}
```
