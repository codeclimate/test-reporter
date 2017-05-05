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
	err := f.Parse()
	r.NoError(err)

	r.Len(f.SourceFiles, 4)

	sf := f.SourceFiles[3]
	r.Equal("github.com/codeclimate/test-reporter/formatters/source_file.go", sf.Name)
	r.InDelta(77.4, sf.CoveredPercent, 1)
	r.Len(sf.Coverage, 116)
}
