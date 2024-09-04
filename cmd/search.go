/*
 * search.go
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

var searchCmd = &cobra.Command{
	Use:   "search keyword",
	Short: "Search AUR packages for the specified keyword",
	Long: `The search command queries the AUR for any packages matching the given keyword.
Searches can be performed on package names, names and descriptions, maintainer,
etc. The default is to match only on package names.`,
	Args: cobra.ExactArgs(1),
	Run:  search,
}

var searchFlag = ""

func init() {
	searchCmd.PersistentFlags().StringVarP(&searchFlag, "by", "b", "name-desc", "Specify the type of search to perform")
}

func search(cmd *cobra.Command, args []string) {
	var t aur.SearchType
	switch searchFlag {
	case "checkdepends":
		t = aur.CheckDepends
	case "depends":
		t = aur.Depends
	case "maintainer":
		t = aur.Maintainer
	case "makedepends":
		t = aur.MakeDepends
	case "name":
		t = aur.Name
	case "name-desc":
		t = aur.NameDesc
	case "optdepends":
		t = aur.OptDepends
	default:
		fmt.Printf("Unrecognized search type: %s\n\n", searchFlag)
		flag := cmd.Flag("by")
		fmt.Printf("%#v\n", flag)
		cmd.Usage()
		os.Exit(1)
	}
	results, err := aur.Search(args[0], t)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, r := range results {
		fmt.Println(r.Name + " " + r.Version)
		fmt.Println("    " + r.Description)
	}
}
