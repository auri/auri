// Package ssh contains different helpers to deal with ssh keys (convert them or determine key types and length)
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
	if len(sshKeyLines) < 3 || strings.TrimSpace(sshKeyLines[0]) != header || strings.TrimSpace(sshKeyLines[len(sshKeyLines)-1]) != footer {
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

	key := ""
	for _, l := range sshKeyLines[keyStart:keyEnd] {
		//we have to append the lines this way, so we can trim possible spaces and CRLF's
		key = key + strings.TrimSpace(l)
	}

	keyType, _, err := DetermineType(key)
	if err != nil {
		return "", err
	}

	keyPrefix := ""

	switch keyType {
	case RSAKey:
		keyPrefix = "ssh-rsa"
	case Ed25519Key:
		keyPrefix = "ssh-ed25519"
	default:
		return "", errors.New("Invalid key")
	}

	newKey := keyPrefix + " " + key
	if comment != "" {
		newKey = newKey + " " + comment
	}
	return newKey, nil
}
