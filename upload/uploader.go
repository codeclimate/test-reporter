package upload

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	Insecure    bool
}

type ErrConflict struct {
	message string
}

func (e *ErrConflict) Error() string {
	return e.message
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
		switch err.(type) {
		case *ErrConflict:
			logrus.Warnf("%s, skipping upload", err.Error())
			return nil
		default:
			return errors.WithStack(err)
		}
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

	postBatchURL, err := u.TransformPostBatchURL(batchLinks.Links.PostBatch)
	if err != nil {
		return errors.WithStack(err)
	}

	return u.SendBatches(testReport, postBatchURL)
}

func (u Uploader) TransformPostBatchURL(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	if u.Insecure {
		parsed.Scheme = "http"
	} else {
		parsed.Scheme = "https"
	}

	return parsed.String(), nil
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
	fmt.Println("Test report uploaded successfully to Code Climate")
	return nil
}

func (u Uploader) doRequest(in io.Reader, url string) (*http.Response, error) {
	c := http.Client{
		Transport: u.newTransport(),
		Timeout:   30 * time.Second,
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
		httpBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		logrus.Debug(string(httpBody))
		errorMessage, err := getErrorMessage(httpBody)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if res.StatusCode == 409 {
			return nil, &ErrConflict{
				message: fmt.Sprintf("Conflict when uploading: %s", errorMessage),
			}
		}

		return res, fmt.Errorf("response from %s.\nHTTP %d: %s", url, res.StatusCode, errorMessage)
	}
	return res, nil
}

func (u Uploader) newTransport() (tr http.RoundTripper) {
	sslCertFile := os.Getenv("SSL_CERT_FILE")
	if sslCertFile == "" {
		return tr
	}

	caCert, err := ioutil.ReadFile(sslCertFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: caCertPool,
		},
	}
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

type apiError struct {
	Detail string `json:"detail"`
}

type errorsResponse struct {
	Errors []apiError `json:"errors"`
}

func getErrorMessage(body []byte) (string, error) {
	var response = new(errorsResponse)
	err := json.Unmarshal(body, &response)
	if err != nil {
		return "", errors.WithStack(err)
	}

	var details []string

	for i := range response.Errors {
		details = append(details, response.Errors[i].Detail)
	}

	return strings.Join(details, ", "), err
}
