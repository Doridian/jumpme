package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type UnixProcess struct {
	Pid      int
	Ppid     int
	State    rune
	Pgrp     int
	Sid      int
	Binary   string
	Uid      int
	Username string
	Gid      int
	Cmdline  []string
	Environ  map[string]string
}

var processes []*UnixProcess

func GetProcesses() ([]*UnixProcess, error) {
	if processes != nil {
		return nil, nil
	}
	processes, err := getProcesses()
	return processes, err
}

func (p *UnixProcess) Refresh() error {
	statPath := fmt.Sprintf("/proc/%d/stat", p.Pid)
	stat, err := os.Stat(statPath)
	if err != nil {
		return err
	}
	dataBytes, err := ioutil.ReadFile(statPath)
	if err != nil {
		return err
	}

	statConv := stat.Sys().(*syscall.Stat_t)
	p.Uid = int(statConv.Uid)
	uname, ok := UIDToName[p.Uid]
	if ok {
		p.Username = uname
	} else {
		p.Username = strconv.FormatInt(int64(p.Uid), 10)
	}
	p.Gid = int(statConv.Gid)

	// First, parse out the image name
	data := string(dataBytes)
	binStart := strings.IndexRune(data, '(') + 1
	binEnd := strings.IndexRune(data[binStart:], ')')
	p.Binary = data[binStart : binStart+binEnd]

	// Move past the image name and start parsing the rest
	data = data[binStart+binEnd+2:]
	_, err = fmt.Sscanf(data,
		"%c %d %d %d",
		&p.State,
		&p.Ppid,
		&p.Pgrp,
		&p.Sid)

	cmdline, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", p.Pid))
	if err != nil {
		return err
	}

	p.Cmdline = strings.Split(string(cmdline), "\000")

	environ, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/environ", p.Pid))
	if err != nil {
		return err
	}

	p.Environ = make(map[string]string)
	environSplit := strings.Split(string(environ), "\000")
	for _, v := range environSplit {
		varSplit := strings.SplitN(v, "=", 2)
		if len(varSplit) < 2 {
			p.Environ[varSplit[0]] = ""
		} else {
			p.Environ[varSplit[0]] = varSplit[1]
		}
	}

	return err
}

func getProcesses() ([]*UnixProcess, error) {
	d, err := os.Open("/proc")
	if err != nil {
		panic(err)
	}
	defer d.Close()

	results := make([]*UnixProcess, 0, 50)
	for {
		names, err := d.Readdirnames(10)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		for _, name := range names {
			// We only care if the name starts with a numeric
			if name[0] < '0' || name[0] > '9' {
				continue
			}

			// From this point forward, any errors we just ignore, because
			// it might simply be that the process doesn't exist anymore.
			pid, err := strconv.ParseInt(name, 10, 0)
			if err != nil {
				continue
			}

			p, err := newUnixProcess(int(pid))
			if err != nil {
				continue
			}

			results = append(results, p)
		}
	}

	return results, nil
}

func newUnixProcess(pid int) (*UnixProcess, error) {
	p := &UnixProcess{Pid: pid}
	return p, p.Refresh()
}
