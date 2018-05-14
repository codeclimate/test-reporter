package cobertura

import (
	"encoding/xml"
	"os"
	"sort"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/env"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
)

var searchPaths = []string{"cobertura.xml", "cobertura.ser"}

type Formatter struct {
	Path string
}

func (f *Formatter) Search(paths ...string) (string, error) {
	paths = append(paths, searchPaths...)
	for _, p := range paths {
		logrus.Debugf("checking search path %s for cobertura formatter", p)
		if _, err := os.Stat(p); err == nil {
			f.Path = p
			return p, nil
		}
	}

	return "", errors.WithStack(errors.Errorf("could not find any files in search paths for cobertura. search paths were: %s", strings.Join(paths, ", ")))
}

func (r Formatter) Format() (formatters.Report, error) {
	rep, err := formatters.NewReport()
	if err != nil {
		return rep, err
	}

	fx, err := os.Open(r.Path)
	if err != nil {
		return rep, errors.WithStack(err)
	}

	coberturaFile := &xmlFile{}
	err = xml.NewDecoder(fx).Decode(coberturaFile)

	if err != nil {
		return rep, errors.WithStack(err)
	}

	gitHead, _ := env.GetHead()
	for _, pp := range coberturaFile.Packages {
		mergedClasses := make(map[string]*xmlClass)
		// merge Classes by filename
		for i, clss := range pp.Classes {
			filename := clss.FileName
			if _, ok := mergedClasses[filename]; ok {
				// Appends lines for mergedClasses with the same filename
				lines := append(mergedClasses[filename].Lines, clss.Lines...)
				mergedClasses[filename].Lines = lines
			} else {
				mergedClasses[filename] = &pp.Classes[i]
			}
		}

		for _, pf := range mergedClasses {
			num := 1
			fileName := coberturaFile.getFullFilePath(pf.FileName)
			logrus.Debugf("creating test file report for %s", fileName)
			sf, err := formatters.NewSourceFile(fileName, gitHead)
			if err != nil {
				return rep, errors.WithStack(err)
			}
			sort.Sort(ByLineNum(pf.Lines))
			for _, l := range pf.Lines {
				if l.Num > 0 {
					for num < l.Num {
						sf.Coverage = append(sf.Coverage, formatters.NullInt{})
						num++
					}
					if l.Num <= len(sf.Coverage) {
						hits := sf.Coverage[l.Num-1].Int + l.Hits
						sf.Coverage[l.Num-1] = formatters.NewNullInt(hits)
					} else {
						ni := formatters.NewNullInt(l.Hits)
						sf.Coverage = append(sf.Coverage, ni)
						num++
					}
				} else {
					logrus.Warnf("Invalid line number %d in file %s", l.Num, fileName)
				}
			}
			err = rep.AddSourceFile(sf)
			if err != nil {
				return rep, errors.WithStack(err)
			}
		}
	}

	return rep, nil
}
