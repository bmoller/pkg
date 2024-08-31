package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/bmoller/pkg/aur"
	"github.com/bmoller/pkg/libalpm"
)

var updatesCmd = &cobra.Command{
	Use:   "updates",
	Short: "Check the AUR for updates to locally-installed foreign packages",
	Long: `The updates command queries the AUR for the currently-published version of any
foreign packages (those returned by 'pacman -Qm'). The total count of foreign
packages not found on the AUR is displayed after outputting any available
updates.`,
	Run:  updates,
	Args: cobra.NoArgs,
}

func updates(cmd *cobra.Command, args []string) {
	syncRepos, err := libalpm.GetConfigRepos(libalpm.DefaultConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	foreignPkgs, err := libalpm.GetForeignPackages(libalpm.DefaultRoot, libalpm.DefaultDBPath, syncRepos)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	notFound := 0
	for pkg, ver := range foreignPkgs {
		switch info, err := aur.Info([]string{pkg}); {
		case err != nil:
			fmt.Println(err)
			notFound += 1
		case len(info) != 1:
			notFound += 1
		case libalpm.CompareVersions(ver, info[0].Version) < 0:
			fmt.Printf("%s %s -> %s\n", pkg, ver, info[0].Version)
		}
	}
	if notFound != 0 {
		fmt.Printf("%d packages not found on the AUR.\n", notFound)
	}
}
