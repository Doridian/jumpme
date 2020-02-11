package main

import (
	"bufio"
	"os"
	"strings"
)

type sshKnownHostsFinder struct {
	HomeFileProcessor
}

func MakeSSHKnownHostsFinder() *sshKnownHostsFinder {
	finder := &sshKnownHostsFinder{}
	finder.FileNames = []string{".ssh/known_hosts"}
	finder.HomeFileProcessor.iFileProcessor = finder
	return finder
}

func (p *sshKnownHostsFinder) RunOnFile(absPath string) []string {
	file, err := os.Open(absPath)
	if err != nil {
		return []string{}
	}
	defer file.Close()

	results := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		host := fields[0]
		if len(host) < 1 || host[0] == '#' {
			continue
		}
		if host[0:3] == "|1|" {
			continue
		}
		results = append(results, host)
	}
	return results
}
