package main

import (
	atlas "github.com/brunoscheufler/atlas/core"
	"github.com/spf13/cobra"
	"os"
)

func prepareStartCmd(rootCmd *cobra.Command) {
	var stack string

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start service",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger := createLogger()

			cwd, err := os.Getwd()
			if err != nil {
				cmd.PrintErrf("could not create logger: %s", err.Error())
				os.Exit(1)
			}

			err = atlas.Start(cmd.Context(), logger, version, cwd, stack, args[0])
			if err != nil {
				cmd.PrintErrf("could not start service: %s", err.Error())
				os.Exit(1)
			}
		},
	}

	startCmd.Flags().StringVarP(&stack, "stack", "s", "", "Stack name (required)")
	_ = startCmd.MarkFlagRequired("stack")

	rootCmd.AddCommand(startCmd)
}
