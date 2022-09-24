package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of the Atlas CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Atlas %s\n", version)
	},
}
