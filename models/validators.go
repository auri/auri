package models

import (
	"auri/config"
	"auri/helpers/ssh"
	"strconv"
	"strings"

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

//ConvertPuttySSH tries to convert the given ssh key to the OpenSSH authorized_key format
// returns the key in authorized_key format or error if conversion wasn't possible
// Conversion is done only on the string basis, there is not cryptographic validation of given key
func ConvertPuttySSH(sshKey string) (string, error) {
	header := "---- BEGIN SSH2 PUBLIC KEY ----"
	footer := "---- END SSH2 PUBLIC KEY ----"
	commentHeader := "Comment:"

	sshKey = strings.TrimSpace(sshKey)
	sshKeyLines := strings.Split(sshKey, "\n")
	if len(sshKeyLines) < 3 || sshKeyLines[0] != header || sshKeyLines[len(sshKeyLines)-1] != footer {
		return "", errors.New("Invalid putty public key")
	}

	keyStart := 1
	keyEnd := len(sshKeyLines) - 1
	comment := ""

	commentLine := strings.TrimSpace(sshKeyLines[1])
	if strings.HasPrefix(commentLine, commentHeader) { // found a comment
		keyStart++
		comment = strings.ReplaceAll(commentLine, commentHeader, "")
		comment = strings.TrimSpace(comment)
		comment = strings.Trim(comment, "\"")
		comment = strings.Trim(comment, "'")
	}

	newKey := "ssh-rsa " + strings.Join(sshKeyLines[keyStart:keyEnd], "")
	if comment != "" {
		newKey = newKey + " " + comment
	}
	return newKey, nil
}
