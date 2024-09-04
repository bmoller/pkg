/*
 * foreign.go
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

package cmd

import (
	"fmt"
	"maps"
	"os"
	"slices"

	"github.com/spf13/cobra"

	"github.com/bmoller/pkg/libalpm"
)

var foreignCmd = &cobra.Command{
	Use:   "foreign",
	Short: "List foreign packages",
	Long: `The foreign command lists locally-installed packages that are not found in any
configured sync repository. This command is equivalent to 'pacman -Qm'.`,
	Run: foreign,
}

func foreign(cmd *cobra.Command, args []string) {
	if pkgs, err := libalpm.GetForeignPackages("/", "/var/lib/pacman", []string{"core", "extra"}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		names := slices.Collect(maps.Keys(pkgs))
		slices.Sort(names)
		for _, pkg := range names {
			fmt.Printf("%s %s\n", pkg, pkgs[pkg])
		}
	}
}
