package cobertura

import (
	"testing"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/codeclimate/test-reporter/env"
	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	gb := env.GitBlob
	defer func() { env.GitBlob = gb }()
	env.GitBlob = func(s string, c *object.Commit) (string, error) {
		return s, nil
	}

	r := require.New(t)

	f := &Formatter{Path: "./example.xml"}
	err := f.Parse()
	r.NoError(err)
	r.Len(f.SourceFiles, 4)

	sf := f.SourceFiles[1]
	r.Equal("search/BinarySearch.java", sf.Name)
	r.InDelta(91.6, sf.CoveredPercent, 1)
	r.Len(sf.Coverage, 31)
	r.False(sf.Coverage[2].Valid)
	r.True(sf.Coverage[11].Valid)
	r.Equal(0, sf.Coverage[10].Int)
	r.Equal(3, sf.Coverage[11].Int)
}
