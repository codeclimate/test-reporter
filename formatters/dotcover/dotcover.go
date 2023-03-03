package dotcover

import (
	"encoding/xml"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/env"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
)

var searchPaths = []string{"dotcover.xml"}

// Formatter is the exported struct to be used on format-coverage.go
type Formatter struct {
	Path string
}

// Search looks for the dotcover test report file in default paths or provided ones.
func (f *Formatter) Search(paths ...string) (string, error) {
	paths = append(paths, searchPaths...)
	for _, p := range paths {
		logrus.Debugf("checking search path %s for dotcover formatter", p)
		if _, err := os.Stat(p); err == nil {
			f.Path = p
			return p, nil
		}
	}

	return "", errors.WithStack(errors.Errorf("could not find any files in search paths for dotcover. search paths were: %s", strings.Join(paths, ", ")))
}

// Format transforms the provided test report into a CC readable report format.
func (f Formatter) Format() (formatters.Report, error) {
	rep, err := formatters.NewReport()
	if err != nil {
		return rep, err
	}

	fx, err := os.Open(f.Path)
	if err != nil {
		return rep, errors.WithStack(err)
	}

	c := &xmlDotCover{}
	err = xml.NewDecoder(fx).Decode(c)
	if err != nil {
		return rep, errors.WithStack(err)
	}

	gitHead, _ := env.GetHead()

	for _, file := range c.Files {
		sf, err := formatters.NewSourceFile(file.Path, gitHead)
		if err != nil {
			return rep, errors.WithStack(err)
		}

		for _, statement := range c.Statements {
			if file.Index == statement.FileIndex {
				if statement.Covered {
					sf.Coverage = append(sf.Coverage, formatters.NewNullInt(1))
				} else {
					sf.Coverage = append(sf.Coverage, formatters.NewNullInt(0))
				}
			}
		}

		err = rep.AddSourceFile(sf)

		if err != nil {
			return rep, errors.WithStack(err)
		}
	}

	return rep, nil
}
