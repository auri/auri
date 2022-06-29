package notifications

import (
	"auri/config"
	"auri/logger"
)

//runAdminEmailTransaction runs the given email admin notification transaction
//if admin email notifications are enabled
func runAdminEmailTransaction(name string, f func() error) error {
	if !config.GetInstance().AdminEmailNotificationsEnable {
		logger.AuriLogger.Infof("Admin email notification: email notifications are disabled, skipping notification for %v", name)
		return nil
	}

	if err := f(); err != nil {
		return err
	}

	logger.AuriLogger.Infof("Admin email notification: admins were successfully notified via email about %v", name)
	return nil
}
