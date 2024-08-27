package libalpm

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DefaultConfig is the default absolute path to the pacman configuration file.
const DefaultConfig = "/etc/pacman.conf"

// DefaultDBPath is the absolute path to the default location where pacman
// stores its local and repository databases.
const DefaultDBPath = "/var/lib/pacman"

// DefaultRoot is the default root path of the pacman installation.
const DefaultRoot = "/"

/*
GetConfigRepos loads the specified pacman configuration file at path and
extracts any configured repository names. Any error is returned with an
explanatory message in err.
*/
func GetConfigRepos(path string) (names []string, err error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("file at %s does not exist", path)
	} else if err != nil {
		return nil, fmt.Errorf("unknown error: %s", err.Error())
	}

	config := bufio.NewScanner(file)
	for config.Scan() {
		line := config.Text()
		if len(line) > 2 && strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if name := line[1 : len(line)-1]; name != "options" {
				names = append(names, name)
			}
		}
	}
	if err = config.Err(); err != nil {
		return nil, fmt.Errorf("error while reading pacman config: %s", err)
	}

	return
}

func CheckSyncDBs(names []string, dbPath string) (found []string) {
	for _, entry := range names {
		dbPath := filepath.Join(dbPath, "sync", fmt.Sprintf("%s.db", entry))
		if info, err := os.Stat(dbPath); err == nil && !info.IsDir() {
			found = append(found, entry)
		}
	}

	return
}
