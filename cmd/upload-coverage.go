package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/codeclimate/test-reporter/formatters"
	"github.com/codeclimate/test-reporter/version"
	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type Uploader struct {
	Input       string
	ReporterID  string
	EndpointURL string
}

var uploadOptions = Uploader{}

// uploadCoverageCmd represents the upload command
var uploadCoverageCmd = &cobra.Command{
	Use:   "upload-coverage",
	Short: "Upload pre-formatted coverage payloads to Code Climate servers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return uploadOptions.Upload()
	},
}

func (u Uploader) Upload() error {
	if u.ReporterID == "" {
		return errors.New("you must supply a CC_TEST_REPORTER_ID ENV variable or pass it via the -r flag")
	}

	var err error
	var in io.Reader
	if u.Input == "-" {
		in = os.Stdin
	} else {
		in, err = os.Open(u.Input)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	// we need to read in the JSON file
	// and set the repo_token to whatever
	// is being set at upload time.
	rep, err := formatters.NewReport()
	if err != nil {
		return errors.WithStack(err)
	}

	err = json.NewDecoder(in).Decode(&rep)
	if err != nil {
		return errors.WithStack(err)
	}
	rep.RepoToken = uploadOptions.ReporterID

	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		json.NewEncoder(pw).Encode(rep)
	}()

	c := http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := u.newRequest(pr)
	if err != nil {
		return errors.WithStack(err)
	}

	res, err := c.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}

	io.Copy(os.Stdout, res.Body)

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("response from %s was %d", u.EndpointURL, res.StatusCode)
	}
	fmt.Printf("Status: %d\n", res.StatusCode)
	return nil
}

func (u Uploader) newRequest(in io.Reader) (*http.Request, error) {
	req, err := http.NewRequest("POST", u.EndpointURL, in)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header.Set("User-Agent", fmt.Sprintf("TestReporter/%s (Code Climate, Inc.)", version.Version))
	req.Header.Set("Content-Type", "application/json")
	return req, err
}

func init() {
	uploadCoverageCmd.Flags().StringVarP(&uploadOptions.Input, "input", "i", "coverage/codeclimate.json", "input path")
	uploadCoverageCmd.Flags().StringVarP(&uploadOptions.ReporterID, "id", "r", os.Getenv("CC_TEST_REPORTER_ID"), "reporter identifier")
	uploadCoverageCmd.Flags().StringVarP(&uploadOptions.EndpointURL, "endpoint", "e", envy.Get("CC_TEST_REPORTER_COVERAGE_ENDPOINT", "https://codeclimate.com/test_reports"), "endpoint to upload coverage information to")
	RootCmd.AddCommand(uploadCoverageCmd)
}
