package login

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

func NewSshKey() (string, string, error) {
	var privateKeyBuffer, publicKeyBuffer strings.Builder

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
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

func Connect(hostPort string, privateKey string, timeout time.Duration, wait time.Duration, maxTries int) (*ssh.Client, *ssh.Session, error) {
	key, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return nil, nil, fmt.Errorf("cannot parse SSH key: %w", err)
	}
	config := &ssh.ClientConfig{
		Config:  ssh.Config{},
		User:    "root",
		Auth:    []ssh.AuthMethod{ssh.PublicKeys(key)},
		Timeout: timeout,
	}
	for try := 0; try < maxTries; try++ {
		client, err := ssh.Dial("tcp", hostPort, config)
		if err == nil {
			session, err := client.NewSession()
			if err == nil {
				return client, session, nil
			}
			client.Close()
		}
		time.Sleep(wait)
	}
	return nil, nil, fmt.Errorf("unable to connect to SSH server %s", hostPort)
}

func Run(hostPort string, privateKey string, timeout time.Duration, wait time.Duration, maxTries int, env map[string]string, commands []string) error {
	client, session, err := Connect(hostPort, privateKey, timeout, wait, maxTries)
	if err != nil {
		return err
	}
	defer client.Close()
	defer session.Close()
	for k, v := range env {
		err = session.Setenv(k, v)
		if err != nil {
			return err
		}
	}
	for _, cmd := range commands {
		err = session.Run(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}
