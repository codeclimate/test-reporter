package xccov

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/env"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
)

var searchPaths = []string{"coverage.json"}

type Formatter struct {
	Path string
}

func (f *Formatter) Search(paths ...string) (string, error) {
	paths = append(paths, searchPaths...)
	for _, p := range paths {
		logrus.Debugf("checking search path %s for xccov formatter", p)
		if _, err := os.Stat(p); err == nil {
			f.Path = p
			return p, nil
		}
	}

	return "", errors.WithStack(errors.Errorf("could not find any files in search paths for xccov. search paths were: %s", strings.Join(paths, ", ")))
}

func (r Formatter) Format() (formatters.Report, error) {
	report, err := formatters.NewReport()
	if err != nil {
		return report, err
	}

	inputXccovFile, err := os.Open(r.Path)
	if err != nil {
		return report, errors.WithStack(errors.Errorf("could not open coverage file %s", r.Path))
	}

	covFile := &xccovFile{}
	err = json.NewDecoder(inputXccovFile).Decode(&covFile)
	if err != nil {
		return report, errors.WithStack(err)
	}

	gitHead, _ := env.GetHead()
	for _, target := range covFile.Targets {
		for _, jsonFile := range target.Files {
			num := 1
			sourceFile, err := formatters.NewSourceFile(jsonFile.Path, gitHead)
			if err != nil {
				logrus.Warnf("Couldn't find file for path \"%s\" from %s coverage data. Ignore if the path doesn't correspond to an existent file in your repo.", jsonFile.Path, r.Path)
				continue
			}

			for _, function := range jsonFile.Functions {
				// fill non executable lines will null
				for num < function.LineNumber {
					sourceFile.Coverage = append(sourceFile.Coverage, formatters.NullInt{})
					num++
				}
				// fill covered lines with 1
				for i := 0; i < function.CoveredLines; i++ {
					sourceFile.Coverage = append(sourceFile.Coverage, formatters.NewNullInt(1))
					num++
				}
				// fill non-covered but executable lines with 0
				nonExecutableLines := function.ExecutableLines - function.CoveredLines
				for i := 0; i < nonExecutableLines; i++ {
					sourceFile.Coverage = append(sourceFile.Coverage, formatters.NewNullInt(0))
					num++
				}
			}

			err = report.AddSourceFile(sourceFile)
			if err != nil {
				return report, errors.WithStack(err)
			}
		}
	}


	return report, nil
}
