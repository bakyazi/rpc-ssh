package main

import (
	"golang.org/x/crypto/ssh"
)

func RunSsh(user, host, command string) (*string, error) {
	hostname, passw, err := Hosts.GetHostInfo(host, user)
	if err != nil {
		return nil, err
	}
	client, session, err := connectToHost(user, passw, hostname)
	defer closeSshClient(client)

	if err != nil {
		return nil, err
	}
	out, err := session.CombinedOutput(command)
	if err != nil {
		return nil, err
	}
	strOut := string(out)
	return &strOut, nil
}

func closeSshClient(c *ssh.Client) {
	if c != nil {
		c.Close()
	}
}

func connectToHost(user, pass, host string) (*ssh.Client, *ssh.Session, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}
