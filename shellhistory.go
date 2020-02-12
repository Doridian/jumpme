package main

import (
	"bufio"
	"fmt"
	"os"
)

type shellHistorySearcher struct {
	HomeFileProcessor
}

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
		if !SshLikeCommand.MatchString(history) {
			continue
		}
		matches := HostLikePattern.FindAllString(history, -1)
		for _, v := range matches {
			results = append(results, fmt.Sprintf("FILE %s; HOST %s", absPath, v))
		}
	}
	return results
}
