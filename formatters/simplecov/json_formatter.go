package simplecov

import (
	"encoding/json"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/env"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
)

type branch struct {
	Type      string      `json:"type"`
	StartLine int         `json:"start_line"`
	EndLine   int         `json:"end_line"`
	Coverage  interface{} `json:"coverage"`
}

type fileCoverage struct {
	LineCoverage []interface{} `json:"lines"`
	Branches     []branch      `json:"branches"`
}

type meta struct {
	SimpleCovVersion string `json:"simplecov_version"`
}

type simplecovJsonFormatterReport struct {
	Meta         meta                    `json:"meta"`
	CoverageType map[string]fileCoverage `json:"coverage"`
}

func transformLineCoverageToCoverage(ln []interface{}) formatters.Coverage {
	coverage := make([]formatters.NullInt, len(ln))
	ignoredLine := formatters.NullInt{-1, false}
	var convertedCoverageValue int
	for i := 0; i < len(ln); i++ {
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

	return coverage
}

func jsonFormat(r Formatter, rep formatters.Report) (error) {
	logrus.Debugf("Analyzing simplecov json output from latest format %s", r.Path)
	jf, err := os.Open(r.Path)
	if err != nil {
		return errors.WithStack(errors.Errorf("could not open coverage file %s", r.Path))
	}

	var m simplecovJsonFormatterReport
	decoder := json.NewDecoder(jf)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&m)

	if err != nil {
		return errors.WithStack(err)
	}

	gitHead, _ := env.GetHead()
	for n, ls := range m.CoverageType {
		fe, err := formatters.NewSourceFile(n, gitHead)
		if err != nil {
			return errors.WithStack(err)
		}
		fe.Coverage = transformLineCoverageToCoverage(ls.LineCoverage)
		err = rep.AddSourceFile(fe)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
