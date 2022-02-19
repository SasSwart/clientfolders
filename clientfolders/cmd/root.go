package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Args struct {
	Source        string
	GroupPattern  string
	EntityPattern string
	YearPattern   string
	Target        string
}

var rootArgs Args = Args{}

func RootCmdFactory(logger *zap.Logger) *cobra.Command {

	rootCmd := cobra.Command{
		Use:   "clientfolders",
		Short: "Manage Clientfolders in bulk",
		Long:  `List, archive and delete items within clientfolders in bulk.`,
	}

	rootCmd.PersistentFlags().StringVar(&rootArgs.Source, "source", "", "clientfolders directory")
	rootCmd.PersistentFlags().StringVar(&rootArgs.GroupPattern, "group", ".*", "")
	rootCmd.PersistentFlags().StringVar(&rootArgs.EntityPattern, "entity", ".*", "")
	rootCmd.PersistentFlags().StringVar(&rootArgs.YearPattern, "year", "\\d\\d\\d\\d", "")
	rootCmd.PersistentFlags().StringVar(&rootArgs.Target, "target", "", "directory to copy to")

	return &rootCmd
}
