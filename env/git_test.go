package env

import (
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/stretchr/testify/require"
)

func Test_FindGitInfo(t *testing.T) {
	r := require.New(t)
	g, err := FindGitInfo()
	r.NoError(err)
	r.NotZero(g.Branch)
	r.NotZero(g.CommitSHA)
	r.NotZero(g.CommittedAt)
}

func Test_loadGitFromENV(t *testing.T) {
	r := require.New(t)
	envy.Temp(func() {
		envy.Set("GIT_BRANCH", "master")
		envy.Set("GIT_COMMIT_SHA", "a12345")
		envy.Set("GIT_COMMITTED_AT", "12:45")
		g, err := loadGitFromENV()
		r.NoError(err)
		r.NotZero(g.Branch)
		r.Equal(g.Branch, "master")
		r.NotZero(g.CommitSHA)
		r.Equal(g.CommitSHA, "a12345")
		r.NotZero(g.CommittedAt)
		r.Equal(g.CommittedAt, "12:45")
	})
}

func Test_Git_String(t *testing.T) {
	r := require.New(t)
	g := Git{
		Branch:      "master",
		CommitSHA:   "a12345",
		CommittedAt: "12:45",
	}
	exp := `GIT_BRANCH=master
GIT_COMMIT_SHA=a12345
GIT_COMMITTED_AT=12:45`
	r.Equal(exp, g.String())
}
