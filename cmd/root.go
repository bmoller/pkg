/*
 * root.go
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
)

var rootCommand = &cobra.Command{
	Use:   "pkg",
	Short: "this is the short info",
	Long:  "this is the long info",
}

func init() {
	rootCommand.AddCommand(fetchCmd)
	rootCommand.AddCommand(foreignCmd)
	rootCommand.AddCommand(infoCmd)
	rootCommand.AddCommand(searchCmd)
	rootCommand.AddCommand(updatesCmd)
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
