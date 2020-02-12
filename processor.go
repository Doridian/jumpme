package main

import "path"

type ProcessorMaker func() IProcessor

type IProcessor interface {
	Run() []string
	GetName() string
}

type Processor struct {
	IProcessor
	Name     string
	DoUnique bool
}

type IFileProcessor interface {
	RunFor(val string) []string
}

type FileProcessor struct {
	Processor
	IFileProcessor
	FileNames []string
}

type HomeFileProcessor struct {
	Processor
	IFileProcessor
	FileNames []string
}

type IProcessProcessor interface {
	RunFor(proc *UnixProcess) []string
}

type ProcessProcessor struct {
	Processor
	IProcessProcessor
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
	for k := range uniqueMap {
		result = append(result, k)
	}
	return result
}

func (p *FileProcessor) Run() []string {
	results := make([]string, 0)
	for _, relPath := range p.FileNames {
		results = append(results, p.RunFor(relPath)...)
	}
	return p.uniqueStrings(results)
}

func (p *HomeFileProcessor) Run() []string {
	results := make([]string, 0)
	for _, relPath := range p.FileNames {
		for _, dir := range HomeDirs {
			results = append(results, p.RunFor(path.Join(dir, relPath))...)
		}
	}
	return p.uniqueStrings(results)
}

func (p *ProcessProcessor) Run() []string {
	results := make([]string, 0)

	procs, err := GetProcesses()
	if err != nil {
		panic(err)
	}

	for _, proc := range procs {
		results = append(results, p.RunFor(proc)...)
	}
	return p.uniqueStrings(results)
}
