/*
 * libalpm.go
 *
 * Copyright (c) 2024 Brandon Moller
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package libalpm

import (
	"bufio"
	"fmt"
	"maps"
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
	if err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("file at %s does not exist", path)
	} else if err != nil {
		return nil, fmt.Errorf("unknown error: %w", err)
	}
	defer file.Close()

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
		return nil, fmt.Errorf("error while reading pacman config: %w", err)
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

func GetForeignPackages(root, dbPath string, repos []string) (pkgs map[string]string, err error) {
	pkgs, err = GetLocalPackages(root, dbPath)
	if err != nil {
		return nil, err
	}
	syncPkgs, err := GetSyncPackages(root, dbPath, repos)
	if err != nil {
		return nil, err
	}

	maps.DeleteFunc(pkgs, func(k, _ string) bool {
		_, ok := syncPkgs[k]
		return ok
	})

	return
}
