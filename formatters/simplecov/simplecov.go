package simplecov

import (
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
)

var simpleCovJsonFormatterPath, legacyPath string = "coverage/coverage.json", "coverage/.resultset.json"

var searchPaths = []string{simpleCovJsonFormatterPath, legacyPath}

type Formatter struct {
	Path string
}

func (f *Formatter) Search(paths ...string) (string, error) {
	paths = append(paths, searchPaths...)
	for _, p := range paths {
		logrus.Debugf("checking search path %s for simplecov formatter", p)
		if _, err := os.Stat(p); err == nil {
			f.Path = p
			return p, nil
		}
	}

	return "", errors.WithStack(errors.Errorf("could not find any files in search paths for simplecov. search paths were: %s", strings.Join(paths, ", ")))
}

func (r Formatter) Format() (formatters.Report, error) {
	rep, err := formatters.NewReport()

	if err != nil {
		return rep, err
	}

	rep, err = jsonFormat(r,rep)

	if err != nil {
		return legacyFormat(r, rep)
	}

	return rep, err
}
