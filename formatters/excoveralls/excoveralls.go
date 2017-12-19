package excoveralls

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/env"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
)

var searchPaths = []string{"cover/excoveralls.json"}

type Formatter struct {
	Path string
}

func (f *Formatter) Search(paths ...string) (string, error) {
	paths = append(paths, searchPaths...)
	for _, p := range paths {
		logrus.Debugf("checking search path %s for excoveralls formatter", p)
		if _, err := os.Stat(p); err == nil {
			f.Path = p
			return p, nil
		}
	}

	return "", errors.WithStack(errors.Errorf("could not find any files in search paths for excoveralls. search paths were: %s", strings.Join(paths, ", ")))
}

func (r Formatter) Format() (formatters.Report, error) {
	report, err := formatters.NewReport()
	if err != nil {
		return report, err
	}

	inputFile, err := os.Open(r.Path)
	if err != nil {
		return report, errors.WithStack(errors.Errorf("could not open coverage file %s", r.Path))
	}

	coverageInput := &jsonExcoveralls{}
	err = json.NewDecoder(inputFile).Decode(&coverageInput)
	if err != nil {
		return report, errors.WithStack(err)
	}

	gitHead, _ := env.GetHead()
	for _, file := range coverageInput.Files {
		sourceFile, err := formatters.NewSourceFile(file.Name, gitHead)
		if err != nil {
			return report, errors.WithStack(err)
		}
		sourceFile.Coverage = file.Coverage
		err = report.AddSourceFile(sourceFile)
		if err != nil {
			return report, errors.WithStack(err)
		}
	}

	return report, nil
}

type jsonSourceFile struct {
	Name     string              `json:"name,attr"`
	Coverage []formatters.NullInt `json:"coverage,attr"`
}

type jsonExcoveralls struct {
	Files []jsonSourceFile `json:"source_files"`
}
