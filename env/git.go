package env

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Git struct {
	Branch      string `json:"branch" structs:"branch"`
	CommitSHA   string `json:"commit_sha" structs:"commit_sha"`
	CommittedAt int    `json:"committed_at" structs:"committed_at"`
}

func (g Git) String() string {
	out := &bytes.Buffer{}
	out.WriteString("GIT_BRANCH=")
	out.WriteString(g.Branch)
	out.WriteString("\nGIT_COMMIT_SHA=")
	out.WriteString(g.CommitSHA)
	out.WriteString("\nGIT_COMMITTED_AT=")
	out.WriteString(fmt.Sprint(g.CommittedAt))
	return out.String()
}

func findGitInfo() (Git, error) {
	_, err := exec.LookPath("git")
	if err != nil {
		// git isn't present, so load from ENV vars:
		return loadGitFromENV()
	}

	g := Git{}

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return g, errors.WithStack(err)
	}
	g.Branch = strings.TrimSpace(string(out))

	cmd = exec.Command("git", "log", "-1", "--pretty=format:%H")
	out, err = cmd.Output()
	if err != nil {
		return g, errors.WithStack(err)
	}
	g.CommitSHA = strings.TrimSpace(string(out))

	cmd = exec.Command("git", "log", "-1", "--pretty=format:%ct")
	out, err = cmd.Output()
	if err != nil {
		return g, errors.WithStack(err)
	}
	g.CommittedAt, err = strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return g, errors.WithStack(err)
	}
	return g, nil
}

func GitSHA(path string) (string, error) {
	args := []string{"log", "-1", "--follow", "--pretty=format:%H"}
	if path != "" {
		args = append(args, path)
	}
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return strings.TrimSpace(string(out)), nil
}

var GitBlob = func(path string) (string, error) {
	sha, err := GitSHA(path)
	if err != nil {
		return "", errors.WithStack(err)
	}
	cmd := exec.Command("git", "ls-tree", sha, "--", path)
	out, err := cmd.Output()
	if err != nil {
		return "", errors.WithStack(err)
	}
	res := strings.TrimSpace(string(out))
	matches := blobRegex.FindStringSubmatch(res)
	if len(matches) == 0 {
		return "", errors.Errorf("could not find blob id for file %s in %s", path, res)
	}
	return matches[1], nil
}

func loadGitFromENV() (Git, error) {
	g := Git{
		Branch:    findVar(gitBranchVars),
		CommitSHA: findVar(gitCommitShaVars),
	}
	var err error
	g.CommittedAt, err = strconv.Atoi(findVar(gitCommittedAtVars))
	return g, err
}

var gitBranchVars = []string{"GIT_BRANCH", "APPVEYOR_REPO_BRANCH", "BRANCH_NAME", "BUILDKITE_BRANCH", "CIRCLE_BRANCH", "CI_BRANCH", "CI_BUILD_REF_NAME", "TRAVIS_BRANCH", "WERCKER_GIT_BRANCH"}

var gitCommitShaVars = []string{"GIT_COMMIT_SHA", "APPVEYOR_REPO_COMMIT", "BUILDKITE_COMMIT", "CIRCLE_SHA1", "CI_BUILD_REF", "CI_BUILD_SHA", "CI_COMMIT", "CI_COMMIT_ID", "GIT_COMMIT", "WERCKER_GIT_COMMIT"}

var gitCommittedAtVars = []string{"GIT_COMMITTED_AT", "GIT_COMMITED_AT", "CI_COMMITTED_AT", "CI_COMMITED_AT"}

var blobRegex = regexp.MustCompile(`^\d.+\s+blob\s(\w+)`)
