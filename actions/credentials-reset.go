package actions

import (
	"auri/actions/apperror"
	"auri/config"
	"auri/helpers"
	"auri/helpers/ssh"
	"auri/ipaclient"
	"auri/logger"
	"auri/mail"
	"auri/models"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/pkg/errors"
)

//CredentialsResetRequest provides the main page for requesting credentials reset
func CredentialsResetRequest(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("credentials-reset-request.plush.html", config.GetInstance().Layout))
}

//CredentialsResetProcessRequest verifies the credentials reset request and submits the reset token via email
func CredentialsResetProcessRequest(c buffalo.Context) error {
	// error handler for internal errors happening in this processing handler
	internalErrorHandler := func(errorType apperror.ErrorType, err error) error {
		apperror.InvokeError(c, errorType, err)
		return c.Render(http.StatusInternalServerError, r.HTML("credentials-reset-request.plush.html", config.GetInstance().Layout))
	}

	reset := &models.Reset{}

	if err := c.Bind(reset); err != nil {
		return internalErrorHandler(apperror.AuriError, err)
	}

	tx, err := getDBConnection(c)
	if err != nil {
		return internalErrorHandler(apperror.DBError, err)
	}

	validationError, err := reset.ValidateEmail(tx)
	if err != nil {
		return internalErrorHandler(apperror.AuriError, err)
	}

	if validationError.HasAny() == true {
		c.Set("errors", validationError.Errors)
		return c.Render(http.StatusBadRequest, r.HTML("credentials-reset-request.plush.html", config.GetInstance().Layout))
	}

	// renderer of resulting confirmation page
	// we do not want to expose any error situations here, as potential attacker could gain information via that
	finalPageRenderer := func() error {
		c.Set("userEmail", reset.Email)
		return c.Render(http.StatusOK, r.HTML("credentials-reset-request-confirmation.plush.html", config.GetInstance().Layout))
	}

	logger.AuriLogger.Infof("Credential reset: new request for account with email %v", reset.Email)

	login, err := ipaclient.UserFind(reset.Email)
	if err != nil {
		logger.AuriLogger.Error(errors.WithMessage(err, "Credential reset: could not find user for provided email "+reset.Email))
		return finalPageRenderer()
	}

	// login wasn't found
	if login == "" {
		logger.AuriLogger.Infof("Credential reset: could not find user account with email %v, ignoring reset request", reset.Email)
		return finalPageRenderer()
	}

	logger.AuriLogger.Debugf("Credential reset: IPA account %v was found for email %v", login, reset.Email)

	// check if we are allowed to process password resets for given groups
	var protectedGroupFound string
	if protectedGroupFound, err = credentialResetVerifyProtectedGroup(login); err != nil {
		return internalErrorHandler(apperror.IPAError, err)
	}

	if protectedGroupFound != "" {
		logger.AuriLogger.Infof("Credential reset: user %v is member of protected group %v, ignoring reset request", login, protectedGroupFound)
		return finalPageRenderer()
	}

	// processing reset request
	// check if we already have a reset request for this login
	requestPresent := false
	if err := tx.Where("login = ? ", login).First(reset); err == nil {
		requestPresent = true
	}
	reset.ExpiryDate = helpers.TimeNowAddHours(config.GetInstance().PasswordResetTimeframe)

	if !requestPresent { // create a new one request or update the existing one
		logger.AuriLogger.Debugf("Credential reset: creating new credentials reset request")
		reset.Login = login
		if err := reset.NewToken(); err != nil {
			return internalErrorHandler(apperror.AuriError, err)
		}
		if err := tx.Create(reset); err != nil {
			return internalErrorHandler(apperror.DBError, err)
		}
	} else {
		logger.AuriLogger.Debugf("Credential reset: credentials reset already present, will update it")
		if err := tx.Update(reset); err != nil {
			return internalErrorHandler(apperror.DBError, err)
		}
	}

	logger.AuriLogger.Debugf("Credential reset: credentials reset token was created/updated for account %v(%v)", reset.Login, reset.Email)

	if err := mail.SendTokenMail(mail.TokenMailOpts{Email: reset.Email, Token: reset.Token}, mail.PasswordReset); err != nil {
		return internalErrorHandler(apperror.MailError, err)
	}

	logger.AuriLogger.Infof("Credential reset: user %v(%v) was informed about credential reset and requested to set authentication data", reset.Login, reset.Email)
	return finalPageRenderer()
}

