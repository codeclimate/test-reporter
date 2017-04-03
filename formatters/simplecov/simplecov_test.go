package simplecov

import (
	"testing"

	"github.com/codeclimate/test-reporter/env"
	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	ogb := env.GitBlob
	defer func() {
		env.GitBlob = ogb
	}()
	env.GitBlob = func(s string) (string, error) {
		return s, nil
	}

	r := require.New(t)

	f := Formatter{
		Path: "./simplecov-example.json",
	}
	err := f.Parse()
	r.NoError(err)

	r.Len(f.Tests, 1)

	tt := f.Tests[0]
	r.Equal("Unit Tests", tt.Name)
	r.NotZero(tt.Timestamp)

	r.Len(tt.SourceFiles, 7)

	cf := tt.SourceFiles[6]
	r.Equal("development/mygem/lib/mygem/wrap.rb", cf.Name)
	r.Len(cf.Coverage, 10)
	for i, x := range []interface{}{1, nil, 1, 17, 20, 16, 16, 12, nil, nil} {
		l := cf.Coverage[i]
		r.Equal(x, l.Interface())
	}
}
