package excoveralls

import (
	"testing"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/codeclimate/test-reporter/env"
	"github.com/stretchr/testify/require"
)

func Test_Format(t *testing.T) {
	gb := env.GitBlob
	defer func() { env.GitBlob = gb }()
	env.GitBlob = func(s string, c *object.Commit) (string, error) {
		return s, nil
	}

	r := require.New(t)

	rb := Formatter{
		Path: "./excoveralls_example.json",
	}
	rep, err := rb.Format()
	r.NoError(err)

	r.InDelta(93.3, rep.CoveredPercent, 1)

	sf := rep.SourceFiles["demo-app/services/update_user_info.ex"]
	r.InDelta(100, sf.CoveredPercent, 1)
	sfLc := sf.LineCounts
	r.Equal(sfLc.Covered, 9)
	r.Equal(sfLc.Missed, 0)
	r.Equal(sfLc.Total, 9)

	lc := rep.LineCounts
	r.Equal(lc.Covered, 14)
	r.Equal(lc.Missed, 1)
	r.Equal(lc.Total, 15)
}

func Test_Format_MissingFile(t *testing.T) {
	gb := env.GitBlob
	defer func() { env.GitBlob = gb }()
	env.GitBlob = func(s string, c *object.Commit) (string, error) {
		return s, nil
	}

	r := require.New(t)

	rb := Formatter{
		Path: "./not_real.json",
	}
	_, err := rb.Format()
	r.Error(err)
	r.Equal("could not open coverage file ./not_real.json", err.Error())
}
