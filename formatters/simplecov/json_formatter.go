package simplecov

import (
	"encoding/json"
	"os"
	"fmt"

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

func reportIsOnLegacyFormat(simplecovVersion string) bool {
	var major, minor, patch int
	fmt.Sscanf(simplecovVersion, "%d.%d.%d", &major, &minor, &patch)

	if major < 1 {
		if minor < 18 {
			return true
		}
	}

	return false
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

func jsonFormat(r Formatter, rep formatters.Report) (formatters.Report, error) {
	logrus.Debugf("Analyzing simplecov json output from latest format %s", r.Path)
	jf, err := os.Open(r.Path)
	if err != nil {
		return rep, errors.WithStack(errors.Errorf("could not open coverage file %s", r.Path))
	}

	var m simplecovJsonFormatterReport
	decoder := json.NewDecoder(jf)

	err = decoder.Decode(&m)

	if err != nil {
		return rep, errors.WithStack(err)
	}

	if reportIsOnLegacyFormat(m.Meta.SimpleCovVersion) {
		return rep, errors.WithStack(errors.Errorf("Simplecov report is on legacy format, falling back to legacy formatter."))
	}

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
