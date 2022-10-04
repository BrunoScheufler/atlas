package main

import (
	"github.com/brunoscheufler/atlas/core"
	"github.com/spf13/cobra"
	"os"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Build artifacts, create networks and volumes, and start service containers",
	Run: func(cmd *cobra.Command, args []string) {
		var stacks []string
		cmd.Flags().StringArrayVarP(&stacks, "stack", "s", nil, "Stack names (required)")

		logger := createLogger()
		cwd, err := os.Getwd()
		if err != nil {
			cmd.PrintErrf("could not create logger: %s", err.Error())
			os.Exit(1)
		}

		err = atlas.Up(cmd.Context(), logger, version, cwd, stacks)
		if err != nil {
			cmd.PrintErrf("could not up stack: %s", err.Error())
			os.Exit(1)
		}
	},
}

func prepareUpCmd(rootCmd *cobra.Command) {
	_ = upCmd.MarkFlagRequired("stack")
	rootCmd.AddCommand(upCmd)
}
