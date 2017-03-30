package simplecov

import (
	"encoding/json"
	"os"
	"sort"
	"time"

	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
)

type Formatter struct {
	Path  string
	Tests []Test
}

func (f *Formatter) Parse() error {
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
			SourceFiles: make([]SourceFile, 0, len(v.Coverage)),
		}
		for n, ls := range v.Coverage {
			fe := SourceFile{
				Name:     n,
				Coverage: ls,
			}
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
	SourceFiles []SourceFile
}

func New(path string) *Formatter {
	return &Formatter{
		Path:  path,
		Tests: []Test{},
	}
}

type rubyTime int64

func (rt rubyTime) Time() time.Time {
	return time.Unix(int64(rt), 0)
}

type input struct {
	Timestamp rubyTime                       `json:"timestamp"`
	Coverage  map[string]formatters.Coverage `json:"coverage"`
}
