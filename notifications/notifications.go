package notifications

import (
	"auri/mail"
)

//AdminNotificationNewRequest notifies admins about new request
func AdminNotificationNewRequest(email, comment string) error {
	return runAdminEmailTransaction("new request from "+email, func() error {
		return mail.SendAdminNewRequestNotification(email, comment)
	})
}
