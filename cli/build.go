package main

import (
	atlas "github.com/brunoscheufler/atlas/core"
	"github.com/spf13/cobra"
	"os"
)

var buildCmd = &cobra.Command{
	Use: "build",
	Run: func(cmd *cobra.Command, args []string) {
		var stacks []string
		cmd.Flags().StringArrayVarP(&stacks, "stack", "s", nil, "Stack names (required)")

		logger := createLogger()
		cwd, err := os.Getwd()
		if err != nil {
			cmd.PrintErrf("could not create logger: %s", err.Error())
			os.Exit(1)
		}

		err = atlas.Build(cmd.Context(), logger, version, cwd, stacks)
		if err != nil {
			cmd.PrintErrf("could not build stacks: %s", err.Error())
			os.Exit(1)
		}
	},
}

func prepareBuildCmd(rootCmd *cobra.Command) {
	_ = buildCmd.MarkFlagRequired("stack")
	rootCmd.AddCommand(buildCmd)
}
