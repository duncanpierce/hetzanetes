package login

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
	"strings"
)

func NewSshKey() (string, string, error) {
	var privateKeyBuffer, publicKeyBuffer strings.Builder

	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return "", "", err
	}
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(&privateKeyBuffer, privateKeyPEM); err != nil {
		return "", "", err
	}
	_, err = publicKeyBuffer.Write(ssh.MarshalAuthorizedKey(publicKey))
	if err != nil {
		return "", "", err
	}

	return publicKeyBuffer.String(), privateKeyBuffer.String(), nil
}
