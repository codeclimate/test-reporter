package gocov

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/cover"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/env"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/gobuffalo/envy"
	"github.com/markbates/pop/nulls"
	"github.com/pkg/errors"
)

var basePackage string

func init() {
	basePackage, _ = os.Getwd()
	for _, gp := range envy.GoPaths() {
		basePackage = strings.TrimPrefix(basePackage, filepath.Join(gp, "src")+string(os.PathSeparator))
	}
}

var searchPaths = []string{"c.out"}

type Formatter struct {
	Path        string
	SourceFiles []formatters.SourceFile
}

func (f *Formatter) Search(paths ...string) (string, error) {
	paths = append(paths, searchPaths...)
	for _, p := range paths {
		logrus.Debugf("checking search path %s for gocov formatter", p)
		if _, err := os.Stat(p); err == nil {
			f.Path = p
			return p, nil
		}
	}

	return "", errors.WithStack(errors.Errorf("could not find any files in search paths for gocov. search paths were: %s", strings.Join(paths, ", ")))
}

func (f *Formatter) Parse() error {
	profiles, err := cover.ParseProfiles(f.Path)
	if err != nil {
		return errors.WithStack(err)
	}

	gitHead, _ := env.GetHead()
	for _, p := range profiles {
		n := strings.TrimPrefix(p.FileName, basePackage+string(os.PathSeparator))
		sf, err := formatters.NewSourceFile(n, gitHead)
		if err != nil {
			return errors.WithStack(err)
		}
		num := 0
		for _, b := range p.Blocks {
			for num < b.StartLine {
				sf.Coverage = append(sf.Coverage, nulls.Int{})
				num++
			}
			for i := 0; i < b.NumStmt; i++ {
				sf.Coverage = append(sf.Coverage, nulls.NewInt(b.Count))
				num++
			}
		}
		sf.CalcLineCounts()
		f.SourceFiles = append(f.SourceFiles, sf)
	}
	return nil
}

func (r Formatter) Format() (formatters.Report, error) {
	rep, err := formatters.NewReport()
	if err != nil {
		return rep, err
	}

	for _, f := range r.SourceFiles {
		rep.AddSourceFile(f)
	}

	return rep, nil
}