//CredentialsResetProcessToken processes the given reset token and provides a reset form if token is valid
func CredentialsResetProcessToken(c buffalo.Context) error {
	tx, err := getDBConnection(c)
	if err != nil {
		return err
	}

	logger.AuriLogger.Debugf("Credential reset: credential reset invoked via token")

	reset, err := credentialResetGetValidToken(tx, c, c.Param("tokenpwr"))
	if reset == nil || err != nil {
		return err
	}

	c.Set("tokenpwr", reset.Token)

	logger.AuriLogger.Infof("Credential reset: valid credential reset token for user %v(%v) found, requesting data for reset of credentials", reset.Login, reset.Email)

	return c.Render(http.StatusOK, r.HTML("credentials-reset-data-collection.plush.html", config.GetInstance().Layout))
}

//CredentialsResetProcessCredentials processes the credential reset
func CredentialsResetProcessCredentials(c buffalo.Context) error {
	// error handler for internal errors happening in this processing handler
	internalErrorHandler := func(errorType apperror.ErrorType, err error) error {
		apperror.InvokeError(c, errorType, err)
		return c.Render(http.StatusInternalServerError, r.HTML("credentials-reset-data-collection.plush.html", config.GetInstance().Layout))
	}

	logger.AuriLogger.Debugf("Credential reset: processing credential reset with provided data")

	request := models.Request{}

	tx, err := getDBConnection(c)
	if err != nil {
		return err
	}

	reset, err := credentialResetGetValidToken(tx, c, c.Param("tokenpwr"))
	if reset == nil || err != nil {
		return err
	}

	c.Set("tokenpwr", reset.Token)

	// check if we are allowed to process password resets for given groups
	var protectedGroupFound string
	if protectedGroupFound, err = credentialResetVerifyProtectedGroup(reset.Login); err != nil {
		return internalErrorHandler(apperror.IPAError, err)
	}
	if protectedGroupFound != "" {
		return internalErrorHandler(apperror.AuriError, errors.Errorf("Credential reset: user %v is member of protected group %v, ignoring existing reset token", reset.Login, protectedGroupFound))
	}

	if err := c.Bind(&request); err != nil {
		return errors.WithStack(err)
	}
	if len(request.PublicKey) != 0 {
		//try to convert the key from putty format
		if newKey, err := ssh.ConvertPuttySSH(request.PublicKey); err == nil {
			logger.AuriLogger.Infof("Credential reset: got public key in Putty format, converted it to OpenSSH authorized_keys format")
			request.PublicKey = newKey
		}
	}
	validate := request.ValidateSSHOrPassword()

	if validate.HasAny() == true {
		c.Set("errors", validate.Errors)
		return c.Render(http.StatusBadRequest, r.HTML("credentials-reset-data-collection.plush.html", config.GetInstance().Layout))
	}

	if err := ipaclient.UserModify(reset.Login, request.Password, request.PublicKey); err != nil {
		return internalErrorHandler(apperror.IPAError, err)
	}
	logger.AuriLogger.Infof("Credential reset: new authentication credentials for user %v(%v) were set in IPA", reset.Login, reset.Email)

	if err := tx.Destroy(reset); err != nil {
		return internalErrorHandler(apperror.DBError, err)
	}
	logger.AuriLogger.Debugf("Credential reset: reset token for user %v(%v) was removed", reset.Login, reset.Email)

	return c.Render(http.StatusOK, r.HTML("credentials-reset-done.plush.html", config.GetInstance().Layout))
}

// credentialResetGetValidToken returns a valid Reset token
// if token is not found or expired, returns nil
func credentialResetGetValidToken(tx *pop.Connection, c buffalo.Context, token string) (*models.Reset, error) {
	reset := &models.Reset{}
	if err := tx.Where("token = ? ", token).First(reset); err != nil {
		logger.AuriLogger.Errorf("Credential reset: credential reset token %v not found", token)
		return nil, c.Render(http.StatusNotFound, tokenExpiredPageRenderer())
	}
	if helpers.TimeNow().After(reset.ExpiryDate) == true {
		logger.AuriLogger.Errorf("Credential reset: credential reset token for user %v(%v) expired", reset.Login, reset.Email)
		return nil, c.Render(http.StatusBadRequest, tokenExpiredPageRenderer())
	}
	return reset, nil
}

// credentialResetVerifyProtectedGroup checks if given user is a member of some protected group
// and isn't allowed to reset credentials
// returns a first group matched if any, empty string otherwise
func credentialResetVerifyProtectedGroup(login string) (string, error) {
	groups, err := ipaclient.GetGroups(login)
	if err != nil {
		return "", err
	}
	logger.AuriLogger.Debugf("Credential reset: user is member of following groups: %v", groups)
	var protectedGroupFound string
protectedLoop:
	for _, confGroup := range config.GetInstance().ProtectedGroups {
		for _, userGroup := range groups {
			if confGroup == userGroup {
				protectedGroupFound = userGroup
				break protectedLoop
			}
		}
	}

	return protectedGroupFound, nil
}
