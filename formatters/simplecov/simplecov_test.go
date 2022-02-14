package simplecov

import (
	"testing"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/codeclimate/test-reporter/env"
	"github.com/stretchr/testify/require"
)

func Test_ParseLegacy(t *testing.T) {
	ogb := env.GitBlob
	defer func() {
		env.GitBlob = ogb
	}()
	env.GitBlob = func(s string, c *object.Commit) (string, error) {
		return s, nil
	}

	assert := require.New(t)

	formatter := Formatter{
		Path: "./simplecov-example-legacy-resultset.json",
	}
	rep, err := formatter.Format()
	assert.NoError(err)

	assert.Len(rep.SourceFiles, 7)

	cf := rep.SourceFiles["development/mygem/lib/mygem/wrap.rb"]
	assert.Len(cf.Coverage, 10)
	for i, x := range []interface{}{1, nil, 1, 17, 20, 16, 16, 12, nil, nil} {
		l := cf.Coverage[i]
		assert.Equal(x, l.Interface())
	}
}

func Test_Parse(t *testing.T) {
	ogb := env.GitBlob
	defer func() {
		env.GitBlob = ogb
	}()
	env.GitBlob = func(s string, c *object.Commit) (string, error) {
		return s, nil
	}

	assert := require.New(t)

	formatter := Formatter{
		Path: "./simplecov-simple-example.json",
	}
	rep, err := formatter.Format()
	assert.NoError(err)

	assert.Len(rep.SourceFiles, 7)

	cf := rep.SourceFiles["development/mygem/lib/mygem/wrap.rb"]
	assert.Len(cf.Coverage, 10)
	for i, x := range []interface{}{1, nil, 1, 17, 20, 16, 16, 12, nil, nil} {
		l := cf.Coverage[i]
		assert.Equal(x, l.Interface())
	}
}

func Test_ParseWithBranch(t *testing.T) {
	ogb := env.GitBlob
	defer func() {
		env.GitBlob = ogb
	}()
	env.GitBlob = func(s string, c *object.Commit) (string, error) {
		return s, nil
	}

	assert := require.New(t)

	formatter := Formatter{
		Path: "./simplecov-branch-example.json",
	}
	rep, err := formatter.Format()
	assert.NoError(err)

	assert.Len(rep.SourceFiles, 7)

	cf := rep.SourceFiles["development/mygem/lib/mygem/wrap.rb"]
	assert.Len(cf.Coverage, 10)
	for i, x := range []interface{}{1, nil, 1, 17, 20, 16, 16, 12, nil, nil} {
		l := cf.Coverage[i]
		assert.Equal(x, l.Interface())
	}
}

func Test_ParseWithBranchWithNoCov(t *testing.T) {
	ogb := env.GitBlob
	defer func() {
		env.GitBlob = ogb
	}()
	env.GitBlob = func(s string, c *object.Commit) (string, error) {
		return s, nil
	}

	assert := require.New(t)

	formatter := Formatter{
		Path: "./simplecov-branch-with-nocov-example.json",
	}
	rep, err := formatter.Format()
	assert.NoError(err)

	assert.Len(rep.SourceFiles, 7)

	cf := rep.SourceFiles["development/mygem/lib/mygem/wrap.rb"]
	assert.Len(cf.Coverage, 12)
	for i, x := range []interface{}{1, nil, 1, 17, nil, 0, nil, 16, 12, 0, nil, nil} {
		l := cf.Coverage[i]
		assert.Equal(x, l.Interface())
	}
}

func Test_FormatLegacy(t *testing.T) {
	gb := env.GitBlob
	defer func() { env.GitBlob = gb }()
	env.GitBlob = func(s string, c *object.Commit) (string, error) {
		return s, nil
	}

	assert := require.New(t)

	formatter := Formatter{
		Path: "./simplecov-example-legacy-resultset.json",
	}
	rep, err := formatter.Format()
	assert.NoError(err)

	assert.InDelta(97.95, rep.CoveredPercent, 1)

	sf := rep.SourceFiles["development/mygem/lib/mygem/wrap.rb"]
	assert.InDelta(100, sf.CoveredPercent, 1)

	lc := rep.LineCounts
	assert.Equal(lc.Covered, 56)
	assert.Equal(lc.Missed, 1)
	assert.Equal(lc.Total, 57)
}

