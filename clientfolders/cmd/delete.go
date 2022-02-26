package cmd

import (
	"fmt"

	"github.com/sasswart/clientfolders/clientfolders/debug"
	"github.com/sasswart/clientfolders/clientfolders/file"
	"github.com/sasswart/clientfolders/clientfolders/find"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func DeleteCmdFactory(logger *zap.Logger, parent *cobra.Command) *cobra.Command {
	cmd := cobra.Command{
		Use:   "delete",
		Short: "Delete a list of files that match a given list of search criteria",
		Run: func(cmd *cobra.Command, args []string) {
			if rootArgs.Debug.Profile {
				debug.Profile(func() {
					delete(logger)
				})
			} else {
				delete(logger)
			}
		},
	}

	parent.AddCommand(&cmd)

	return &cmd
}

func delete(logger *zap.Logger) {
	logger.Info(
		"Running with Root Arguments:",
		zap.Any("source", rootArgs),
	)

	action := NewDeleteAction(*logger, rootArgs)

	subfilesChan := make(chan []string)
	subErrChan := make(chan error)
	patterns := []string{rootArgs.GroupPattern, rootArgs.EntityPattern, rootArgs.YearPattern}
	find.Find(rootArgs.Source, patterns, action, subfilesChan, subErrChan)

	select {
	case files := <-subfilesChan:
		logger.Info(
			"Found files",
			zap.Strings("files", files),
		)
	case err := <-subErrChan:
		logger.Error(
			"Could not delete files",
			zap.Error(fmt.Errorf("could not delete files: %w", err)),
		)
	}
}

func NewDeleteAction(logger zap.Logger, args Args) find.Action {
	return func(path string, errchan chan error) {
		logger.Info(fmt.Sprintf("Deleting %s", path))

		err := file.Delete(path)
		if err != nil {
			errchan <- err
		}

		errchan <- nil
	}
}
