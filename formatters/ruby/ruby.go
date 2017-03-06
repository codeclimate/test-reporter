package ruby

import (
	"encoding/json"
	"os"
	"sort"
	"time"

	"github.com/markbates/pop/nulls"
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
			Name:      k,
			Timestamp: v.Timestamp.Time(),
			Coverage:  make([]File, 0, len(v.Coverage)),
		}
		for n, ls := range v.Coverage {
			fe := File{
				Name:  n,
				Lines: make([]Line, 0, len(ls)),
			}
			for i, l := range ls {
				fe.Lines = append(fe.Lines, Line{Number: i + 1, Coverage: l})
			}
			tt.Coverage = append(tt.Coverage, fe)
		}
		sort.Slice(tt.Coverage, func(a, b int) bool {
			return tt.Coverage[a].Name < tt.Coverage[b].Name
		})
		f.Tests = append(f.Tests, tt)
	}
	return nil
}

type Test struct {
	Name      string
	Timestamp time.Time
	Coverage  []File
}

type File struct {
	Name  string
	Lines []Line
}

type Line struct {
	Number   int
	Coverage nulls.Int
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
	Timestamp rubyTime               `json:"timestamp"`
	Coverage  map[string][]nulls.Int `json:"coverage"`
}
