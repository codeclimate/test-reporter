package cmd

import (
	"encoding/json"
	"fmt"
	"os"

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
		f, err := cmd.Flags().GetString("format")
		if err != nil {
			return err
		}
		switch f {
		case "json":
			json.NewEncoder(os.Stdout).Encode(e)
		default:
			fmt.Println(e.String())
		}
		return nil
	},
}

func init() {
	envCmd.Flags().StringP("format", "f", "string", "formats the output")
	RootCmd.AddCommand(envCmd)
}
