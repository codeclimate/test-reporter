package cmd

import (
	"encoding/json"
	"fmt"
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
		fmt.Println(1)
		if err != nil {
			return errors.WithStack(err)
		}
		for _, n := range args {
			f, err := os.Open(n)
			fmt.Println(2)
			if err != nil {
				return errors.WithStack(err)
			}
			rr, err := formatters.NewReport()
			fmt.Println(3)
			if err != nil {
				return errors.WithStack(err)
			}
			err = json.NewDecoder(f).Decode(&rr)
			fmt.Println(4)
			if err != nil {
				return errors.WithStack(err)
			}
			fmt.Printf("### rr -> %+v\n", rr)
			rep.Merge(&rr)
		}

		var out io.Writer
		if summerOptions.Print {
			out = os.Stdout
		} else {
			os.MkdirAll(filepath.Dir(summerOptions.Output), 0755)
			out, err = os.Create(summerOptions.Output)
			fmt.Println(5)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		err = rep.Save(out)
		fmt.Println(6)
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
