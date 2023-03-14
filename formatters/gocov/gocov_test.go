package gocov

import (
	"testing"
	"fmt"

	"path/filepath"

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

		sf := rep.SourceFiles[filepath.FromSlash("github.com/codeclimate/test-reporter/formatters/source_file.go")]

		r.InDelta(75.8, sf.CoveredPercent, 1)
		r.Len(sf.Coverage, 115)
		r.False(sf.Coverage[5].Valid)
		r.True(sf.Coverage[54].Valid)
		r.Equal(0, sf.Coverage[52].Int)
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

		f := &Formatter{Path: "./example/foobar_test.out"}
		rep, err := f.Format()
		r.NoError(err)

		r.Len(rep.SourceFiles, 2)

		fmt.Println(rep.SourceFiles)
		sfFoo := rep.SourceFiles[filepath.Join("example","foo","foo.go")]
		sfBar := rep.SourceFiles[filepath.Join("example","bar","bar.go")]

		r.EqualValues(85, rep.CoveredPercent)
		r.EqualValues(100, sfFoo.CoveredPercent)
		r.InDelta(66.66, sfBar.CoveredPercent, 0.01)
	})
}
