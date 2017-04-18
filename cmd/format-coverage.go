package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/codeclimate/test-reporter/formatters/lcov"
	"github.com/codeclimate/test-reporter/formatters/simplecov"
	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type CoverageFormatter struct {
	In        formatters.Formatter
	InputType string
	Output    string
	Prefix    string
}

var formatOptions = CoverageFormatter{}

// a prioritized list of the formatters to use
var formatterList = []string{"simplecov", "lcov"}

// a map of the formatters to use
var formatterMap = map[string]formatters.Formatter{
	"simplecov": &simplecov.Formatter{},
	"lcov":      &lcov.Formatter{},
}

// formatCoverageCmd represents the format command
var formatCoverageCmd = &cobra.Command{
	Use:   "format-coverage",
	Short: "Locate, parse, and re-format supported coverage sources.",
	RunE: func(cmd *cobra.Command, args []string) error {
		envy.Set("PREFIX", formatOptions.Prefix)
		// if a type is specified use that
		if formatOptions.InputType != "" {
			if f, ok := formatterMap[formatOptions.InputType]; ok {
				logrus.Debugf("using formatter %s", formatOptions.InputType)
				formatOptions.In = f
			} else {
				return errors.WithStack(errors.Errorf("could not find a formatter of type %s", formatOptions.InputType))
			}
		} else {
			logrus.Debug("searching for a formatter to use")
			// else start searching for files:
			for n, f := range formatterMap {
				logrus.Debugf("checking %s formatter", n)
				if p, err := f.Search(); err == nil {
					logrus.Debugf("found file %s for %s formatter", p, n)
					formatOptions.In = f
					break
				}
			}
		}

		if formatOptions.In == nil {
			return errors.WithStack(errors.Errorf("could not find any viable formatter. available formatters: %s", strings.Join(formatterList, ", ")))
		}

		return formatOptions.Save()
	},
}

func (f CoverageFormatter) Save() error {
	err := f.In.Parse()
	if err != nil {
		return errors.WithStack(err)
	}

	var out io.Writer
	if formatOptions.Output == "-" {
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

	rep, err := f.In.Format()
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
	pwd, _ := os.Getwd()
	formatCoverageCmd.Flags().StringVarP(&formatOptions.Prefix, "prefix", "p", pwd, "the root directory where the coverage analysis was performed")
	formatCoverageCmd.Flags().StringVarP(&formatOptions.Output, "output", "o", ccDefaultCoveragePath, "output path")
	formatCoverageCmd.Flags().StringVarP(&formatOptions.InputType, "input-type", "t", "", fmt.Sprintf("type of input source to use [%s]", strings.Join(formatterList, ", ")))
	RootCmd.AddCommand(formatCoverageCmd)
}
