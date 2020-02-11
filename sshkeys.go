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

func MakeSSHKeyFinder() IProcessor {
	finder := &sshKeyFinder{}
	finder.Name = "SSH keys"
	finder.DoUnique = true
	finder.FileNames = []string{".ssh/id_rsa", ".ssh/id_dsa", ".ssh/id_ed25519"}
	finder.HomeFileProcessor.IProcessor = finder
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
