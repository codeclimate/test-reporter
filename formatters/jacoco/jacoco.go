package jacoco

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codeclimate/test-reporter/env"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
)

var searchPaths = []string{"jacoco.xml"}

type Formatter struct {
	Path string
}

func (f *Formatter) Search(paths ...string) (string, error) {
	paths = append(paths, searchPaths...)
	for _, p := range paths {
		logrus.Debugf("checking search path %s for jacoco formatter", p)
		if _, err := os.Stat(p); err == nil {
			f.Path = p
			return p, nil
		}
	}

	return "", errors.WithStack(errors.Errorf("could not find any files in search paths for jacoco. search paths were: %s", strings.Join(paths, ", ")))
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

	xmlJacoco := &xmlFile{}
	err = xml.NewDecoder(fx).Decode(xmlJacoco)
	if err != nil {
		return rep, errors.WithStack(err)
	}

	gitHead, _ := env.GetHead()
	for _, xmlPackage := range xmlJacoco.Packages {
		for _, xmlSF := range xmlPackage.SourceFile {
			num := 1
			filepath := fmt.Sprintf("%s/%s", xmlPackage.Name, xmlSF.Name)
			sf, err := formatters.NewSourceFile(filepath, gitHead)
			if err != nil {
				return rep, errors.WithStack(err)
			}
			for _, l := range xmlSF.Lines {
				for num < l.Num {
					sf.Coverage = append(sf.Coverage, formatters.NullInt{})
					num++
				}
				ni := formatters.NewNullInt(l.Hits)
				sf.Coverage = append(sf.Coverage, ni)
				num++
			}
			err = rep.AddSourceFile(sf)
			if err != nil {
				return rep, errors.WithStack(err)
			}
		}
	}

	return rep, nil
}
