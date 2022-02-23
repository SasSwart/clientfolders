package cmd

import (
	"fmt"
	"strings"

	copypkg "github.com/sasswart/clientfolders/clientfolders/copy"
	"github.com/sasswart/clientfolders/clientfolders/debug"
	"github.com/sasswart/clientfolders/clientfolders/find"
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

	patterns := []string{rootArgs.GroupPattern, rootArgs.EntityPattern, rootArgs.YearPattern}
	subfiles, errChan := find.Find(rootArgs.Source, patterns, action)

	var files []string
	var err error
	select {
	case subfiles <- files:
		logger.Info(
			"Found files",
			zap.Strings("files", files),
		)
	case errChan <- err:
		errChan <- fmt.Errorf("could not copy files: %w", err)
		return
	}
}

func NewCopyAction(logger zap.Logger, args Args) find.Action {
	return func(path string, errchan chan error) {
		destination := strings.ReplaceAll(path, args.Source, args.Target)
		logger.Info(fmt.Sprintf("Copying %s to %s", path, destination))

		err := copypkg.Copy(destination, path)
		if err != nil {
			errchan <- err
		}

		errchan <- nil
	}
}
