package cmd

import (
	"fmt"
	"strings"

	"github.com/sasswart/clientfolders/clientfolders/copy"
	"github.com/sasswart/clientfolders/clientfolders/find"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func CopyCmdFactory(logger *zap.Logger, parent *cobra.Command) *cobra.Command {
	cmd := cobra.Command{
		Use:   "copy",
		Short: "Copy a list of files that match a given list of search criteria",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info(
				"Running with Root Arguments:",
				zap.Any("source", rootArgs),
			)

			action := NewCopyAction(*logger, rootArgs)

			patterns := []string{rootArgs.GroupPattern, rootArgs.EntityPattern, rootArgs.YearPattern}
			files, err := find.Find(rootArgs.Source, patterns, action)
			if err != nil {
				logger.Error(
					"Could not copy files:",
					zap.Error(fmt.Errorf("could not copy files: %w", err)),
				)
			}
			logger.Info(
				"Found files",
				zap.Strings("files", files),
			)
		},
	}

	parent.AddCommand(&cmd)

	return &cmd
}

func NewCopyAction(logger zap.Logger, args Args) func(path string) error {
	return func(path string) error {
		destination := strings.ReplaceAll(path, args.Source, args.Target)
		logger.Info(fmt.Sprintf("Copying %s to %s", path, destination))

		err := copy.Copy(destination, path)
		if err != nil {
			return err
		}

		return nil
	}
}
