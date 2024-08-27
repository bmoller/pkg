package cmd

import (
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use: "aur search",
	Run: runSearch,
}

func runSearch(cmd *cobra.Command, args []string) {

}