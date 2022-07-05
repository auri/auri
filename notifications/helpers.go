package notifications

import (
	"auri/config"
	"auri/logger"

	"os"
	"os/exec"
	"strings"
	"unicode"

	"golang.org/x/exp/slices"
)

//runAdminEmailTransaction runs the given email admin notification transaction
//if admin email notifications are enabled
func runAdminEmailTransaction(name string, f func() error) error {
	if !config.GetInstance().AdminEmailNotificationsEnable {
		logger.AuriLogger.Infof("Admin email notification: email notifications are disabled, skipping notification for %v", name)
		return nil
	}

	if err := f(); err != nil {
		logger.AuriLogger.Errorf("Admin email notification: %v", err)
		return nil
	}

	logger.AuriLogger.Infof("Admin email notification: admins were successfully notified via email about %v", name)
	return nil
}

//runAdminShellTransaction runs the given shell commands for admin notification transaction
//if shell hook notifications are enabled
func runAdminShellTransaction(
	name,
	hook string,
	hookType string,
	requestEmail string,
	requestComment string,
	adminUser string,
) error {
	if !config.GetInstance().AdminShellNotificationsEnable {
		logger.AuriLogger.Infof("Admin shell notification: shell notifications are disabled, skipping notification for %v", name)
		return nil
	}

	if hook == "" {
		logger.AuriLogger.Errorf("Admin shell notification: shell hook is not defined, skipping notification for %v", name)
		return nil
	}

	cmd := exec.Command(hook)
	cmd.Env = append(os.Environ(),
		"HOOK_TYPE="+sanitizeShellInput(hookType),
		"REQUEST_EMAIL="+sanitizeShellInput(requestEmail),
		"REQUEST_COMMENT="+sanitizeShellInput(requestComment),
		"ADMIN_USER="+sanitizeShellInput(adminUser),
	)
	if err := cmd.Run(); err != nil {
		logger.AuriLogger.Errorf("Admin shell notification: %v", err)
		return nil
	}

	logger.AuriLogger.Infof("Admin shell notification: admins were successfully notified via shell hook about %v", name)
	return nil
}

// sanitizeShellInput ensures that provided input is more or less safe to be used in the shellout context
func sanitizeShellInput(input string) string {
	allowedChars := []rune{'_', '-', '.', ',', '@', ' '}

	s := strings.FieldsFunc(input, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r) && !slices.Contains(allowedChars, r)
	})

	return strings.Join(s, "")
}
