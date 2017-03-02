package cmd

import (
	"fmt"

	"github.com/codeclimate/test-reporter/env"
	"github.com/spf13/cobra"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		e, err := env.New()
		if err != nil {
			return err
		}
		fmt.Println(e.String())
		return nil
	},
}

func init() {
	RootCmd.AddCommand(envCmd)
}
