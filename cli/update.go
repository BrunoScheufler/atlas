package main

import (
	"github.com/brunoscheufler/atlas/core"
	"github.com/spf13/cobra"
	"os"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update dependencies of Atlasfiles in the current workspace",
	Run: func(cmd *cobra.Command, args []string) {
		logger := createLogger()
		cwd, err := os.Getwd()
		if err != nil {
			logger.Fatal(err)
		}

		err = atlas.Update(cmd.Context(), logger, cwd)
		if err != nil {
			logger.Fatal(err)
		}
	},
}
