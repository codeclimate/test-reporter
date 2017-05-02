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
	l := Formatter{Path: "./lcov-example.info"}
	err := l.Parse()
	r.NoError(err)

	r.Len(l.SourceFiles, 1)
	sf := l.SourceFiles[0]
	r.Equal("/Users/markbates/Dropbox/development/javascript-test-reporter/formatter.js", sf.Name)
	r.Len(sf.Coverage, 104)
	// for i, c := range sf.Coverage {
	// 	fmt.Printf("%d,%d\n", i+1, c.Int)
	// }
	// r.False(true)
}
