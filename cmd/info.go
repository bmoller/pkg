package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/bmoller/pkg/aur"
)

var infoCommand = &cobra.Command{
	Use:   "info package",
	Short: "Display details of an AUR package",
	Long: `The info command queries the AUR for details about a user-uploaded package. The
information, if found, is formatted and displayed similar to pacman's output
for official packages.`,
	Args: cobra.ExactArgs(1),
	Run:  info,
}

func info(cmd *cobra.Command, args []string) {
	switch results, err := aur.Info([]string{args[0]}); {
	case err != nil:
		fmt.Println(err)
		os.Exit(1)
	case len(results) != 1:
		fmt.Printf("No package with matching name '%s' found\n", args[0])
	default:
		fmt.Print(results[0].Formatted())
	}
}
