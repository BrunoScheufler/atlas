package main

import (
	atlas "github.com/brunoscheufler/atlas/core"
	"github.com/spf13/cobra"
	"os"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop all running service containers and remove volumes and networks",
	Run: func(cmd *cobra.Command, args []string) {
		var stacks []string
		cmd.Flags().StringArrayVarP(&stacks, "stack", "s", nil, "Stack names (required)")

		logger := createLogger()
		cwd, err := os.Getwd()
		if err != nil {
			cmd.PrintErrf("could not create logger: %s", err.Error())
			os.Exit(1)
		}

		err = atlas.Down(cmd.Context(), logger, cwd, stacks)
		if err != nil {
			cmd.PrintErrf("could not up stack: %s", err.Error())
			os.Exit(1)
		}
	},
}

func prepareDownCmd(rootCmd *cobra.Command) {
	_ = downCmd.MarkFlagRequired("stack")
	rootCmd.AddCommand(downCmd)
}
