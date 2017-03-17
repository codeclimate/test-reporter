package cmd

import (
	"io"
	"os"
	"path/filepath"

	"github.com/codeclimate/test-reporter/formatters"
	"github.com/codeclimate/test-reporter/formatters/ruby"
	"github.com/spf13/cobra"
)

type CoverageFormatter struct {
	Output string
	Print  bool
}

var formatOptions = CoverageFormatter{}

// formatCoverageCmd represents the format command
var formatCoverageCmd = &cobra.Command{
	Use:   "format-coverage",
	Short: "Locate, parse, and re-format supported coverage sources.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return formatOptions.Save()
	},
}

func (f CoverageFormatter) Save() error {
	var in formatters.Formatter
	_, err := os.Stat("coverage/.resultset.json")
	if err == nil {
		in = ruby.New("coverage/.resultset.json")
	}
	// } else {
	// 	in = ruby.New("./formatters/ruby/ruby-example.json")
	// }
	err = in.Parse()
	if err != nil {
		return err
	}

	var out io.Writer
	if formatOptions.Print {
		out = os.Stdout
	} else {
		os.MkdirAll(filepath.Dir(formatOptions.Output), 0755)
		out, err = os.Create(formatOptions.Output)
		if err != nil {
			return err
		}
	}

	rep, err := in.Format()
	if err != nil {
		return err
	}

	return rep.Save(out)
}

func init() {
	formatCoverageCmd.Flags().BoolVarP(&formatOptions.Print, "print", "p", false, "prints to standard out only")
	formatCoverageCmd.Flags().StringVarP(&formatOptions.Output, "output", "o", "coverage/codeclimate.json", "output path")
	RootCmd.AddCommand(formatCoverageCmd)
}
