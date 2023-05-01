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

func Connect(hostPort string, privateKey string, timeout time.Duration) (*ssh.Client, error) {
	key, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return nil, fmt.Errorf("cannot parse SSH key: %w", err)
	}
	config := &ssh.ClientConfig{
		Config:          ssh.Config{},
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(key)},
		Timeout:         timeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", hostPort, config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to SSH server %s", hostPort)
	}
	return client, nil
}

func RunOne(client *ssh.Client, env map[string]string, command string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	for k, v := range env {
		err = session.Setenv(k, v)
		if err != nil {
			fmt.Printf("error setting env var %s=%s %s", k, v, err.Error())
			return err
		}
	}

	err = session.Run(command)
	if err != nil {
		fmt.Printf("error running %s: %s", command, err.Error())
		return err
	}

	return nil
}

func Run(hostPort string, privateKey string, timeout time.Duration, env map[string]string, commands []string) error {
	client, err := Connect(hostPort, privateKey, timeout)
	if err != nil {
		return err
	}
	defer client.Close()

	for _, command := range commands {
		err = RunOne(client, env, command)
		if err != nil {
			return err
		}
	}

	return nil
}

func AwaitCloudInit(hostPort string, privateKey string) {
	wait := time.NewTicker(10 * time.Second)
	defer wait.Stop()
	for range wait.C {
		fmt.Printf("polling cloudinit on %s\n", hostPort)
		done := PollCloudInit(hostPort, privateKey)
		if done {
			return
		}
	}
}

func PollCloudInit(hostPort string, privateKey string) (done bool) {
	client, err := Connect(hostPort, privateKey, 3*time.Second)
	if err != nil {
		return false
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return false
	}
	defer session.Close()
	err = session.Run("cloud-init status --format=json | jq -e '.status==\"done\"'")
	return err == nil
}
