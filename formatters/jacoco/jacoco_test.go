package jacoco

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
	rep, err := f.Format()
	r.NoError(err)
	r.Len(rep.SourceFiles, 3)

	sf := rep.SourceFiles["be/apo/basic/Application.java"]
	r.InDelta(33.3, sf.CoveredPercent, 1)
	r.Len(sf.Coverage, 11)
	r.True(sf.Coverage[6].Valid)
	r.False(sf.Coverage[8].Valid)
	r.Equal(3, sf.Coverage[6].Int)
	r.Equal(0, sf.Coverage[8].Int)
}
