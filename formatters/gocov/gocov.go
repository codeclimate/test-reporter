package gocov

import (
	"os"
	"path/filepath"
	"strings"

	"fmt"
	"golang.org/x/tools/cover"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/env"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
)

var basePackage string

func init() {
	basePackage, _ = os.Getwd()
	fmt.Println(basePackage)
	for _, gp := range envy.GoPaths() {
		fmt.Println(filepath.Join(gp, "src")+string(os.PathSeparator))
		basePackage = strings.TrimPrefix(basePackage, filepath.Join(gp, "src")+string(os.PathSeparator))
	}
	fmt.Println(basePackage)
}

var searchPaths = []string{"c.out"}

type Formatter struct {
	Path string
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

func (r Formatter) Format() (formatters.Report, error) {
	rep, err := formatters.NewReport()
	if err != nil {
		return rep, err
	}

	profiles, err := cover.ParseProfiles(r.Path)
	if err != nil {
		return rep, errors.WithStack(err)
	}

	gitHead, _ := env.GetHead()
	for _, p := range profiles {
		fmt.Println(basePackage+string(os.PathSeparator))
		fmt.Println(filepath.FromSlash(p.FileName))
		n := strings.TrimPrefix(filepath.FromSlash(p.FileName), basePackage+string(os.PathSeparator))
		sf, err := formatters.NewSourceFile(n, gitHead)
		if err != nil {
			return rep, errors.WithStack(err)
		}
		blocks := []cover.ProfileBlock{}
		for _, b := range p.Blocks {
			lstIdx := len(blocks) - 1
			if lstIdx < 0 || blocks[lstIdx].StartLine != b.StartLine || blocks[lstIdx].EndLine != b.EndLine {
				blocks = append(blocks, b)
				continue
			}
			blocks[lstIdx].Count += b.Count
		}
		lineNum := 1
		for _, b := range blocks {
			for lineNum < b.StartLine {
				sf.Coverage = append(sf.Coverage, formatters.NullInt{})
				lineNum++
			}
			for lineNum <= b.EndLine {
				sf.Coverage = append(sf.Coverage, formatters.NewNullInt(b.Count))
				lineNum++
			}
		}
		err = rep.AddSourceFile(sf)
		if err != nil {
			return rep, errors.WithStack(err)
		}
	}

	return rep, nil
}
