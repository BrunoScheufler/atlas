package main

import (
	atlas "github.com/brunoscheufler/atlas/core"
	"github.com/spf13/cobra"
	"os"
)

func prepareDownCmd(rootCmd *cobra.Command) {
	var all bool
	var stacks []string

	var downCmd = &cobra.Command{
		Use:   "down",
		Short: "Stop running service containers and remove volumes and networks for specified or all stacks",
		Run: func(cmd *cobra.Command, args []string) {
			logger := createLogger()
			cwd, err := os.Getwd()
			if err != nil {
				cmd.PrintErrf("could not create logger: %s", err.Error())
				os.Exit(1)
			}

			err = atlas.Down(cmd.Context(), logger, cwd, version, stacks, all)
			if err != nil {
				cmd.PrintErrf("could not up stack: %s", err.Error())
				os.Exit(1)
			}
		},
	}

	downCmd.Flags().StringArrayVarP(&stacks, "stacks", "s", []string{}, "Stack names")
	downCmd.Flags().BoolVarP(&all, "all", "a", false, "Clean up all containers and networks")

	rootCmd.AddCommand(downCmd)
}
