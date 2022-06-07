package grifts

import (
	"auri/logger"
	"auri/models"

	"github.com/gobuffalo/pop/v6"
	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
)

var _ = grift.Desc("cleanup_reset_tokens", "Remove expired tokens for credential reset")

var _ = grift.Add("cleanup_reset_tokens", func(c *grift.Context) error {

	var err error
	resets := []models.Reset{}

	if err := models.DB.Transaction(func(tx *pop.Connection) error {

		if err := tx.Where("expiry_date < now()").All(&resets); err != nil {
			return errors.WithMessage(err, "Can not find tokens")
		}

		for _, reset := range resets {
			logger.AuriLogger.Infof("Cleanup: credential reset token from email %v is expired and will be removed", reset.Email)
		}

		return tx.Destroy(&resets)
	}); err != nil {
		return errors.WithMessage(err, "Can not delete reset token")
	}

	return err
})
