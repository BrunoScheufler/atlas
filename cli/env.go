package main

import (
	"errors"
	"github.com/brunoscheufler/atlas/core"
	"github.com/spf13/cobra"
	"os"
)

func prepareEnvCmd(rootCmd *cobra.Command) {
	var stack string

	var envCmd = &cobra.Command{
		Use:   "env",
		Short: "Export service environment variables to .env.local",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("requires at least one arg")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger := createLogger()

			cwd, err := os.Getwd()
			if err != nil {
				cmd.PrintErrf("could not create logger: %s", err.Error())
				os.Exit(1)
			}

			err = atlas.Env(cmd.Context(), logger, version, cwd, stack, args[0])
			if err != nil {
				cmd.PrintErrf("could not sync stack env: %s", err.Error())
				os.Exit(1)
			}
		},
	}

	envCmd.Flags().StringVarP(&stack, "stack", "s", "", "Stack name (required)")
	_ = envCmd.MarkFlagRequired("stack")

	rootCmd.AddCommand(envCmd)
}
