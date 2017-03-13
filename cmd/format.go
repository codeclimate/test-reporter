package cmd

import (
	"io"
	"os"
	"path/filepath"

	"github.com/codeclimate/test-reporter/formatters"
	"github.com/codeclimate/test-reporter/formatters/ruby"
	"github.com/spf13/cobra"
)

type Formatter struct {
	Output string
	Print  bool
}

var formatOptions = Formatter{}

// formatCmd represents the format command
var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		return formatOptions.Save()
	},
}

func (f Formatter) Save() error {
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
	formatCmd.Flags().BoolVarP(&formatOptions.Print, "print", "p", false, "prints to standard out only")
	formatCmd.Flags().StringVarP(&formatOptions.Output, "output", "o", "codeclimate.json", "output path")
	RootCmd.AddCommand(formatCmd)
}
