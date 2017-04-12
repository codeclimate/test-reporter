package upload

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/codeclimate/test-reporter/version"
	"github.com/pkg/errors"
)

type Uploader struct {
	Input       string
	ReporterID  string
	EndpointURL string
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
	rep.RepoToken = u.ReporterID

	err = json.NewDecoder(in).Decode(&rep)
	if err != nil {
		return errors.WithStack(err)
	}

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

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("response from %s was %d", u.EndpointURL, res.StatusCode)
	}
	logrus.Infof("Status: %d", res.StatusCode)
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
