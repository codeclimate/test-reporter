package cmd

import (
	"io"
	"os"
	"path/filepath"

	"github.com/codeclimate/test-reporter/formatters"
	"github.com/codeclimate/test-reporter/formatters/ruby"
	"github.com/spf13/cobra"
)

// formatCmd represents the format command
var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		var in formatters.Formatter
		in = ruby.New("./formatters/ruby/ruby-example.json")
		err := in.Parse()
		if err != nil {
			return err
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		var out io.Writer
		if output == "" {
			out = os.Stdout
		} else {
			out, err = os.Open(filepath.Join(output, "codeclimate.json"))
			if err != nil {
				return err
			}
		}

		rep, err := in.Format()
		if err != nil {
			return err
		}

		return rep.Save(out)
	},
}

func init() {
	formatCmd.Flags().StringP("output", "o", "", "output path")
	RootCmd.AddCommand(formatCmd)
}
