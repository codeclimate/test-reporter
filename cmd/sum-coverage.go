package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type CoverageSummer struct {
	Output string
	Parts  int
}

var summerOptions = CoverageSummer{}

var sumCoverageCmd = &cobra.Command{
	Use:   "sum-coverage",
	Short: "Combine (sum) multiple pre-formatted coverage payloads into one.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you must pass in one or more files to be summarized")
		}
		if summerOptions.Parts != 0 && len(args) != summerOptions.Parts {
			return errors.Errorf("expected %d parts, received only %d parts", summerOptions.Parts, len(args))
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
			err = rep.Merge(&rr)
      if err != nil {
				return errors.WithStack(err)
      }
		}

		err = os.MkdirAll(filepath.Dir(summerOptions.Output), 0755)
		if err != nil {
			return errors.WithStack(err)
		}
		out, err := os.Create(summerOptions.Output)
		if err != nil {
			return errors.WithStack(err)
		}

		err = rep.Save(out)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	},
}

func init() {
	sumCoverageCmd.Flags().IntVarP(&summerOptions.Parts, "parts", "p", 0, "total number of parts to sum")
	sumCoverageCmd.Flags().StringVarP(&summerOptions.Output, "output", "o", ccDefaultCoveragePath, "output path")
	RootCmd.AddCommand(sumCoverageCmd)
}
