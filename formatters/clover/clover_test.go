package clover

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
	r.Len(rep.SourceFiles, 13)

	sf := rep.SourceFiles["/Users/markbates/Dropbox/development/php-test-reporter/src/TestReporter/Entity/CiInfo.php"]
	r.InDelta(91.78, sf.CoveredPercent, 1)
	r.Len(sf.Coverage, 194)
	r.False(sf.Coverage[51].Valid)
	r.True(sf.Coverage[54].Valid)
	r.Equal(4, sf.Coverage[53].Int)
	r.Equal(0, sf.Coverage[55].Int)
}

func Test_Parse_Without_Package(t *testing.T) {
	gb := env.GitBlob
	defer func() { env.GitBlob = gb }()
	env.GitBlob = func(s string, c *object.Commit) (string, error) {
		return s, nil
	}

	r := require.New(t)

	f := &Formatter{Path: "./example_without_package.xml"}
	rep, err := f.Format()
	r.NoError(err)
	r.Len(rep.SourceFiles, 4)

	sf := rep.SourceFiles["/Users/markbates/Dropbox/development/php-test-reporter/src/ConsoleCommands/SelfUpdateCommand.php"]
	r.InDelta(15.2, sf.CoveredPercent, 1)
	r.Len(sf.Coverage, 80)
	r.False(sf.Coverage[2].Valid)
	r.True(sf.Coverage[43].Valid)
	r.Equal(0, sf.Coverage[38].Int)
	r.Equal(5, sf.Coverage[62].Int)
}
