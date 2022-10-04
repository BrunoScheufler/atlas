package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var version string

var rootCmd = &cobra.Command{
	Use:   "atlas",
	Short: "atlas makes local development easy",
	Long:  `a set of tools to make local development easy`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func createLogger() logrus.FieldLogger {
	logLevelEnv := os.Getenv("LOG_LEVEL")
	if logLevelEnv == "" {
		logLevelEnv = "info"
	}

	parsedLevel, err := logrus.ParseLevel(logLevelEnv)
	if err != nil {
		logrus.Fatal(err)
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(parsedLevel)
	return logger
}

func main() {
	ctx := context.Background()

	prepareUpCmd(rootCmd)
	prepareDownCmd(rootCmd)
	prepareBuildCmd(rootCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
