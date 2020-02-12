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
	p.Username = UIDToName[p.Uid]
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
