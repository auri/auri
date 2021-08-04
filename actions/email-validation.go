package actions

import (
	"auri/actions/apperror"
	"auri/config"
	"auri/helpers"
	"auri/logger"
	"auri/mail"
	"auri/models"
	"net/http"

	"github.com/gobuffalo/buffalo"
)

//EmailValidationProcessRequest validates if token expired or not expired
func EmailValidationProcessRequest(c buffalo.Context) error {
	logger.AuriLogger.Debugf("Email validation: email validation was invoked via token")

	tx, err := getDBConnection(c)
	if err != nil {
		return err
	}

	request := models.Request{}
	token := c.Param("token")
	request.Token = token
	if err := tx.Where("token = ? ", request.Token).First(&request); err != nil {
		logger.AuriLogger.Errorf("Email validation: token %v could not be found", token)
		return c.Render(http.StatusNotFound, tokenExpiredPageRenderer())
	}
	if request.MailVerification == true {
		logger.AuriLogger.Errorf("Email validation: validation for account request from email %v was already done", request.Email)
		return c.Render(http.StatusBadRequest, tokenExpiredPageRenderer())
	}
	if helpers.TimeNow().After(request.ExpiryDate) == true {
		logger.AuriLogger.Errorf("Email validation: validation time for account request from email %v is expired", request.Email)
		return c.Render(http.StatusBadRequest, tokenExpiredPageRenderer())
	}
	request.MailVerification = true
	tx.ValidateAndSave(&request)
	logger.AuriLogger.Infof("Email validation: account request from email %v was succesfully valdiated", request.Email)

	if err := mail.SendAdminNotification(request.Email, request.CommentField.String); err != nil {
		apperror.InvokeError(c, apperror.MailError, err)
		return c.Error(http.StatusBadRequest, err)
	}
	logger.AuriLogger.Info("Email validation: admins were successfully notified about new request")

	return c.Render(http.StatusOK, r.HTML("account-request-mail-verified.plush.html", config.GetInstance().Layout))
}
