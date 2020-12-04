package coveragepy

import (
	"encoding/xml"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/codeclimate/test-reporter/env"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
)

var searchPaths = []string{"coverage.xml"}

type Formatter struct {
	Path string
}

func (f *Formatter) Search(paths ...string) (string, error) {
	paths = append(paths, searchPaths...)
	for _, p := range paths {
		logrus.Debugf("checking search path %s for coverage.py formatter", p)
		if _, err := os.Stat(p); err == nil {
			f.Path = p
			return p, nil
		}
	}

	return "", errors.WithStack(errors.Errorf("could not find any files in search paths for coverage.py. search paths were: %s", strings.Join(paths, ", ")))
}

func (r *Formatter) Format() (formatters.Report, error) {
	rep, err := formatters.NewReport()
	if err != nil {
		return rep, err
	}

	fx, err := os.Open(r.Path)
	if err != nil {
		return rep, errors.WithStack(err)
	}

	coverageFile := &xmlFile{}
	err = xml.NewDecoder(fx).Decode(coverageFile)
	if err != nil {
		return rep, errors.WithStack(err)
	}

	gitHead, _ := env.GetHead()
	for _, xmlPackage := range coverageFile.Packages {
		for _, xmlClass := range xmlPackage.Classes {
			fileName := coverageFile.getFullFilePath(xmlClass.FileName)
			logrus.Debugf("creating test file report for %s", fileName)
			sourceFile, err := formatters.NewSourceFile(fileName, gitHead)
			if err != nil {
				return rep, errors.WithStack(err)
			}
			num := 1
			for _, l := range xmlClass.Lines {
				for num < l.Number {
					sourceFile.Coverage = append(sourceFile.Coverage, formatters.NullInt{})
					num++
				}
				ni := formatters.NewNullInt(l.Hits)
				sourceFile.Coverage = append(sourceFile.Coverage, ni)
				num++
			}
			err = rep.AddSourceFile(sourceFile)
			if err != nil {
				return rep, errors.WithStack(err)
			}
		}
	}

	return rep, nil
}
