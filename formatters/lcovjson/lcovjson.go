package lcovjson

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/env"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
)

type Formatter struct {
	Path string
}

func (f *Formatter) Search(paths ...string) (string, error) {
	paths = append(paths)
	for _, p := range paths {
		logrus.Debugf("checking search path %s for lcov-json formatter", p)
		if _, err := os.Stat(p); err == nil {
			f.Path = p
			return p, nil
		}
	}

	return "", errors.WithStack(errors.Errorf("could not find any files in search paths for lcov-json. search paths were: %s", strings.Join(paths, ", ")))
}

func (r Formatter) Format() (formatters.Report, error) {
	report, err := formatters.NewReport()
	if err != nil {
		return report, err
	}

	inputLcovJsonFile, err := os.Open(r.Path)
	if err != nil {
		return report, errors.WithStack(errors.Errorf("could not open coverage file %s", r.Path))
	}

	covFile := &lcovJsonFile{}
	err = json.NewDecoder(inputLcovJsonFile).Decode(&covFile)
	if err != nil {
		return report, errors.WithStack(err)
	}

	gitHead, _ := env.GetHead()
	for _, target := range covFile.Data {
		report.CoveredPercent = target.Totals.Lines.Percent
		regionsByFilename := make(map[string][]region)

		for _, function := range target.Functions {
			for _, filename := range function.Filenames {
				// Ignore dependencies.
				if strings.Contains(filename, ".build/checkouts") {
					logrus.Warnf("Ignored dependency file at path \"%s\".", filename)
					continue
				}

				regionsByFilename[filename] = append(regionsByFilename[filename], function.Regions...)
			}
		}

		for filename, regions := range regionsByFilename {
			sourceFile, err := formatters.NewSourceFile(filename, gitHead)
			if err != nil {
				logrus.Warnf("Couldn't find file at path \"%s\" from %s coverage data. Ignore if the path doesn't correspond to an existent file in your repo.", filename, r.Path)
				continue
			}

			coverage := make(map[int]formatters.NullInt)
			lastLine := 1

			for _, region := range regions {
				for line := region.LineStart; line <= region.LineEnd; line++ {
					coverage[line] = formatters.NewNullInt(1)

					if region.ExecutionCount == 0 {
						coverage[line] = formatters.NewNullInt(0)
					}

					if line > lastLine {
						lastLine = line
					}
				}
			}

			for line := 0; line <= lastLine; line++ {
				executionCount, isPresent := coverage[line]

				if isPresent {
					sourceFile.Coverage = append(sourceFile.Coverage, executionCount)
				} else {
					sourceFile.Coverage = append(sourceFile.Coverage, formatters.NullInt{})
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
