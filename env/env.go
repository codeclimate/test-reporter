package env

import (
	"bytes"

	"github.com/gobuffalo/envy"
)

type CI struct {
	Name     string
	BuildID  string
	BuildURL string
}

func (c CI) String() string {
	out := &bytes.Buffer{}
	// TODO fill in
	return out.String()
}

type Environment struct {
	Git Git
	CI  CI
}

func (e Environment) String() string {
	out := &bytes.Buffer{}
	out.WriteString(e.Git.String())
	out.WriteString(e.CI.String())
	return out.String()
}

func New() (Environment, error) {
	e := Environment{}
	git, err := FindGitInfo()
	if err != nil {
		return e, err
	}
	e.Git = git
	return e, nil
}

func loadCIInfo() CI {
	return CI{}
}

func findVar(names []string) string {
	for _, n := range names {
		v := envy.Get(n, "")
		if v != "" {
			return v
		}
	}
	return ""
}
