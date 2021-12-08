package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

type Host string
type Username string
type Password string

type HostDefinitions map[Host]HostDefinition

type HostDefinition struct {
	Port  uint32
	Users UserCredentials
}

type UserCredentials map[Username]Password

func LoadHosts() (HostDefinitions, error) {
	path := relative2abs()
	hosts := HostDefinitions{}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(b, &hosts); err != nil {
		return nil, err
	}
	return hosts, nil
}

func relative2abs() string {
	_, fileName, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(fileName), "config/hosts.yml")
}

// GetHostInfo returns <host:port>, <password>, error
func (hds *HostDefinitions) GetHostInfo(host, username string) (string, string, error) {
	hd, ok := (*hds)[Host(host)]
	if !ok {
		return "", "", errors.New("unknown host")
	}

	passw, ok := hd.Users[Username(username)]
	if !ok {
		return "", "", errors.New("unknown username")
	}

	return fmt.Sprintf("%s:%d", host, hd.Port),
		string(passw),
		nil
}
