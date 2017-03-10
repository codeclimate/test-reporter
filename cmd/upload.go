package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var uploadOptions = struct {
	input       string
	reporterID  string
	endpointURL string
}{}

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := os.Open(uploadOptions.input)
		if err != nil {
			return err
		}

		c := http.Client{
			Timeout: 30 * time.Second,
		}
		res, err := c.Post(uploadOptions.endpointURL, "application/json", f)
		if err != nil {
			return err
		}
		b, _ := ioutil.ReadAll(res.Body)
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return fmt.Errorf("response from %s was %d: %s", uploadOptions.endpointURL, res.StatusCode, string(b))
		}
		fmt.Printf("Status: %d\n", res.StatusCode)
		fmt.Println(string(b))
		return nil
	},
}

func init() {
	uploadCmd.Flags().StringVarP(&uploadOptions.input, "input", "i", "codeclimate.json", "input path")
	uploadCmd.Flags().StringVarP(&uploadOptions.reporterID, "id", "r", os.Getenv("CC_TEST_REPORTER_ID"), "reporter identifier")
	uploadCmd.Flags().StringVarP(&uploadOptions.endpointURL, "endpoint", "e", "https://codeclimate.com/test_reports", "endpoint to upload coverage information to")
	RootCmd.AddCommand(uploadCmd)
}
