package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/codeclimate/test-reporter/version"
	"github.com/pkg/errors"
)

type Uploader struct {
	ReporterID  string
	EndpointURL string
	BatchSize   int
	Input       io.Reader
}

func (u Uploader) Upload() error {
	if u.ReporterID == "" {
		return errors.New("you must supply a CC_TEST_REPORTER_ID ENV variable or pass it via the -r flag")
	}

	rep := formatters.Report{
		SourceFiles: formatters.SourceFiles{},
	}

	err := json.NewDecoder(u.Input).Decode(&rep)
	if err != nil {
		return errors.WithStack(err)
	}

	testReport := NewTestReport(rep)

	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		bb := &bytes.Buffer{}
		w := io.MultiWriter(pw, bb)
		err := json.NewEncoder(w).Encode(JSONWraper{Data: testReport})
		if err != nil {
			logrus.Error(err)
			return
		}
		logrus.Debug(bb.String())
	}()

	res, err := u.doRequest(pr, u.EndpointURL)
	if err != nil {
		return errors.WithStack(err)
	}

	batchLinks := struct {
		Links struct {
			PostBatch string `json:"post_batch"`
		} `json:"links"`
	}{}

	err = json.NewDecoder(res.Body).Decode(&batchLinks)
	if err != nil {
		return errors.WithStack(err)
	}

	return u.SendBatches(testReport, batchLinks.Links.PostBatch)
}

func (u Uploader) SendBatches(rep *TestReport, url string) error {
	batch := [][]SourceFile{}

	pos := 0
	count := len(rep.SourceFiles) / u.BatchSize
	remainder := len(rep.SourceFiles) % u.BatchSize
	for i := 0; i < count; i++ {
		end := pos + u.BatchSize
		batch = append(batch, rep.SourceFiles[pos:end])
		pos = end
	}
	if remainder > 0 {
		batch = append(batch, rep.SourceFiles[pos:])
	}

	for i, b := range batch {
		pr, pw := io.Pipe()
		go func() {
			defer pw.Close()
			bb := &bytes.Buffer{}
			w := io.MultiWriter(pw, bb)
			err := json.NewEncoder(w).Encode(JSONWraper{
				Data: b,
				Meta: map[string]int{
					"current": i + 1,
					"total":   len(batch),
				},
			})
			if err != nil {
				logrus.Error(err)
				return
			}
			logrus.Debug(bb.String())
		}()
		_, err := u.doRequest(pr, url)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (u Uploader) doRequest(in io.Reader, url string) (*http.Response, error) {
	c := http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := u.newRequest(in, url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	logrus.Debugf("posting request to %s", url)
	res, err := c.Do(req)
	if err != nil {
		return res, errors.WithStack(err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return res, fmt.Errorf("response from %s was %d", url, res.StatusCode)
	}
	return res, nil
}

func (u Uploader) newRequest(in io.Reader, url string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, in)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header.Set("User-Agent", fmt.Sprintf("TestReporter/%s (Code Climate, Inc.)", version.Version))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CC-Test-Reporter-Id", u.ReporterID)
	req.Header.Set("Accept", "application/vnd.api+json")
	return req, err
}
