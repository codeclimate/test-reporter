package cmd

import (
	"io"
	"os"
	"path/filepath"

	"github.com/codeclimate/test-reporter/formatters"
	"github.com/codeclimate/test-reporter/formatters/ruby"
	"github.com/pkg/errors"
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
	if in == nil {
		return errors.New("no coverage found to format")
	}
	err = in.Parse()
	if err != nil {
		return errors.WithStack(err)
	}

	var out io.Writer
	if formatOptions.Print || formatOptions.Output == "-" {
		out = os.Stdout
	} else {
		err = os.MkdirAll(filepath.Dir(formatOptions.Output), 0755)
		if err != nil {
			return errors.WithStack(err)
		}
		out, err = os.Create(formatOptions.Output)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	rep, err := in.Format()
	if err != nil {
		return errors.WithStack(err)
	}

	err = rep.Save(out)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func init() {
	formatCoverageCmd.Flags().BoolVarP(&formatOptions.Print, "print", "p", false, "prints to standard out only")
	formatCoverageCmd.Flags().StringVarP(&formatOptions.Output, "output", "o", ccDefaultCoveragePath, "output path")
	RootCmd.AddCommand(formatCoverageCmd)
}
