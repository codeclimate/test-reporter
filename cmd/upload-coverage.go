package cmd

import (
	"os"

	"github.com/codeclimate/test-reporter/upload"
	"github.com/gobuffalo/envy"
	"github.com/spf13/cobra"
)

var uploadOptions = upload.Uploader{}

// uploadCoverageCmd represents the upload command
var uploadCoverageCmd = &cobra.Command{
	Use:   "upload-coverage",
	Short: "Upload pre-formatted coverage payloads to Code Climate servers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return uploadOptions.Upload()
	},
}

func init() {
	uploadCoverageCmd.Flags().StringVarP(&uploadOptions.Input, "input", "i", ccDefaultCoveragePath, "input path")
	uploadCoverageCmd.Flags().StringVarP(&uploadOptions.ReporterID, "id", "r", os.Getenv("CC_TEST_REPORTER_ID"), "reporter identifier")
	uploadCoverageCmd.Flags().StringVarP(&uploadOptions.EndpointURL, "endpoint", "e", envy.Get("CC_TEST_REPORTER_COVERAGE_ENDPOINT", "https://codeclimate.com/test_reports"), "endpoint to upload coverage information to")
	RootCmd.AddCommand(uploadCoverageCmd)
}
