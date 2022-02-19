package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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

			action := func(path string) error {
				destination := strings.ReplaceAll(path, rootArgs.Source, rootArgs.Target)
				logger.Info(fmt.Sprintf("Copying %s to %s", path, destination))

				err := os.MkdirAll(filepath.Dir(destination), os.ModePerm)
				if err != nil {
					return fmt.Errorf("could not create destination parent directories: %w", err)
				}

				sourceFile, err := os.Open(path)
				if err != nil {
					return fmt.Errorf("could not open source file: %w", err)
				}

				destinationFile, err := os.OpenFile(destination, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
				if err != nil {
					return fmt.Errorf("could not open destination file: %w", err)
				}

				_, err = io.Copy(destinationFile, sourceFile)
				if err != nil {
					return fmt.Errorf("could not copy file: %w", err)
				}
				return nil
			}

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

	cmd.PersistentFlags().StringVar(&rootArgs.Target, "target", "", "")

	parent.AddCommand(&cmd)

	return &cmd
}
