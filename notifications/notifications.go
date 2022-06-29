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

//AdminNotificationRequestApproved notifies admins about approval of account request
func AdminNotificationRequestApproved(admin, email, comment string) error {
	return runAdminEmailTransaction("approval of request from "+email, func() error {
		return mail.SendAdminRequestApprovedNotification(admin, email, comment)
	})
}

//AdminNotificationRequestDeclined notifies admins about approval of account request
func AdminNotificationRequestDeclined(admin, email, comment string) error {
	return runAdminEmailTransaction("rejection of request from "+email, func() error {
		return mail.SendAdminRequestDeclinedNotification(admin, email, comment)
	})
}
