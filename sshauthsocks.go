package main

import (
	"fmt"
)

type sshAuthSockSearcher struct {
	ProcessProcessor
}

func MakeSSHAuthSockSearcher() IProcessor {
	finder := &sshAuthSockSearcher{}
	finder.Name = "SSH_AUTH_SOCK"
	finder.DoUnique = true
	finder.ProcessProcessor.IProcessProcessor = finder
	return finder
}

func (p *sshAuthSockSearcher) RunFor(proc *UnixProcess) []string {
	sock, ok := proc.Environ["SSH_AUTH_SOCK"]
	if ok {
		return []string{fmt.Sprintf("USER %s; PATH %s", proc.Username, sock)}
	}
	return []string{}
}
