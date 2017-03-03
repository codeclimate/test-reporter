package env

import (
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/stretchr/testify/require"
)

func Test_FindGitInfo(t *testing.T) {
	r := require.New(t)
	g, err := findGitInfo()
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
		r.Equal(g.Branch, "master")
		r.Equal(g.CommitSHA, "a12345")
		r.Equal(g.CommittedAt, "12:45")
	})
}

func Test_loadGitFromENV_Alt_Vars(t *testing.T) {
	r := require.New(t)
	envy.Temp(func() {
		envy.Set("CIRCLE_BRANCH", "circle")
		envy.Set("WERCKER_GIT_COMMIT", "b12345")
		envy.Set("CI_COMMITED_AT", "13:45")
		g, err := loadGitFromENV()
		r.NoError(err)
		r.Equal(g.Branch, "circle")
		r.Equal(g.CommitSHA, "b12345")
		r.Equal(g.CommittedAt, "13:45")
	})
}

func Test_Git_String(t *testing.T) {
	r := require.New(t)
	g := git{
		Branch:      "master",
		CommitSHA:   "a12345",
		CommittedAt: "12:45",
	}
	exp := `GIT_BRANCH=master
GIT_COMMIT_SHA=a12345
GIT_COMMITTED_AT=12:45`
	r.Equal(exp, g.String())
}
