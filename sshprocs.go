package main

import (
	"fmt"
	"strings"
)

type sshProcsSearcher struct {
	ProcessProcessor
}

func MakeSSHProcsSearcher() IProcessor {
	finder := &sshProcsSearcher{}
	finder.Name = "SSH processes"
	finder.DoUnique = true
	finder.ProcessProcessor.IProcessProcessor = finder
	return finder
}

func (p *sshProcsSearcher) RunFor(proc *UnixProcess) []string {
	if SshLikeCommand.MatchString(proc.Binary) {
		return []string{fmt.Sprintf("USER %s; EXEC %s; CMDLINE %s", proc.Username, proc.Binary, strings.Join(proc.Cmdline, " "))}
	}
	return []string{}
}
