package grifts

import (
	"auri/logger"
	"auri/models"

	"github.com/gobuffalo/pop/v6"
	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
)

var _ = grift.Desc("cleanup_requests", "Remove expired and not validated account requests")

var _ = grift.Add("cleanup_requests", func(c *grift.Context) error {

	var err error
	requests := models.Requests{}

	if err := models.DB.Transaction(func(tx *pop.Connection) error {

		if err := tx.Where("expiry_date < now() AND email_verification = FALSE").All(&requests); err != nil {
			return errors.WithMessage(err, "could not find requests")
		}

		for _, request := range requests {
			logger.AuriLogger.Infof("Cleanup: account request from email %v is expired and will be removed", request.Email)
		}

		return tx.Destroy(&requests)
	}); err != nil {
		return errors.WithMessage(err, "Can not delete user")
	}

	return err
})
