package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

const PASSWD_FIELD_USERNAME = 0
const PASSWD_FIELD_UID = 2
const PASSWD_FIELD_HOMEDIR = 5

var typeFlag = flag.String("type", "all", "Choose which type of processor to run")

var processorMakers map[string]ProcessorMaker

func main() {
	flag.Parse()
	processorMakers = make(map[string]ProcessorMaker)

	loadData()

	processorMakers["history"] = MakeShellHistorySearcher
	processorMakers["known"] = MakeSSHKnownHostsFinder
	processorMakers["keys"] = MakeSSHKeyFinder
	processorMakers["procs"] = MakeSSHProcsSearcher

	typeProc := *typeFlag
	if typeProc == "all" {
		for _, proc := range processorMakers {
			runProc(proc)
		}
		return
	}
	runProc(processorMakers[typeProc])
}

func runProc(procMaker ProcessorMaker) {
	proc := procMaker()
	log.Printf("Running processor \"%s\"", proc.GetName())
	for _, f := range proc.Run() {
		fmt.Printf("%s\n", f)
	}
}

func loadData() {
	IsRoot = os.Geteuid() == 0

	homeFolderSet := make(map[string]bool)

	addPath := func(subPath string) {
		homeFolderSet[path.Clean(subPath)] = true
	}

	// Grab all home folders from /home
	dirs, err := ioutil.ReadDir("/home")
	if err == nil {
		for _, dir := range dirs {
			addPath(path.Join("/home", dir.Name()))
		}
	}

	// Add /root to be mix
	addPath("/root")

	// Add current user's homedir
	addPath(os.Getenv("HOME"))

	UIDToName = make(map[int]string)

	// Parse /etc/passwd to find more homes
	passwdStream, err := os.Open("/etc/passwd")
	if err == nil {
		passwdScanner := bufio.NewScanner(passwdStream)
		for passwdScanner.Scan() {
			passwdFields := strings.Split(passwdScanner.Text(), ":")
			uid, err := strconv.ParseInt(passwdFields[PASSWD_FIELD_UID], 10, 32)
			if err == nil {
				UIDToName[int(uid)] = passwdFields[PASSWD_FIELD_USERNAME]
			}
			addPath(passwdFields[PASSWD_FIELD_HOMEDIR])
		}
		passwdStream.Close()
	}

	HomeDirs = make([]string, 0, len(homeFolderSet))
	for k := range homeFolderSet {
		HomeDirs = append(HomeDirs, k)
	}
}
