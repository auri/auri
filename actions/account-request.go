package actions

import (
	"auri/actions/apperror"
	"auri/config"
	"auri/logger"
	"auri/mail"
	"auri/models"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/x/defaults"
	"github.com/pkg/errors"
)

//AccountRequestPage serves the account request page
func AccountRequestPage(c buffalo.Context) error {
	c.Set("mail", defaults.String(c.Request().FormValue("mail"), ""))
	c.Set("name", defaults.String(c.Request().FormValue("name"), ""))
	c.Set("lastname", defaults.String(c.Request().FormValue("lastname"), ""))
	c.Set("comment", defaults.String(c.Request().FormValue("comment"), ""))
	return c.Render(http.StatusOK, accountRequestPageRenderer())
}

//AccountRequestProcessRequest to validate user inputs
func AccountRequestProcessRequest(c buffalo.Context) error {
	c.Set("mail", defaults.String(c.Request().FormValue("mail"), ""))
	c.Set("name", defaults.String(c.Request().FormValue("name"), ""))
	c.Set("lastname", defaults.String(c.Request().FormValue("lastname"), ""))
	c.Set("comment", defaults.String(c.Request().FormValue("comment"), ""))

	request := &models.Request{}
	if err := c.Bind(request); err != nil {
		return errors.WithStack(err)
	}

	tx, err := getDBConnection(c)
	if err != nil {
		return err
	}

	validationError, err := request.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}
	if validationError.HasAny() == true {
		c.Set("errors", validationError.Errors)
		return c.Render(http.StatusBadRequest, accountRequestPageRenderer())
	}

	logger.AuriLogger.Infof("Account request: new account request for email %v", request.Email)

	if err := mail.SendTokenMail(mail.TokenMailOpts{Email: request.Email, Token: request.Token}, mail.EMailValidation); err != nil {
		apperror.InvokeError(c, apperror.MailError, err)
		return c.Render(http.StatusBadRequest, accountRequestPageRenderer())
	}
	logger.AuriLogger.Infof("Account request: requester with email %v was notified about required validation", request.Email)

	return c.Render(http.StatusCreated, r.HTML("account-request-created.plush.html", config.GetInstance().Layout))
}

//accountRequestPageRenderer returns the renderer of account request page
func accountRequestPageRenderer() render.Renderer {
	return r.HTML("account-request.plush.html", config.GetInstance().Layout)
}
