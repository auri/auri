package notifications

import (
	"auri/config"
	"auri/logger"
	"auri/mail"
)

//AdminNotificationNewRequest notifies admins about new request
func AdminNotificationNewRequest(email, comment string) error {
	if config.GetInstance().AdminEmailNotificationsEnable {
		if err := mail.SendAdminNotification(email, comment); err != nil {
			return err
		}
		logger.AuriLogger.Info("Admin notification: admins were successfully notified about new request via email")
	} else {
		logger.AuriLogger.Info("Admin notification: email notifications are disabled, skipping notification")
	}

	return nil
}
