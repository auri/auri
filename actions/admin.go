package actions

import (
	"auri/actions/apperror"
	"auri/config"
	"auri/helpers"
	"auri/ipaclient"
	"auri/logger"
	"auri/mail"
	"auri/models"
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
)

// RequestResource type struct
type RequestResource struct {
	buffalo.BaseResource
}

//List implements the listing of account requests
func (v RequestResource) List(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		c.Logger().Error("no connection to the database present")
		err := errors.New("An error has occurred. Please reload the page and try again")
		c.Flash().Add("danger", err.Error())
		return c.Error(500, err)
	}
	requests := models.Requests{}
	err := tx.Order("created_at desc").All(&requests)
	if err != nil {
		c.Flash().Add("danger", err.Error())
		return c.Error(500, err)
	}
	c.Set("requests", requests)
	return c.Render(http.StatusOK, r.HTML("admin.html", config.GetInstance().Layout))
}

// Update implements the approval action in the admin interface
func (v RequestResource) Update(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		c.Logger().Error("no connection to the database present")
		err := errors.New("An error has occurred. Please reload the page and try again")
		c.Flash().Add("danger", err.Error())
		return c.Error(500, err)
	}

	request := models.Request{}

	if err := tx.Find(&request, c.Param("request_id")); err != nil {
		return c.Error(400, err)
	}

	// try to find a new unused login, max 5 tries
	i := 0
	var login string
	for {
		if i >= 5 {
			logger.AuriLogger.Errorf("Admin area: can't find unused login for user %v %v after %v tries", request.Name, request.LastName, i)
			c.Flash().Add("danger", fmt.Sprintf("Can't find unused login for user %v %v after %v tries", request.Name, request.LastName, i))
			return c.Redirect(http.StatusFound, "adminPath()")
		}

		login = request.Login(i)

		loginExist, err := ipaclient.LoginExists(login)
		if err != nil {
			return err
		}

		if !loginExist {
			break
		}
		i++
	}

	if err := ipaclient.AddUser(login, request.Name, request.LastName, request.Email, request.PublicKey); err != nil {
		apperror.InvokeError(c, apperror.IPAError, err)
		return c.Redirect(http.StatusFound, "adminPath()")
	}

	logger.AuriLogger.Infof("Admin area: account %v with email %v was created in IPA (without configured authentication data)", login, request.Email)
	c.Flash().Add("success", "Account "+login+" with email "+request.Email+" was created in IPA (without configured authentication data)")

	// after we created account in IPA, we should destroy the request regardless of future errors: otherwise we would get multiple accounts if something in the later steps fails
	if err := tx.Destroy(&request); err != nil {
		return c.Error(400, err)
	}

	//After the user been added to ipa, we send reset notification to user to set his password
	//or his ssh key
	reset := models.Reset{}

	if err := reset.NewToken(); err != nil {
		return err
	}

	reset.ExpiryDate = helpers.TimeNowAddHours(config.GetInstance().InitialCredentialSetupTimeframe)
	reset.Email = request.Email
	reset.Login = login
	//We save the email in reset table
	if err := tx.Create(&reset); err != nil {
		apperror.InvokeError(c, apperror.DBError, err)
		return c.Redirect(http.StatusFound, "adminPath()")
	}
	logger.AuriLogger.Infof("Admin area: credentials reset token was created for account with email %v", request.Email)

	// Send approval notification to the user
	if err := mail.SendTokenMail(mail.TokenMailOpts{Email: reset.Email, Token: reset.Token, Login: login}, mail.ApprovalMail); err != nil {
		apperror.InvokeError(c, apperror.MailError, err)
		return c.Redirect(http.StatusFound, "adminPath()")
	}
	logger.AuriLogger.Infof("Admin area: user %v(%v) was informed about approval of account request and requested to set authentication data", login, request.Email)
	c.Flash().Add("success", "User "+login+"("+request.Email+") has been informed and requested to set authentication data")

	return c.Redirect(http.StatusFound, "adminPath()")
}

// Destroy implements the account rejection
func (v RequestResource) Destroy(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		c.Logger().Error(500, "no connection to the database present")
		err := errors.New("An error has occurred. Please reload the page and try again")
		c.Flash().Add("danger", err.Error())
		return c.Error(500, err)
	}

	request := models.Request{}

	if err := tx.Find(&request, c.Param("request_id")); err != nil {
		return c.Error(400, err)
	}

	if err := mail.SendDecline(request.Email); err != nil {
		apperror.InvokeError(c, apperror.MailError, err)
		return c.Redirect(http.StatusFound, "adminPath()")
	}
	logger.AuriLogger.Infof("Admin area: user with email %v was informed about rejection of account request", request.Email)

	if err := tx.Destroy(&request); err != nil {
		return c.Error(400, err)
	}

	logger.AuriLogger.Infof("Admin area: account request with email %v was removed", request.Email)
	c.Flash().Add("danger", "Account request from "+request.Email+" was rejected, requester was notified")

	return c.Redirect(http.StatusFound, "adminPath()")
}
