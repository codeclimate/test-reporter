package simplecov

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
)

var searchPaths = []string{"coverage/.resultset.json"}

type Formatter struct {
	Path  string
	Tests []Test
}

func (f *Formatter) Search(paths ...string) (string, error) {
	paths = append(paths, searchPaths...)
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			f.Path = p
			return p, nil
		}
	}

	return "", errors.WithStack(errors.Errorf("could not find any files in search paths for simplecov. search paths were: %s", strings.Join(paths, ", ")))
}

func (f *Formatter) Parse() error {
	if f.Path == "" {
		_, err := f.Search()
		if err != nil {
			return errors.WithStack(err)
		}
	}
	jf, err := os.Open(f.Path)
	if err != nil {
		return errors.WithStack(err)
	}
	m := map[string]input{}
	err = json.NewDecoder(jf).Decode(&m)
	if err != nil {
		return errors.WithStack(err)
	}
	f.Tests = make([]Test, 0, len(m))
	for k, v := range m {
		tt := Test{
			Name:        k,
			Timestamp:   v.Timestamp.Time(),
			SourceFiles: make([]formatters.SourceFile, 0, len(v.Coverage)),
		}
		for n, ls := range v.Coverage {
			fe, err := formatters.NewSourceFile(n)
			if err != nil {
				return errors.WithStack(err)
			}
			fe.Coverage = ls
			fe.CalcLineCounts()
			tt.SourceFiles = append(tt.SourceFiles, fe)
		}
		sort.Slice(tt.SourceFiles, func(a, b int) bool {
			return tt.SourceFiles[a].Name < tt.SourceFiles[b].Name
		})
		f.Tests = append(f.Tests, tt)
	}
	return nil
}

type Test struct {
	Name        string
	Timestamp   time.Time
	SourceFiles []formatters.SourceFile
}

type rubyTime int64

func (rt rubyTime) Time() time.Time {
	return time.Unix(int64(rt), 0)
}

type input struct {
	Timestamp rubyTime                       `json:"timestamp"`
	Coverage  map[string]formatters.Coverage `json:"coverage"`
}
