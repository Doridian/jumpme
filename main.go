package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const PASSWD_FIELD_HOMEDIR = 5

func main() {
	IsRoot = os.Geteuid() == 0
	HomeDirs = getHomeDirs()

	proc := MakeSSHKnownHostsFinder()
	for _, f := range proc.Run() {
		print(f)
		print("\n")
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
