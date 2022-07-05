package notifications

import (
	"auri/config"
	"auri/mail"
)

//AdminNotificationNewRequest notifies admins about new request
func AdminNotificationNewRequest(email, comment string) error {
	if err := runAdminEmailTransaction("new request from "+email, func() error {
		return mail.SendAdminNewRequestNotification(email, comment)
	}); err != nil {
		return err
	}

	return runAdminShellTransaction(
		"new request from "+email,
		config.GetInstance().AdminShellNewRequest,
		"NEW_REQUEST",
		email,
		comment,
		"",
	)
}

//AdminNotificationRequestApproved notifies admins about approval of account request
func AdminNotificationRequestApproved(admin, email, comment string) error {
	if err := runAdminEmailTransaction("approval of request from "+email, func() error {
		return mail.SendAdminRequestApprovedNotification(admin, email, comment)
	}); err != nil {
		return err
	}

	return runAdminShellTransaction(
		"approval of request from "+email,
		config.GetInstance().AdminShellRequestApproved,
		"REQUEST_APPROVED",
		email,
		comment,
		admin,
	)
}

//AdminNotificationRequestDeclined notifies admins about rejection of account request
func AdminNotificationRequestDeclined(admin, email, comment string) error {
	if err := runAdminEmailTransaction("rejection of request from "+email, func() error {
		return mail.SendAdminRequestDeclinedNotification(admin, email, comment)
	}); err != nil {
		return err
	}

	return runAdminShellTransaction(
		"rejection of request from "+email,
		config.GetInstance().AdminShellRequestDeclined,
		"REQUEST_DECLINED",
		email,
		comment,
		admin,
	)
}
