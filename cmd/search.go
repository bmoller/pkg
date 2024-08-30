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

func search(cmd *cobra.Command, args []string) {
	results, err := aur.Search(args[0], aur.Name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, r := range results {
		fmt.Println(r.Name + " " + r.Version)
		fmt.Println("    " + r.Description)
	}
}
