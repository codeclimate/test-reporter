package gocov

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

	f := &Formatter{Path: "./example.out"}
	rep, err := f.Format()
	r.NoError(err)

	r.Len(rep.SourceFiles, 4)

	sf := rep.SourceFiles["github.com/codeclimate/test-reporter/formatters/source_file.go"]
	r.InDelta(77.4, sf.CoveredPercent, 1)
	r.Len(sf.Coverage, 116)
	r.False(sf.Coverage[5].Valid)
	r.True(sf.Coverage[54].Valid)
	r.Equal(0, sf.Coverage[53].Int)
	r.Equal(1, sf.Coverage[55].Int)
}
