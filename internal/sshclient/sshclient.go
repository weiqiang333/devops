package sshclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)


// SSHClient
func SSHCline(privateKey, username, address, port, cmd string) (string, error) {
	key, err := ioutil.ReadFile(privateKey)
	if err != nil {
		log.Printf("unable to read private key: %v", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Printf("unable to parse private key: %v", err)
	}
	// An SSH client is represented with a ClientConn.
	//
	// To authenticate with the remote server you must pass at least one
	// implementation of AuthMethod via the Auth field in ClientConfig,
	// and provide a HostKeyCallback.
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: 5 * time.Second,
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", address, port), config)
	if err != nil {
		log.Printf("Failed to dial for SSHCline: ", err)
		return "", err
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		log.Printf("Failed to create session for SSHCline: ", err)
		return "", err
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(cmd); err != nil {
		log.Printf("Failed to run for SSHCline: %s. %s", cmd, err.Error())
		return "", err
	}
	return b.String(), nil
}