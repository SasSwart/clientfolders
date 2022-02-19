package main

import (
	"fmt"
	"os"

	"github.com/sasswart/clientfolders/clientfolders/cmd"
	"go.uber.org/zap"
)

func main() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	logger, _ := config.Build()
	defer logger.Sync()

	rootCmd := cmd.RootCmdFactory(logger)
	cmd.ListCmdFactory(logger, rootCmd)
	cmd.CopyCmdFactory(logger, rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
