// +build !prodBuild

package actions

import (
	"auri/config"
	"auri/helpers"
	"auri/logger"
	"auri/models"

	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/nulls"
	"github.com/pkg/errors"
)

//ValidateTokenIPA this will parse token and check for expiry date
func ValidateTokenIPA(c buffalo.Context) error {
	tx, err := getDBConnection(c)
	if err != nil {
		return err
	}

	user := models.IpaUser{}
	token := c.Param("ipa_token_account_validation")

	if err := tx.Where("token = ? ", token).First(&user); err != nil {
		logger.AuriLogger.Errorf("Account revalidation: confirmation token %v not found", token)
		return c.Render(http.StatusBadRequest, tokenExpiredPageRenderer())
	}

	if user.NotifiedAt.Time.After(helpers.TimeNowAddHours(config.GetInstance().AccountValidationTimeframe)) {
		logger.AuriLogger.Errorf("Account revalidation: confirmation reset token for user %v expired", user.UID)
		return c.Render(http.StatusBadRequest, tokenExpiredPageRenderer())
	}

	user.NotifiedAt = nulls.Time{}
	user.NotificationsSent = 0
	user.RemindedAt = helpers.TimeNow()
	user.Token = ""

	if err := tx.Update(&user); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(http.StatusOK, r.HTML("ipa-validation-success.html", config.GetInstance().Layout))
}
