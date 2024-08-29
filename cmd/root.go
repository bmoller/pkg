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
	rootCommand.AddCommand(infoCommand)
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
