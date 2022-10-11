package main

import (
	atlas "github.com/brunoscheufler/atlas/core"
	"github.com/spf13/cobra"
	"os"
)

func prepareBuildCmd(rootCmd *cobra.Command) {
	var stacks []string

	var buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Build all artifacts required for stacks",
		Run: func(cmd *cobra.Command, args []string) {
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

	buildCmd.Flags().StringArrayVarP(&stacks, "stacks", "s", nil, "Stack names")
	_ = buildCmd.MarkFlagRequired("stacks")
	rootCmd.AddCommand(buildCmd)
}
