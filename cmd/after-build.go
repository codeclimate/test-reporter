package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/codeclimate/test-reporter/upload"
	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var afterBuildOptions = struct {
	InputType   string
	Prefix      string
	BatchSize   int
	EndpointURL string
	ReporterID  string
	ExitCode    int
	Insecure    bool
}{}

var afterBuildCmd = &cobra.Command{
	Use:   "after-build",
	Short: "Locate, parse, and re-format supported coverage sources. Upload pre-formatted coverage payloads to Code Climate servers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if afterBuildOptions.ExitCode != 0 {
			return errors.Errorf("will not run after-build due to previous exit code of %d", afterBuildOptions.ExitCode)
		}

		bb := &bytes.Buffer{}

		cf := CoverageFormatter{
			Prefix:    afterBuildOptions.Prefix,
			InputType: afterBuildOptions.InputType,
			writer:    bb,
		}

		logrus.Debug("about to run format-coverage")
		err := runFormatter(cf)
		if err != nil {
			return errors.WithStack(err)
		}

		uploader := upload.Uploader{
			Input:       bb,
			ReporterID:  afterBuildOptions.ReporterID,
			EndpointURL: afterBuildOptions.EndpointURL,
			BatchSize:   afterBuildOptions.BatchSize,
			Insecure:    afterBuildOptions.Insecure,
		}

		logrus.Debug("about to run upload-coverage")
		return uploader.Upload()
	},
}

func init() {
	pwd, _ := os.Getwd()
	afterBuildCmd.Flags().IntVar(&afterBuildOptions.ExitCode, "exit-code", 0, "exit code of the test run")
	afterBuildCmd.Flags().StringVarP(&afterBuildOptions.Prefix, "prefix", "p", pwd, "the root directory where the coverage analysis was performed")
	afterBuildCmd.Flags().StringVarP(&afterBuildOptions.InputType, "coverage-input-type", "t", "", fmt.Sprintf("type of input source to use [%s]", strings.Join(formatterList, ", ")))
	afterBuildCmd.Flags().StringVarP(&afterBuildOptions.ReporterID, "id", "r", os.Getenv("CC_TEST_REPORTER_ID"), "reporter identifier")
	afterBuildCmd.Flags().StringVarP(&afterBuildOptions.EndpointURL, "coverage-endpoint", "e", envy.Get("CC_TEST_REPORTER_COVERAGE_ENDPOINT", "https://api.codeclimate.com/v1/test_reports"), "endpoint to upload coverage information to")
	afterBuildCmd.Flags().IntVarP(&afterBuildOptions.BatchSize, "batch-size", "s", 500, "batch size for source files")
	afterBuildCmd.Flags().BoolVar(&afterBuildOptions.Insecure, "insecure", false, "send coverage insecurely (without HTTPS)")
	RootCmd.AddCommand(afterBuildCmd)
}
