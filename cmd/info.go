/*
 * info.go
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
	"os"

	"github.com/spf13/cobra"

	"github.com/bmoller/pkg/aur"
)

var infoCmd = &cobra.Command{
	Use:   "info package",
	Short: "Display details of an AUR package",
	Long: `The info command queries the AUR for details about a user-uploaded package. The
information, if found, is formatted and displayed similar to pacman's output
for official packages.`,
	Args: cobra.ExactArgs(1),
	Run:  info,
}

func info(cmd *cobra.Command, args []string) {
	switch results, err := aur.Info(args); {
	case err != nil:
		fmt.Println(err)
		os.Exit(1)
	case len(results) != 1:
		fmt.Printf("No package with matching name '%s' found\n", args[0])
	default:
		fmt.Print(results[0].Formatted())
	}
}
