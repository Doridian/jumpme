package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

const PASSWD_FIELD_HOMEDIR = 5

var typeFlag = flag.String("type", "help", "Choose which type of processor to run")

var processorMakers map[string]ProcessorMaker

func main() {
	flag.Parse()
	processorMakers = make(map[string]ProcessorMaker)

	IsRoot = os.Geteuid() == 0
	HomeDirs = getHomeDirs()

	processorMakers["history"] = MakeShellHistorySearcher
	processorMakers["known"] = MakeSSHKnownHostsFinder
	processorMakers["keys"] = MakeSSHKeyFinder

	proc := processorMakers[*typeFlag]()

	log.Printf("Running processor \"%s\"", proc.GetName())

	for _, f := range proc.Run() {
		fmt.Printf("%s\n", f)
	}
}

func getHomeDirs() []string {
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

	// Parse /etc/passwd to find more homes
	passwdStream, err := os.Open("/etc/passwd")
	if err == nil {
		passwdScanner := bufio.NewScanner(passwdStream)
		for passwdScanner.Scan() {
			passwdFields := strings.Split(passwdScanner.Text(), ":")
			addPath(passwdFields[PASSWD_FIELD_HOMEDIR])
		}
		passwdStream.Close()
	}

	keys := make([]string, 0, len(homeFolderSet))
	for k := range homeFolderSet {
		keys = append(keys, k)
	}
	return keys
}
