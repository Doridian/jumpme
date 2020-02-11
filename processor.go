package main

import "path"

var IsRoot bool
var HomeDirs []string

type ProcessorMaker func() IProcessor

type IProcessor interface {
	RunOnFile(absPath string) []string
	Run() []string
	GetName() string
}

type Processor struct {
	IProcessor
	Name     string
	DoUnique bool
}

type FileProcessor struct {
	Processor
	FileNames []string
}

type HomeFileProcessor struct {
	Processor
	FileNames []string
}

func (p *Processor) GetName() string {
	return p.Name
}

func (p *Processor) uniqueStrings(arr []string) []string {
	if !p.DoUnique {
		return arr
	}

	uniqueMap := make(map[string]bool)
	for _, v := range arr {
		uniqueMap[v] = true
	}
	result := make([]string, 0, len(uniqueMap))
	for k, _ := range uniqueMap {
		result = append(result, k)
	}
	return result
}

func (p *FileProcessor) Run() []string {
	results := make([]string, 0)
	for _, relPath := range p.FileNames {
		results = append(results, p.RunOnFile(relPath)...)
	}
	return p.uniqueStrings(results)
}

func (p *HomeFileProcessor) Run() []string {
	results := make([]string, 0)
	for _, relPath := range p.FileNames {
		for _, dir := range HomeDirs {
			results = append(results, p.RunOnFile(path.Join(dir, relPath))...)
		}
	}
	return p.uniqueStrings(results)
}
