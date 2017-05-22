package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var debug bool

const ccDefaultCoveragePath = "coverage/codeclimate.json"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cc-test-reporter",
	Short: "Report information about tests to Code Climate",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		v, err := cmd.Flags().GetBool("version")
		if err != nil {
			return err
		}
		if v {
			fmt.Printf("Code Climate Test Reporter %s\n", version.FormattedVersion())
			return nil
		}
		return cmd.Help()
	},
}

func writer(p string) (io.Writer, error) {
	if p == "-" {
		return os.Stdout, nil
	}
	err := os.MkdirAll(filepath.Dir(p), 0755)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	out, err := os.Create(p)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return out, err
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "run in debug mode")
	RootCmd.Flags().BoolP("version", "v", false, "Show version information")
	logrus.SetLevel(logrus.WarnLevel)
}
