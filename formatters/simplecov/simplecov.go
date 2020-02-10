package simplecov

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/env"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
)

var searchPaths = []string{"coverage/.resultset.json"}

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

	jf, err := os.Open(r.Path)
	if err != nil {
		return rep, errors.WithStack(errors.Errorf("could not open coverage file %s", r.Path))
	}

	m := map[string]input{}
	err = json.NewDecoder(jf).Decode(&m)
	if err != nil {
		return rep, errors.WithStack(err)
	}

	gitHead, _ := env.GetHead()
	for _, v := range m {
		for n, ls := range v.CoverageType {
			fe, err := formatters.NewSourceFile(n, gitHead)
			if err != nil {
				return rep, errors.WithStack(err)
			}
			fe.Coverage = ls.LineCoverage
			err = rep.AddSourceFile(fe)
			if err != nil {
				return rep, errors.WithStack(err)
			}
		}
	}

	return rep, nil
}

type jsonSourceFileCoverage struct {
  LineCoverage formatters.Coverage `json:"lines"`
  BranchCoverage formatters.Coverage `json:"branches"`
}

type input struct {
	CoverageType map[string]jsonSourceFileCoverage `json:"coverage"`
}
