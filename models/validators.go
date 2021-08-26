package models

import (
	"auri/config"
	"auri/helpers/ssh"
	"strconv"

	"regexp"

	"github.com/pkg/errors"
)

//VerifyPassword this func will validate the password strength
// We require it to be 12 characters long and it needs to have at least 3 out of four different character classes:
// lower case letters, upper case letters, numbers or special characters (e.g. ?,!,.,, ...)
func VerifyPassword(password string) error {

	categories := 0

	if len(password) < 12 {
		return errors.New("Password should be at least 12 characters")
	}
	if regexp.MustCompile("[^a-zA-Z0-9]").MatchString(password) {
		categories++
	}
	if regexp.MustCompile("[a-z]").MatchString(password) {
		categories++
	}
	if regexp.MustCompile("[A-Z]").MatchString(password) {
		categories++
	}
	if regexp.MustCompile("[0-9]").MatchString(password) {
		categories++
	}
	if categories < 3 {
		return errors.New("Password do not meet 3 of 4 categories")
	}
	return nil
}

//ValidateSSH this validates the given ssh key
// We only allow RSA or ED25519 keys, depending on what is set in the configuration file. Other key-types will be denied here.
// RSA Keys also need to full fill a length requirement, which can be set via the configuration file
func ValidateSSH(sshKey string) error {
	keyType, length, err := ssh.DetermineType(sshKey)
	if err != nil {
		return err
	}

	switch keyType {
	case ssh.RSAKey:
		if !config.GetInstance().AllowRSAKeys {
			return errors.New("RSA key is not allowed")
		}

		if length < config.GetInstance().RSAKeyLength {
			return errors.New("RSA key is too short, it should be at least " + strconv.Itoa(config.GetInstance().RSAKeyLength) + " bit")
		}

	case ssh.Ed25519Key:
		if !config.GetInstance().AllowED25519Keys {
			return errors.New("ED25519 key is not allowed")
		}

	default:
		return errors.New("Invalid key format")
	}

	return nil
}
