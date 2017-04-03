package lcov

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/codeclimate/test-reporter/formatters"
	"github.com/markbates/pop/nulls"
	"github.com/pkg/errors"
)

var searchPaths = []string{"coverage/lcov.info"}

type Formatter struct {
	Path        string
	SourceFiles []formatters.SourceFile
}

func (f *Formatter) Search(paths ...string) (string, error) {
	paths = append(paths, searchPaths...)
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			f.Path = p
			return p, nil
		}
	}

	return "", errors.WithStack(errors.Errorf("could not find any files in search paths for lcov. search paths were: %s", strings.Join(paths, ", ")))
}
func (f *Formatter) Parse() error {
	b, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return errors.WithStack(err)
	}

	var sf formatters.SourceFile
	curLine := 1

	for _, line := range bytes.Split(b, []byte("\n")) {
		if bytes.HasPrefix(line, []byte("SF:")) {
			name := string(bytes.TrimSpace(bytes.TrimPrefix(line, []byte("SF:"))))
			sf, err = formatters.NewSourceFile(name)
			if err != nil {
				return errors.WithStack(err)
			}
			continue
		}
		if bytes.HasPrefix(line, []byte("DA:")) {
			lineInfo := bytes.Split(bytes.TrimSpace(bytes.TrimPrefix(line, []byte("DA:"))), []byte(","))
			ln, err := strconv.Atoi(string(lineInfo[0]))
			if err != nil {
				return errors.WithStack(err)
			}
			for ln-curLine >= 1 {
				sf.Coverage = append(sf.Coverage, nulls.Int{})
				curLine++
			}
			lh, err := strconv.Atoi(string(lineInfo[1]))
			if err != nil {
				return errors.WithStack(err)
			}
			sf.Coverage = append(sf.Coverage, nulls.NewInt(lh))
			curLine++
			continue
		}
		if bytes.HasPrefix(line, []byte("end_of_record")) {
			f.SourceFiles = append(f.SourceFiles, sf)
			curLine = 1
			continue
		}
	}
	return nil
}
