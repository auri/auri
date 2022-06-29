package mail

import (
	"auri/config"
	"reflect"
)

//TokenMailType defined the type of token email
type TokenMailType uint8

const (
	//EMailValidation defined the email sent to the user to validate his email address
	EMailValidation TokenMailType = iota
	//ApprovalMail defines the email about account approval sent to the user
	ApprovalMail
	//PasswordReset defines the email sent by password reset
	PasswordReset
	//AccountValidation defined the email sent to the user to confirm his account
	AccountValidation
)

type tokenMailConf struct {
	template string
	subject  string
}

//tokenMailConfigs sets up the specifics of different mails with tokens
var tokenMailConfigs = map[TokenMailType]tokenMailConf{
	EMailValidation:   {"email_validation.plush.html", config.GetInstance().EmailSubjectMailValidation},
	ApprovalMail:      {"request_approved.plush.html", config.GetInstance().EmailSubjectApproved},
	PasswordReset:     {"password_reset.plush.html", config.GetInstance().EmailSubjectReset},
	AccountValidation: {"account_validation.plush.html", config.GetInstance().EmailSubjectAccountValidation},
}

// TokenMailOpts specifies options which are possible for token mail
type TokenMailOpts struct {
	Login string
	Email string
	Token string
}

// GetMap returns the render.Data representation of options struct
func (to TokenMailOpts) GetMap() map[string]interface{} {
	t := reflect.TypeOf(to)
	v := reflect.ValueOf(to)
	rm := map[string]interface{}{}
	for i := 0; i < t.NumField(); i++ {
		rm[t.Field(i).Name] = v.Field(i).String()
	}
	return rm
}

//SendTokenMail sends a mail of given type to the given email address including given token
func SendTokenMail(opts TokenMailOpts, mailType TokenMailType) error {
	return Send(
		tokenMailConfigs[mailType].subject,
		[]string{opts.Email},
		tokenMailConfigs[mailType].template,
		opts.GetMap(),
	)
}

//SendDecline sends the account rejection mail to the user
func SendDecline(email string) error {
	return Send(
		config.GetInstance().EmailSubjectDeclined,
		[]string{email},
		"request_declined.plush.html",
		map[string]interface{}{
			"userEmail": email,
		},
	)
}

//SendAdminNotification sends the notification about new account request to the admins
func SendAdminNotification(userEmail, comment string) error {
	return Send(
		config.GetInstance().EmailSubjectAdminNotification,
		config.GetInstance().AdminEmailNotificationAddresses,
		"admin_notification.plush.html",
		map[string]interface{}{
			"userEmail": userEmail,
			"comment":   comment,
		},
	)
}
