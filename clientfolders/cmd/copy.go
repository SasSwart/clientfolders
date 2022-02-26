package cmd

import (
	"fmt"
	"strings"

	"github.com/sasswart/clientfolders/clientfolders/debug"
	"github.com/sasswart/clientfolders/clientfolders/file"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func CopyCmdFactory(logger *zap.Logger, parent *cobra.Command) *cobra.Command {
	cmd := cobra.Command{
		Use:   "copy",
		Short: "Copy a list of files that match a given list of search criteria",
		Run: func(cmd *cobra.Command, args []string) {
			if rootArgs.Debug.Profile {
				debug.Profile(func() {
					copy(logger)
				})
			} else {
				copy(logger)
			}
		},
	}

	parent.AddCommand(&cmd)

	return &cmd
}

func copy(logger *zap.Logger) {
	logger.Info(
		"Running with Root Arguments:",
		zap.Any("source", rootArgs),
	)

	action := NewCopyAction(*logger, rootArgs)

	subfilesChan := make(chan []interface{})
	subErrChan := make(chan error)
	patterns := []string{rootArgs.GroupPattern, rootArgs.EntityPattern, rootArgs.YearPattern}
	file.Find(rootArgs.Source, patterns, action, subfilesChan, subErrChan)

	select {
	case files := <-subfilesChan:
		logger.Info(
			"Found files",
			zap.Any("files", files),
		)
	case err := <-subErrChan:
		logger.Error(
			"Could not copy files",
			zap.Error(fmt.Errorf("could not copy files: %w", err)),
		)
	}
}

func NewCopyAction(logger zap.Logger, args Args) file.Action {
	return func(path string) (interface{}, error) {
		destination := strings.ReplaceAll(path, args.Source, args.Target)
		logger.Info(fmt.Sprintf("Copying %s to %s", path, destination))

		err := file.Copy(destination, path)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}
