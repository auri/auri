package ssh

import (
	"errors"
	"strings"
)

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
