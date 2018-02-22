package cmd

import (
	"os"

	"github.com/codeclimate/test-reporter/upload"
	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var uploadInput string
var uploadOptions = upload.Uploader{}

// uploadCoverageCmd represents the upload command
var uploadCoverageCmd = &cobra.Command{
	Use:   "upload-coverage",
	Short: "Upload pre-formatted coverage payloads to Code Climate servers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if uploadInput == "-" {
			uploadOptions.Input = os.Stdin
		} else {
			uploadOptions.Input, err = os.Open(uploadInput)
			if err != nil {
				return errors.WithStack(err)
			}
		}
		return uploadOptions.Upload()
	},
}

func init() {
	uploadCoverageCmd.Flags().StringVarP(&uploadInput, "input", "i", ccDefaultCoveragePath, "input path")
	uploadCoverageCmd.Flags().StringVarP(&uploadOptions.ReporterID, "id", "r", os.Getenv("CC_TEST_REPORTER_ID"), "reporter identifier")
	uploadCoverageCmd.Flags().StringVarP(&uploadOptions.EndpointURL, "endpoint", "e", envy.Get("CC_TEST_REPORTER_COVERAGE_ENDPOINT", "https://api.codeclimate.com/v1/test_reports"), "endpoint to upload coverage information to")
	uploadCoverageCmd.Flags().IntVarP(&uploadOptions.BatchSize, "batch-size", "s", 500, "batch size for source files")
	uploadCoverageCmd.Flags().BoolVar(&uploadOptions.Insecure, "insecure", false, "send coverage insecurely (without HTTPS)")

	RootCmd.AddCommand(uploadCoverageCmd)
}
