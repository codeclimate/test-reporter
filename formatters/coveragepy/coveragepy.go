package coveragepy

import (
	"encoding/xml"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/env"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/markbates/pop/nulls"
	"github.com/pkg/errors"
)

var searchPaths = []string{"coverage.xml"}

type Formatter struct {
	Path        string
	SourceFiles []formatters.SourceFile
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

func (f *Formatter) Parse() error {
	fx, err := os.Open(f.Path)
	if err != nil {
		return errors.WithStack(err)
	}

	c := &xmlFile{}
	err = xml.NewDecoder(fx).Decode(c)
	if err != nil {
		return errors.WithStack(err)
	}

	gitHead, _ := env.GetHead()
	for _, pp := range c.Packages {
		for _, cc := range pp.Classes {
			sf, err := formatters.NewSourceFile(cc.FileName, gitHead)
			if err != nil {
				return errors.WithStack(err)
			}
			num := 1
			for _, l := range cc.Lines {
				for num < l.Number {
					sf.Coverage = append(sf.Coverage, nulls.Int{})
					num++
				}
				ni := nulls.NewInt(l.Hits)
				sf.Coverage = append(sf.Coverage, ni)
				num++
			}
			sf.CalcLineCounts()
			f.SourceFiles = append(f.SourceFiles, sf)
		}
	}

	return nil
}

func (r *Formatter) Format() (formatters.Report, error) {
	rep, err := formatters.NewReport()
	if err != nil {
		return rep, err
	}

	for _, f := range r.SourceFiles {
		rep.AddSourceFile(f)
	}

	return rep, nil
}
