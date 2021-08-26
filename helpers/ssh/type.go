package ssh

import (
	"crypto/ed25519"
	"crypto/rsa"
	"errors"

	"golang.org/x/crypto/ssh"
)

//KeyType defines the type of ssh key
type KeyType uint8

const (
	//RSAKey represents the ssh rsa key type
	RSAKey KeyType = iota
	//Ed25519Key represents the ssh ed25519 key type
	Ed25519Key
)

//DetermineType returns the type of given ssh key and it's length in bits (only RSA)
// and an error if something goes wrong
func DetermineType(sshKey string) (KeyType, int, error) {
	key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(sshKey))
	if err != nil {
		return 0, 0, errors.New("Invalid key format")
	}

	switch keyType := key.(ssh.CryptoPublicKey).CryptoPublicKey().(type) {

	case *rsa.PublicKey:
		return RSAKey, (keyType.Size() * 8), nil

	case ed25519.PublicKey:
		return Ed25519Key, 0, nil

	default:
		return 0, 0, errors.New("Invalid key format")
	}
}
