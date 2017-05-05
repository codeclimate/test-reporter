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
	err := f.Parse()
	r.NoError(err)
	r.Len(f.SourceFiles, 12)

	sf := f.SourceFiles[10]
	r.Equal("/Users/markbates/Dropbox/development/php-test-reporter/src/TestReporter/Entity/CiInfo.php", sf.Name)
	r.InDelta(91.78, sf.CoveredPercent, 1)
	r.Len(sf.Coverage, 194)
	r.False(sf.Coverage[51].Valid)
	r.True(sf.Coverage[54].Valid)
	r.Equal(4, sf.Coverage[53].Int)
	r.Equal(0, sf.Coverage[55].Int)
}
