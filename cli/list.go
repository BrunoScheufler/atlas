package main

import (
	"github.com/brunoscheufler/atlas/core"
	"github.com/spf13/cobra"
	"os"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		logger := createLogger()
		cwd, err := os.Getwd()
		if err != nil {
			logger.Fatal(err)
		}

		err = atlas.List(cmd.Context(), logger, version, cwd)
		if err != nil {
			logger.Fatal(err)
		}
	},
}
