package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

/*
repos, err := libalpm.GetConfigRepos(libalpm.DefaultConfig)
	if err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
	found := libalpm.CheckSyncDBs(repos, libalpm.DefaultDBPath)
	fmt.Println(found)

	result, err := aur.Search("auracle", aur.Name)
	if err != nil {
		fmt.Println(err)
	} else if result.Type == "error" {
		fmt.Println(result.Error)
	} else {
		fmt.Println(result)
	}
*/

var rootCommand = &cobra.Command{
	Use:   "lksjdflkjdsaf",
	Short: "",
	Long:  "",
}

func init() {
	rootCommand.AddCommand(versionCommand)
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
