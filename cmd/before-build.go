package cmd

import (
	"github.com/spf13/cobra"
)

var beforeBuildCmd = &cobra.Command{
	Use:   "before-build",
	Short: "To be run before a build",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	RootCmd.AddCommand(beforeBuildCmd)
}
