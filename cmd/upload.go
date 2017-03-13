package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type Uploader struct {
	Input       string
	ReporterID  string
	EndpointURL string
}

var uploadOptions = Uploader{}

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		return uploadOptions.Upload()
	},
}

func (u Uploader) Upload() error {
	f, err := os.Open(u.Input)
	if err != nil {
		return err
	}

	c := http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := c.Post(u.EndpointURL, "application/json", f)
	if err != nil {
		return err
	}
	b, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("response from %s was %d: %s", u.EndpointURL, res.StatusCode, string(b))
	}
	fmt.Printf("Status: %d\n", res.StatusCode)
	fmt.Println(string(b))
	return nil
}

func init() {
	uploadCmd.Flags().StringVarP(&uploadOptions.Input, "input", "i", "codeclimate.json", "input path")
	uploadCmd.Flags().StringVarP(&uploadOptions.ReporterID, "id", "r", os.Getenv("CC_TEST_REPORTER_ID"), "reporter identifier")
	uploadCmd.Flags().StringVarP(&uploadOptions.EndpointURL, "endpoint", "e", "https://codeclimate.com/test_reports", "endpoint to upload coverage information to")
	RootCmd.AddCommand(uploadCmd)
}
