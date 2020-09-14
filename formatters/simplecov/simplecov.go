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

        jf, err := os.Open(r.Path)
        if err != nil {
              return rep, errors.WithStack(errors.Errorf("could not open coverage file %s", r.Path))
        }

        var m simplecovJsonFormatterReport
        decoder := json.NewDecoder(jf)
        decoder.DisallowUnknownFields()

        err = decoder.Decode(&m)

        if err != nil {
              return legacyFormat(r, rep)
        } else {
            logrus.Debugf("Analyzing simplecov json output from latest format %s", r.Path)
            gitHead, _ := env.GetHead()
            for n, ls := range m.CoverageType {
                fe, err := formatters.NewSourceFile(n, gitHead)
                if err != nil {
                        return rep, errors.WithStack(err)
                }
                fe.Coverage = transformLineCoverageToCoverage(ls.LineCoverage)
                err = rep.AddSourceFile(fe)
                if err != nil {
                        return rep, errors.WithStack(err)
                }
            }

            return rep, nil
        }
}

func transformLineCoverageToCoverage(ln []interface{}) formatters.Coverage {
        coverage := make([]formatters.NullInt, len(ln))
        ignoredLine := formatters.NullInt{-1, false}
        var convertedCoverageValue int
        for i:=0; i<len(ln) ; i++ {
          _, ok := ln[i].(string)
          if ok {
              coverage[i] = ignoredLine
          } else {
              if ln[i] == nil {
              coverage[i] = ignoredLine
              } else {
                convertedCoverageValue = int(ln[i].(float64))
                coverage[i] = formatters.NewNullInt(convertedCoverageValue)
              }
          }
        }

        return  coverage
}

type branch struct {
        Type string `json:"type"`
        StartLine int `json:"start_line"`
        EndLine int `json:"end_line"`
        Coverage interface{} `json:"coverage"`
}

type fileCoverage struct {
	LineCoverage []interface{} `json:"lines"`
        Branches []branch `json:"branches"`
}


type meta struct {
        SimpleCovVersion string `json:"simplecov_version"`
}

type simplecovJsonFormatterReport struct {
        Meta meta `json:"meta"`
	CoverageType map[string]fileCoverage `json:"coverage"`
}
