// Package config provides parsed and validated configuration for application
//
// In order to use it in the code:
// - import config package
// - get initialized configuration data via GetInstance() and use the struct fields
//
// Example:
//
//   import "auri/config"
//
//   var myvar := config.GetInstance().Layout
//
// In order to add new configuration options just add a new struct field with proper type.
//
package config

// Config contains the structure with all configuration options of application
// Tags description:
// - envyname - name of according configuration variable in envy, used as configuration source
// - default - default value used if nothing was set via envy
type Config struct {
	Environment                     string   `envyname:"ENVIRONMENT" default:"production"`
	SessionSecret                   string   `envyname:"SESSION_SECRET"`
	Layout                          string   `envyname:"LAYOUT" default:"application.html"`
	WebTemplatePath                 string   `envyname:"WEB_TEMPLATE_PATH" default:""`
	WebAssetPath                    string   `envyname:"WEB_ASSET_PATH" default:""`
	BaseURL                         string   `envyname:"BASE_URL"`
	ContactEmail                    string   `envyname:"CONTACT_EMAIL"`
	SMTPEnable                      bool     `envyname:"SMTP_ENABLE" default:"true"`
	SMTPHost                        string   `envyname:"SMTP_HOST"`
	SMTPPort                        uint     `envyname:"SMTP_PORT" default:"587"`
	SMTPUser                        string   `envyname:"SMTP_USER" default:""`
	SMTPPassword                    string   `envyname:"SMTP_PASSWORD" default:""`
	SMTPEnforceStartTLS             bool     `envyname:"SMTP_ENFORCE_STARTTLS" default:"true"`
	SMTPTLSSkipVerify               bool     `envyname:"SMTP_TLS_SKIP_VERIFY" default:"false"`
	FromEmail                       string   `envyname:"FROM_EMAIL"`
	NotificationAddresses           []string `envyname:"NOTIFICATION_ADDRESSES"`
	AllowedDomains                  []string `envyname:"ALLOWED_DOMAINS"`
	PasswordResetTimeframe          int      `envyname:"PASSWORD_RESET_TIMEFRAME" default:"2"`
	ProtectedGroups                 []string `envyname:"PROTECTED_GROUPS" default:""`
	InitialCredentialSetupTimeframe int      `envyname:"INITIAL_CREDENTIAL_SETUP_TIMEFRAME" default:"48"`
	EmailValidationTimeframe        int      `envyname:"EMAIL_VALIDATION_TIMEFRAME" default:"24"`
	EmailTemplatePath               string   `envyname:"EMAIL_TEMPLATE_PATH" default:""`
	EmailSubjectApproved            string   `envyname:"EMAIL_SUBJECT_APPROVED" default:"Account request was approved"`
	EmailSubjectDeclined            string   `envyname:"EMAIL_SUBJECT_DECLINED" default:"Account request was declined"`
	EmailSubjectReset               string   `envyname:"EMAIL_SUBJECT_RESET" default:"Password reset for your account"`
	EmailSubjectAccountValidation   string   `envyname:"EMAIL_SUBJECT_ACCOUNT_VALIDATION" default:"Account validation required"`
	EmailSubjectMailValidation      string   `envyname:"EMAIL_SUBJECT_MAIL_VALIDATION" default:"Validation of email address required"`
	EmailSubjectAdminNotification   string   `envyname:"EMAIL_SUBJECT_ADMIN_NOTIFICATION" default:"New account request"`
	AccountValidationTimeframe      int      `envyname:"ACCOUNT_VALIDATION_TIMEFRAME" default:"48"`
	IPAEnable                       bool     `envyname:"IPA_ENABLE" default:"true"`
	IPAHost                         string   `envyname:"IPA_HOST"`
	IPAUser                         string   `envyname:"IPA_USER"`
	IPAPassword                     string   `envyname:"IPA_PASSWORD"`
	IPATLSSkipVerify                bool     `envyname:"IPA_TLS_SKIP_VERIFY" default:"false"`
	RSAKeyLength                    int      `envyname:"RSA_KEY_LENGTH" default:"4096"`
	AllowRSAKeys                    bool     `envyname:"ALLOW_RSA_KEYS" default:"true"`
	AllowED25519Keys                bool     `envyname:"ALLOW_ED25519_KEYS" default:"true"`
	AccountCheckTimeframe           int      `envyname:"ACCOUNT_CHECK_TIMEFRAME" default:"6"`
	AccountCheckReminderAfter       int      `envyname:"ACCOUNT_CHECK_REMINDER_AFTER" default:"2"`
	AccountCheckMaxReminders        int      `envyname:"ACCOUNT_CHECK_MAX_REMINDERS" default:"1"`
	LogIPA                          string   `envyname:"LOG_IPA" default:"/var/log/auri-ipa.log"`
	LogAuri                         string   `envyname:"LOG_AURI" default:"/var/log/auri.log"`
}
