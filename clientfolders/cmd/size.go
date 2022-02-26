package cmd

import (
	"fmt"

	"github.com/sasswart/clientfolders/clientfolders/debug"
	"github.com/sasswart/clientfolders/clientfolders/file"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func SizeCmdFactory(logger *zap.Logger, parent *cobra.Command) *cobra.Command {
	cmd := cobra.Command{
		Use:   "size",
		Short: "Size a list of files that match a given list of search criteria",
		Run: func(cmd *cobra.Command, args []string) {
			if rootArgs.Debug.Profile {
				debug.Profile(func() {
					size(logger)
				})
			} else {
				size(logger)
			}
		},
	}

	parent.AddCommand(&cmd)

	return &cmd
}

func size(logger *zap.Logger) {
	logger.Info(
		"Running with Root Arguments:",
		zap.Any("source", rootArgs),
	)

	action := NewSizeAction(*logger, rootArgs)

	subfilesChan := make(chan []interface{})
	subErrChan := make(chan error)
	patterns := []string{rootArgs.GroupPattern, rootArgs.EntityPattern, rootArgs.YearPattern}
	file.Find(rootArgs.Source, patterns, action, subfilesChan, subErrChan)

	select {
	case files := <-subfilesChan:
		var total int64
		for _, file := range files {
			total += file.(int64)
		}
		fmt.Printf("Total Size: %d bytes\n", total)
	case err := <-subErrChan:
		logger.Error(
			"Could not size files",
			zap.Error(fmt.Errorf("could not size files: %w", err)),
		)
	}
}

func NewSizeAction(logger zap.Logger, args Args) file.Action {
	return func(path string) (interface{}, error) {
		logger.Info(fmt.Sprintf("Sizing %s", path))

		size, err := file.Size(path)
		if err != nil {
			return nil, err
		}

		return size, nil
	}
}