func Test_Format(t *testing.T) {
	gb := env.GitBlob
	defer func() { env.GitBlob = gb }()
	env.GitBlob = func(s string, c *object.Commit) (string, error) {
		return s, nil
	}

	assert := require.New(t)

	formatter := Formatter{
		Path: "./simplecov-simple-example.json",
	}
	rep, err := formatter.Format()
	assert.NoError(err)

	assert.InDelta(97.95, rep.CoveredPercent, 1)

	sf := rep.SourceFiles["development/mygem/lib/mygem/wrap.rb"]
	assert.InDelta(100, sf.CoveredPercent, 1)

	lc := rep.LineCounts
	assert.Equal(lc.Covered, 56)
	assert.Equal(lc.Missed, 1)
	assert.Equal(lc.Total, 57)
}

func Test_FormatWithBranch(t *testing.T) {
	gb := env.GitBlob
	defer func() { env.GitBlob = gb }()
	env.GitBlob = func(s string, c *object.Commit) (string, error) {
		return s, nil
	}

	assert := require.New(t)

	formatter := Formatter{
		Path: "./simplecov-branch-example.json",
	}
	rep, err := formatter.Format()
	assert.NoError(err)

	assert.InDelta(97.95, rep.CoveredPercent, 1)

	sf := rep.SourceFiles["development/mygem/lib/mygem/wrap.rb"]
	assert.InDelta(100, sf.CoveredPercent, 1)

	lc := rep.LineCounts
	assert.Equal(lc.Covered, 56)
	assert.Equal(lc.Missed, 1)
	assert.Equal(lc.Total, 57)
}

func Test_FormatWithBranchWithNoCovLines(t *testing.T) {
	gb := env.GitBlob
	defer func() { env.GitBlob = gb }()
	env.GitBlob = func(s string, c *object.Commit) (string, error) {
		return s, nil
	}

	assert := require.New(t)

	formatter := Formatter{
		Path: "./simplecov-branch-with-nocov-example.json",
	}
	rep, err := formatter.Format()
	assert.NoError(err)

	assert.InDelta(94.73, rep.CoveredPercent, 1)

	sf := rep.SourceFiles["development/mygem/lib/mygem/wrap.rb"]
	assert.InDelta(71.42, sf.CoveredPercent, 1)

	lc := rep.LineCounts
	assert.Equal(lc.Covered, 54)
	assert.Equal(lc.Missed, 3)
	assert.Equal(lc.Total, 57)
}

func Test_Format_Merged(t *testing.T) {
	gb := env.GitBlob
	defer func() { env.GitBlob = gb }()
	env.GitBlob = func(s string, c *object.Commit) (string, error) {
		return s, nil
	}

	assert := require.New(t)

	formatter := Formatter{
		Path: "./simplecov-merged-legacy-resultset.json",
	}

	rep, err := formatter.Format()
	assert.NoError(err)

	assert.InDelta(75.0, rep.CoveredPercent, 1)
	assert.Len(rep.SourceFiles, 1)

	sf := rep.SourceFiles["/home/patrick/code/codeclimate/ruby-test-reporter/lib/codeclimate-test-reporter.rb"]
	assert.InDelta(75.0, sf.CoveredPercent, 1)

	lc := rep.LineCounts
	assert.Equal(lc.Covered, 6)
	assert.Equal(lc.Missed, 2)
	assert.Equal(lc.Total, 8)
}

func Test_ParseWithGroups(t *testing.T) {
	ogb := env.GitBlob
	defer func() {
		env.GitBlob = ogb
	}()
	env.GitBlob = func(s string, c *object.Commit) (string, error) {
		return s, nil
	}

	assert := require.New(t)

	formatter := Formatter{
		Path: "./simplecov-with-groups.json",
	}
	rep, err := formatter.Format()
	assert.NoError(err)

	assert.Len(rep.SourceFiles, 7)

	cf := rep.SourceFiles["development/mygem/lib/mygem/wrap.rb"]
	assert.Len(cf.Coverage, 10)
	for i, x := range []interface{}{1, nil, 1, 17, 20, 16, 16, 12, nil, nil} {
		l := cf.Coverage[i]
		assert.Equal(x, l.Interface())
	}
}
