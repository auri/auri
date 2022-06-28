package mail

import (
	"auri/config"
	"auri/logger"
	"auri/public"
	template "auri/templates"

	"crypto/tls"
	"strconv"

	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/buffalo/render"
	"github.com/pkg/errors"
)

var sender mail.Sender
var renderer *render.Engine

// New links to the constructor, which is used to create the sender
var New = NewSender

//NewSender returns initialized mail.Sender
func NewSender(host, port, user, password string, enforceStartTLS, TLSSkipVerify bool) (mail.Sender, error) {
	sender, err := mail.NewSMTPSender(
		host,
		port,
		user,
		password,
	)

	if err != nil {
		return nil, err
	}

	if enforceStartTLS {
		// require encryption via MandatoryStartTLS
		// https://github.com/gobuffalo/buffalo/blob/6bbc80e6c4fe8606145e4a887f0ec59a5cc41ef9/mail/internal/mail/smtp.go#L156-L159
		sender.Dialer.StartTLSPolicy = 1
	}

	if TLSSkipVerify {
		sender.Dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	} else {
		sender.Dialer.TLSConfig = &tls.Config{ServerName: host}
	}

	return sender, nil
}

// GetSender returns the current initialized instance of mail.Sender
func GetSender() (mail.Sender, error) {
	if sender == nil {
		conf := config.GetInstance()

		lsender, err := New(
			conf.SMTPHost,
			strconv.FormatUint(uint64(conf.SMTPPort), 10),
			conf.SMTPUser,
			conf.SMTPPassword,
			conf.SMTPEnforceStartTLS,
			conf.SMTPTLSSkipVerify,
		)

		if err != nil {
			return nil, err
		}

		sender = lsender
	}
	return sender, nil
}

// getRenderer returns the current initialized instance of template renderer
func getRenderer() *render.Engine {
	if renderer == nil {
		renderer = render.New(render.Options{
			HTMLLayout: "mail/layout.html",
			// fs.FS containing templates
			TemplatesFS: template.FS(),

			// fs.FS containing assets
			AssetsFS: public.FS(),
			Helpers:  render.Helpers{},
		})
	}
	return renderer
}

//Send sends an email with given data
func Send(subject string, to []string, template string, data render.Data) error {
	m := mail.NewMessage()
	m.From = config.GetInstance().FromEmail
	m.Subject = subject
	m.To = to

	// some generic settings for our templates
	conf := config.GetInstance()
	data["BaseURL"] = conf.BaseURL
	data["ContactEmail"] = conf.ContactEmail

	if err := m.AddBodies(data, getRenderer().HTML("mail/"+template)); err != nil {
		return errors.WithMessage(err, "Can't load the email template")
	}

	if !config.GetInstance().SMTPEnable {
		logger.BuffaloLogger.Warnf("SMTP functionality is disabled. Skipping email with subject '%v'", subject)
		return nil
	}

	sender, err := GetSender()
	if err != nil {
		return errors.WithMessage(err, "Can't send the email")
	}

	if err := sender.Send(m); err != nil {
		return errors.WithMessage(err, "Can't send the email")
	}

	return nil
}
