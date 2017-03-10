package env

import (
	"bytes"
	"encoding/json"

	"github.com/fatih/structs"
	"github.com/gobuffalo/envy"
)

// Environment represent the current testing environment
type Environment struct {
	Git       Git
	CI        ci
	RepoToken string `json:"-"`
}

func (e Environment) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}

	g := structs.Map(e.Git)
	for k, v := range g {
		m[k] = v
	}

	g = structs.Map(e.CI)
	for k, v := range g {
		m[k] = v
	}

	return json.Marshal(m)
}

func (e Environment) String() string {
	out := &bytes.Buffer{}
	out.WriteString(e.Git.String())
	out.WriteString("\n")
	out.WriteString(e.CI.String())
	return out.String()
}

// New environment. If there are problems loading parts of
// the environment an error will be returned. Validation errors
// are not considered an "error" here, but should be checked
// further down the chain, when validation of the environment
// is required.
func New() (Environment, error) {
	e := Environment{
		RepoToken: findVar([]string{"CC_REPO_TOKEN", "REPO_TOKEN"}),
	}
	git, err := findGitInfo()
	if err != nil {
		return e, err
	}
	e.Git = git
	return e, nil
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
