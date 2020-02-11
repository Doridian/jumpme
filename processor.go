package main

import "path"

var IsRoot bool
var HomeDirs []string

type Processor interface {
	RunOnFile(absPath string) []string
	Run() []string
}

type iFileProcessor interface {
	Processor
	GetFileNames() []string
}

type FileProcessor struct {
	iFileProcessor
	FileNames []string
}

type HomeFileProcessor struct {
	iFileProcessor
	FileNames []string
}

func (p *FileProcessor) Run() []string {
	results := make([]string, 0)
	for _, relPath := range p.FileNames {
		results = append(results, p.RunOnFile(relPath)...)
	}
	return results
}

func (p *HomeFileProcessor) Run() []string {
	results := make([]string, 0)
	for _, relPath := range p.FileNames {
		for _, dir := range HomeDirs {
			results = append(results, p.RunOnFile(path.Join(dir, relPath))...)
		}
	}
	return results
}
