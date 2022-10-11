package main

import (
	atlas "github.com/brunoscheufler/atlas/core"
	"github.com/spf13/cobra"
	"os"
)

func prepareStopCmd(rootCmd *cobra.Command) {
	var stack string

	var stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop service",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger := createLogger()

			cwd, err := os.Getwd()
			if err != nil {
				cmd.PrintErrf("could not create logger: %s", err.Error())
				os.Exit(1)
			}

			err = atlas.Stop(cmd.Context(), logger, version, cwd, stack, args[0])
			if err != nil {
				cmd.PrintErrf("could not stop service: %s", err.Error())
				os.Exit(1)
			}
		},
	}

	stopCmd.Flags().StringVarP(&stack, "stack", "s", "", "Stack name (required)")
	_ = stopCmd.MarkFlagRequired("stack")

	rootCmd.AddCommand(stopCmd)
}
