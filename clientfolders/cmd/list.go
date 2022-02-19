package cmd

import (
	"github.com/sasswart/clientfolders/clientfolders/find"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func ListCmdFactory(logger *zap.Logger, parent *cobra.Command) *cobra.Command {

	cmd := cobra.Command{
		Use:   "list",
		Short: "Find a list of files that match a given list of search criteria",
		Long:  `Useful as a test or dry run to see which files will be acted on before running archive or delete`,
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info(
				"Running with Root Arguments:",
				zap.Any("source", rootArgs),
			)

			patterns := []string{rootArgs.GroupPattern, rootArgs.EntityPattern, rootArgs.YearPattern}
			files, _ := find.Find(rootArgs.Source, patterns, nil)

			logger.Info(
				"Found files",
				zap.Strings("files", files),
			)
		},
	}

	parent.AddCommand(&cmd)

	return &cmd
}
