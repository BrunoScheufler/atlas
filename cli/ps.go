package main

import (
	atlas "github.com/brunoscheufler/atlas/core"
	"github.com/spf13/cobra"
	"os"
)

func preparePsCmd(rootCmd *cobra.Command) {
	var stacks []string

	var psCmd = &cobra.Command{
		Use:   "ps",
		Short: "Shows status of running containers",
		Run: func(cmd *cobra.Command, args []string) {
			logger := createLogger()
			cwd, err := os.Getwd()
			if err != nil {
				cmd.PrintErrf("could not create logger: %s", err.Error())
				os.Exit(1)
			}

			err = atlas.Ps(cmd.Context(), logger, cwd, version, stacks)
			if err != nil {
				cmd.PrintErrf("could not build stacks: %s", err.Error())
				os.Exit(1)
			}
		},
	}

	psCmd.Flags().StringArrayVarP(&stacks, "stacks", "s", []string{}, "Stack names")

	rootCmd.AddCommand(psCmd)
}
