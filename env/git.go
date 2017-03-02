package env

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/gobuffalo/envy"
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
		Branch:      envy.Get("GIT_BRANCH", ""),
		CommitSHA:   envy.Get("GIT_COMMIT_SHA", ""),
		CommittedAt: envy.Get("GIT_COMMITTED_AT", ""),
	}, nil
}
