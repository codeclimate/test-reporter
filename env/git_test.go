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
		envy.Set("GIT_COMMITTED_AT", "1234")
		g, err := findGitInfo()
		r.NoError(err)
		r.Equal(g.Branch, "master")
		r.Equal(g.CommitSHA, "a12345")
		r.Equal(g.CommittedAt, 1234)
	})
}

func Test_loadGitFromENVOrGit(t *testing.T) {
	r := require.New(t)
	envy.Temp(func() {
		envy.Set("GIT_BRANCH", "master")
		envy.Set("GIT_COMMIT_SHA", "a12345")
		g, err := findGitInfo()
		r.NoError(err)
		r.Equal(g.Branch, "master")
		r.Equal(g.CommitSHA, "a12345")
		r.NotZero(g.CommittedAt)
	})
}

func Test_loadGitFromENV_GitHub_Vars(t *testing.T) {
	r := require.New(t)
	envy.Temp(func() {
		envy.MustSet("GITHUB_REF", "master")
		envy.MustSet("GITHUB_SHA", "a12345")
		g, err := findGitInfo()
		r.NoError(err)
		r.Equal(g.Branch, "master")
		r.Equal(g.CommitSHA, "a12345")
		r.NotZero(g.CommittedAt)
	})
}

func Test_loadGitFromENV_Alt_Vars(t *testing.T) {
	r := require.New(t)
	envy.Temp(func() {
		envy.Set("CIRCLE_BRANCH", "circle")
		envy.Set("CIRCLE_SHA1", "b12345")
		envy.Set("CI_COMMITED_AT", "1345")
		g, err := findGitInfo()
		r.NoError(err)
		r.Equal(g.Branch, "circle")
		r.Equal(g.CommitSHA, "b12345")
		r.Equal(g.CommittedAt, 1345)
	})
}

func Test_Git_String(t *testing.T) {
	r := require.New(t)
	g := Git{
		Branch:      "master",
		CommitSHA:   "a12345",
		CommittedAt: 1234,
	}
	exp := `GIT_BRANCH=master
GIT_COMMIT_SHA=a12345
GIT_COMMITTED_AT=1234`
	r.Equal(exp, g.String())
}

func Test_loadGitFromENV_CI_TIMESTAMP(t *testing.T) {
	r := require.New(t)
	envy.Temp(func() {
		envy.Set("CI_TIMESTAMP", "1345")
		g, err := findGitInfo()
		r.NoError(err)
		r.Equal(g.CommittedAt, 1345)
	})
}
