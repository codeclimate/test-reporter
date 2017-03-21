package cmd

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type CoverageSummer struct {
	Output string
	Print  bool
}

var summerOptions = CoverageSummer{}

var sumCoverageCmd = &cobra.Command{
	Use:   "sum-coverage",
	Short: "Combine (sum) multiple pre-formatted coverage payloads into one.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you must pass in one or more files to be summarized")
		}
		rep, err := formatters.NewReport()
		if err != nil {
			return errors.WithStack(err)
		}
		for _, n := range args {
			f, err := os.Open(n)
			if err != nil {
				return errors.WithStack(err)
			}
			rr, err := formatters.NewReport()
			if err != nil {
				return errors.WithStack(err)
			}
			err = json.NewDecoder(f).Decode(&rr)
			if err != nil {
				return errors.WithStack(err)
			}
			rep.Merge(&rr)
		}

		var out io.Writer
		if summerOptions.Print {
			out = os.Stdout
		} else {
			os.MkdirAll(filepath.Dir(summerOptions.Output), 0755)
			out, err = os.Create(summerOptions.Output)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		err = rep.Save(out)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	},
}

func init() {
	sumCoverageCmd.Flags().BoolVarP(&summerOptions.Print, "print", "p", false, "prints to standard out only")
	sumCoverageCmd.Flags().StringVarP(&summerOptions.Output, "output", "o", "codeclimate.json", "output path")
	RootCmd.AddCommand(sumCoverageCmd)
}
