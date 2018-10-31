package env

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
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

func GetHead() (*object.Commit, error) {
	r, err := git.PlainOpen(".")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ref, err := r.Head()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return commit, nil
}

func findGitInfo() (Git, error) {
	g := Git{}
	var err error

	g.Branch = findVar(gitBranchVars)
	if g.Branch == "" {
		logrus.Debug("couldn't load branch from ENV, trying git...")
		g.Branch, err = loadFromGit("rev-parse", "--abbrev-ref", "HEAD")
		if err != nil {
			return g, errors.WithStack(err)
		}
	}

	g.CommitSHA = findVar(gitCommitShaVars)
	if g.CommitSHA == "" {
		logrus.Debug("couldn't load commit sha from ENV, trying git...")
		g.CommitSHA, err = loadFromGit("log", "-1", "--pretty=format:%H")
		if err != nil {
			return g, errors.WithStack(err)
		}
	}

	committedAt := findVar(gitCommittedAtVars)
	if committedAt == "" {
		logrus.Debug("couldn't load committed at from ENV, trying git...")
		committedAt, err = loadFromGit("log", "-1", "--pretty=format:%ct")
		if err != nil {
			return g, errors.WithStack(err)
		}
	}
	g.CommittedAt, err = strconv.Atoi(committedAt)
	if err != nil {
		return g, errors.WithStack(err)
	}

	return g, nil
}

func GitSHA(path string) (string, error) {
	args := []string{"log", "-1", "--follow", "--pretty=format:%H"}
	if path != "" {
		if pwd, err := os.Getwd(); err == nil {
			path = strings.TrimPrefix(path, pwd)
			path = filepath.Join(".", path)
		}
		args = append(args, path)
	}
	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return strings.TrimSpace(string(out)), nil
}

var GitBlob = func(path string, commit *object.Commit) (string, error) {
	if commit != nil {
		file, err := commit.File(path)
		if err == nil {
			logrus.Debugf("getting git blob_id for source file %s", path)

			blob := strings.TrimSpace(file.Hash.String())
			return blob, nil
		}
	}

	blob, err := fallbackBlob(path)

	if err != nil {
		return "", errors.WithStack(err)
	}

	return blob, nil
}

func fallbackBlob(path string) (string, error) {
	logrus.Debugf("getting fallback blob_id for source file %s", path)
	file, err := ioutil.ReadFile(path)

	if err != nil {
		logrus.Errorf("failed to read file %s\n%s", path, err)
		return "", errors.WithStack(err)
	}

	hash := plumbing.ComputeHash(plumbing.BlobObject, []byte(file))
	res := hash.String()

	return res, nil
}

func loadFromGit(gitArgs ...string) (string, error) {
	cmd := exec.Command("git", gitArgs...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return "", errors.WithStack(err)
	}

	return strings.TrimSpace(string(out)), nil
}

var gitBranchVars = []string{"GIT_BRANCH", "APPVEYOR_REPO_BRANCH", "BRANCH_NAME", "BUILDKITE_BRANCH", "CIRCLE_BRANCH", "CI_BRANCH", "CI_BUILD_REF_NAME", "DRONE_BRANCH", "HEROKU_TEST_RUN_BRANCH", "TRAVIS_PULL_REQUEST_BRANCH", "TRAVIS_BRANCH", "WERCKER_GIT_BRANCH", "CI_COMMIT_REF_NAME"}

var gitCommitShaVars = []string{"GIT_COMMIT_SHA", "APPVEYOR_REPO_COMMIT", "BUILDKITE_COMMIT", "CIRCLE_SHA1", "CI_BUILD_REF", "CI_BUILD_SHA", "CI_COMMIT", "CI_COMMIT_ID", "DRONE_COMMIT", "GIT_COMMIT", "HEROKU_TEST_RUN_COMMIT_VERSION", "WERCKER_GIT_COMMIT", "TRAVIS_PULL_REQUEST_SHA", "TRAVIS_COMMIT", "CI_COMMIT_SHA"}

var gitCommittedAtVars = []string{"GIT_COMMITTED_AT", "GIT_COMMITED_AT", "CI_COMMITTED_AT", "CI_COMMITED_AT", "CI_TIMESTAMP"}

var blobRegex = regexp.MustCompile(`^\d.+\s+blob\s(\w+)`)
