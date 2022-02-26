package cmd

import (
	"fmt"

	"github.com/sasswart/clientfolders/clientfolders/debug"
	"github.com/sasswart/clientfolders/clientfolders/file"
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

	subfilesChan := make(chan []interface{})
	subErrChan := make(chan error)
	patterns := []string{rootArgs.GroupPattern, rootArgs.EntityPattern, rootArgs.YearPattern}
	file.Find(rootArgs.Source, patterns, action, subfilesChan, subErrChan)

	select {
	case files := <-subfilesChan:
		logger.Info(
			"Deleted files",
			zap.Any("files", files),
		)
	case err := <-subErrChan:
		logger.Error(
			"Could not delete files",
			zap.Error(fmt.Errorf("could not delete files: %w", err)),
		)
	}
}

func NewDeleteAction(logger zap.Logger, args Args) file.Action {
	return func(path string) (interface{}, error) {
		logger.Info(fmt.Sprintf("Deleting %s", path))

		err := file.Delete(path)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}
