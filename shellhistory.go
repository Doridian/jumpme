package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type shellHistorySearcher struct {
	HomeFileProcessor
}

var sshLikeCommand = regexp.MustCompilePOSIX("(ssh|scp|sftp|rsync|git remote|git clone) ")
var hostLikePattern = regexp.MustCompilePOSIX("([a-zA-Z0-9]+@)?[a-zA-Z0-9\\-\\_]+\\.[a-zA-Z0-9\\-\\_\\.]+")

func MakeShellHistorySearcher() IProcessor {
	finder := &shellHistorySearcher{}
	finder.Name = "Shell history"
	finder.DoUnique = true
	finder.FileNames = []string{".bash_history", ".zsh_history", ".ash_history"}
	finder.HomeFileProcessor.IFileProcessor = finder
	return finder
}

func (p *shellHistorySearcher) RunFor(absPath string) []string {
	file, err := os.Open(absPath)
	if err != nil {
		return []string{}
	}
	defer file.Close()

	results := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		history := scanner.Text()
		if !sshLikeCommand.MatchString(history) {
			continue
		}
		matches := hostLikePattern.FindAllString(history, -1)
		for _, v := range matches {
			results = append(results, fmt.Sprintf("FILE %s; HOST %s", absPath, v))
		}
	}
	return results
}
