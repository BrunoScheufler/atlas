package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

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
		logLevelEnv = "warn"
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
	prepareUpCmd(rootCmd)
	prepareDownCmd(rootCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(updateCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
