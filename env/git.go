package env

import (
	"bytes"
	"os/exec"
	"strings"
)

type Git struct {
	Branch      string
	CommitSHA   string
	CommittedAt string
}

func (g Git) String() string {
	out := &bytes.Buffer{}
	out.WriteString("GIT_BRANCH=")
	out.WriteString(g.Branch)
	out.WriteString("\nGIT_COMMIT_SHA=")
	out.WriteString(g.CommitSHA)
	out.WriteString("\nGIT_COMMITTED_AT=")
	out.WriteString(g.CommittedAt)
	return out.String()
}

func FindGitInfo() (Git, error) {
	_, err := exec.LookPath("git")
	if err != nil {
		// git isn't present, so load from ENV vars:
		return loadGitFromENV()
	}

	g := Git{}

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return g, err
	}
	g.Branch = strings.TrimSpace(string(out))

	cmd = exec.Command("git", "log", "-1", "--pretty=format:'%H'")
	out, err = cmd.Output()
	if err != nil {
		return g, err
	}
	g.CommitSHA = strings.TrimSpace(string(out))

	cmd = exec.Command("git", "log", "-1", "--pretty=format:'%ct'")
	out, err = cmd.Output()
	if err != nil {
		return g, err
	}
	g.CommittedAt = strings.TrimSpace(string(out))
	return g, nil
}

func loadGitFromENV() (Git, error) {
	// TODO: find via other variables:
	return Git{
		Branch:      findVar(gitBranchVars),
		CommitSHA:   findVar(gitCommitShaVars),
		CommittedAt: findVar(gitCommittedAtVars),
	}, nil
}

var gitBranchVars = []string{"GIT_BRANCH", "APPVEYOR_REPO_BRANCH", "BRANCH_NAME", "BUILDKITE_BRANCH", "CIRCLE_BRANCH", "CI_BRANCH", "CI_BUILD_REF_NAME", "TRAVIS_BRANCH", "WERCKER_GIT_BRANCH"}

var gitCommitShaVars = []string{"GIT_COMMIT_SHA", "APPVEYOR_REPO_COMMIT", "BUILDKITE_COMMIT", "CIRCLE_SHA1", "CI_BUILD_REF", "CI_BUILD_SHA", "CI_COMMIT", "CI_COMMIT_ID", "GIT_COMMIT", "WERCKER_GIT_COMMIT"}

var gitCommittedAtVars = []string{"GIT_COMMITTED_AT", "GIT_COMMITED_AT", "CI_COMMITTED_AT", "CI_COMMITED_AT"}
