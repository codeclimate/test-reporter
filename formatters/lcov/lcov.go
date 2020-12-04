package lcov

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/codeclimate/test-reporter/env"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
)

var searchPaths = []string{"coverage/lcov.info"}

type Formatter struct {
	Path string
}

func (f *Formatter) Search(paths ...string) (string, error) {
	paths = append(paths, searchPaths...)
	for _, p := range paths {
		logrus.Debugf("checking search path %s for lcov formatter", p)
		if _, err := os.Stat(p); err == nil {
			f.Path = p
			return p, nil
		}
	}

	return "", errors.WithStack(errors.Errorf("could not find any files in search paths for lcov. search paths were: %s", strings.Join(paths, ", ")))
}

func (r Formatter) Format() (formatters.Report, error) {
	rep, err := formatters.NewReport()
	if err != nil {
		return rep, err
	}

	b, err := ioutil.ReadFile(r.Path)
	if err != nil {
		return rep, errors.WithStack(err)
	}

	var gitHead, _ = env.GetHead()

	var sf formatters.SourceFile
	curLine := 1

	for _, line := range bytes.Split(b, []byte("\n")) {
		if bytes.HasPrefix(line, []byte("SF:")) {
			name := string(bytes.TrimSpace(bytes.TrimPrefix(line, []byte("SF:"))))
			sf, err = formatters.NewSourceFile(name, gitHead)
			if err != nil {
				return rep, errors.WithStack(err)
			}
			continue
		}
		if bytes.HasPrefix(line, []byte("DA:")) {
			lineInfo := bytes.Split(bytes.TrimSpace(bytes.TrimPrefix(line, []byte("DA:"))), []byte(","))
			ln, err := strconv.Atoi(string(lineInfo[0]))
			if err != nil {
				return rep, errors.WithStack(err)
			}
			for ln-curLine >= 1 {
				sf.Coverage = append(sf.Coverage, formatters.NullInt{})
				curLine++
			}
			lh, err := strconv.Atoi(string(lineInfo[1]))
			if err != nil {
				return rep, errors.WithStack(err)
			}
			sf.Coverage = append(sf.Coverage, formatters.NewNullInt(lh))
			curLine++
			continue
		}
		if bytes.HasPrefix(line, []byte("end_of_record")) {
			err = rep.AddSourceFile(sf)
			if err != nil {
				return rep, errors.WithStack(err)
			}
			curLine = 1
			continue
		}
	}

	return rep, nil
}
