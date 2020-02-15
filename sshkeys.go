package main

import (
	"fmt"
	"io/ioutil"
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
	finder.HomeFileProcessor.IFileProcessor = finder
	return finder
}

func (p *sshKeyFinder) RunFor(absPath string) []string {
	contents, err := ioutil.ReadFile(absPath)
	if err != nil {
		return []string{}
	}

	state := "UNENCRYPTED"
	if strings.Index(string(contents), "ENCRYPTED") > 0 {
		state = "ENCRYPTED"
	}
	return []string{fmt.Sprintf("FILE %s; STATE %s", absPath, state)}
}
