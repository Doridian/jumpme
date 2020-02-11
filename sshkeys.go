package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type sshKeyFinder struct {
	HomeFileProcessor
}

func MakeSSHKeyFinder() *sshKeyFinder {
	finder := &sshKeyFinder{}
	finder.FileNames = []string{".ssh/id_rsa", ".ssh/id_dsa", ".ssh/id_ed25519"}
	finder.HomeFileProcessor.iFileProcessor = finder
	return finder
}

func (p *sshKeyFinder) RunOnFile(absPath string) []string {
	file, err := os.Open(absPath)
	if err != nil {
		return []string{}
	}
	defer file.Close()

	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return []string{}
	}

	state := "UNENCRYPTED"
	if strings.Index(string(contents), "ENCRYPTED") > 0 {
		state = "ENCRYPTED"
	}
	return []string{fmt.Sprintf("KEY: %s; STATE: %s", absPath, state)}
}
