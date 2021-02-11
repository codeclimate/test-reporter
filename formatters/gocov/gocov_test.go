package gocov

import (
	"testing"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/codeclimate/test-reporter/env"
	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {

	t.Run("should parse report from single package", func(t *testing.T) {
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
	})

	//
	// Test parsing report that was generated using `-coverpkg` flag.
	// 
	// go test \
	// -coverpkg="github.com/codeclimate/test-reporter/formatters/gocov/example/foo,github.com/codeclimate/test-reporter/formatters/gocov/example/bar" \
	// -coverprofile=example_foobar.out \
	// ./...
	//
	t.Run("should parse coverage report from multiple packages", func(t *testing.T) {
		gb := env.GitBlob
		defer func() { env.GitBlob = gb }()
		env.GitBlob = func(s string, c *object.Commit) (string, error) {
			return s, nil
		}

		r := require.New(t)

		f := &Formatter{Path: "./example_foobar.out"}
		rep, err := f.Format()
		r.NoError(err)

		r.Len(rep.SourceFiles, 2)

		sfFoo := rep.SourceFiles["example/foo/foo.go"]
		sfBar := rep.SourceFiles["example/bar/bar.go"]

		r.EqualValues(87.5, rep.CoveredPercent)
		r.EqualValues(100, sfFoo.CoveredPercent)
		r.InDelta(66.66, sfBar.CoveredPercent, 0.01)
	})

}
