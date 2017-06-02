package lcov

import (
	"testing"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/codeclimate/test-reporter/env"
	"github.com/stretchr/testify/require"
)

func Test_Formatter_Parse(t *testing.T) {
	gb := env.GitBlob
	defer func() { env.GitBlob = gb }()
	env.GitBlob = func(s string, c *object.Commit) (string, error) {
		return s, nil
	}

	r := require.New(t)
	l := Formatter{Path: "./example.info"}
	rep, err := l.Format()
	r.NoError(err)

	r.Len(rep.SourceFiles, 1)
	sf := rep.SourceFiles["/Users/markbates/Dropbox/development/javascript-test-reporter/formatter.js"]
	r.Len(sf.Coverage, 104)
}

func Test_Format(t *testing.T) {
	gb := env.GitBlob
	defer func() { env.GitBlob = gb }()
	env.GitBlob = func(s string, c *object.Commit) (string, error) {
		return s, nil
	}

	r := require.New(t)

	rb := Formatter{
		Path: "./example.info",
	}
	rep, err := rb.Format()
	r.NoError(err)

	r.InDelta(90.38, rep.CoveredPercent, 1)

	sf := rep.SourceFiles["/Users/markbates/Dropbox/development/javascript-test-reporter/formatter.js"]
	r.InDelta(90.19, sf.CoveredPercent, 1)

	lc := rep.LineCounts
	r.Equal(47, lc.Covered)
	r.Equal(5, lc.Missed)
	r.Equal(52, lc.Total)
}
