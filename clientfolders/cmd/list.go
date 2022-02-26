package cmd

import (
	"fmt"

	"github.com/sasswart/clientfolders/clientfolders/debug"
	"github.com/sasswart/clientfolders/clientfolders/file"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func ListCmdFactory(logger *zap.Logger, parent *cobra.Command) *cobra.Command {

	cmd := cobra.Command{
		Use:   "list",
		Short: "Find a list of files that match a given list of search criteria",
		Long:  `Useful as a test or dry run to see which files will be acted on before running archive or delete`,
		Run: func(cmd *cobra.Command, args []string) {
			if rootArgs.Debug.Profile {
				debug.Profile(func() {
					list(logger)
				})
			} else {
				list(logger)
			}
		},
	}

	parent.AddCommand(&cmd)

	return &cmd
}

func list(logger *zap.Logger) {
	logger.Info(
		"Running with Root Arguments:",
		zap.Any("source", rootArgs),
	)

	subfilesChan := make(chan []string)
	subErrChan := make(chan error)
	patterns := []string{rootArgs.GroupPattern, rootArgs.EntityPattern, rootArgs.YearPattern}
	file.Find(rootArgs.Source, patterns, nil, subfilesChan, subErrChan)

	select {
	case files := <-subfilesChan:
		logger.Info(
			"Found files",
			zap.Strings("files", files),
		)
		return
	case err := <-subErrChan:
		logger.Error(
			"could not list files",
			zap.Error(fmt.Errorf("could not list files: %w", err)),
		)
		return
	}

}
