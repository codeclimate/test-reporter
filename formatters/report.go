package formatters

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/codeclimate/test-reporter/env"
	"github.com/codeclimate/test-reporter/version"
	"github.com/gobuffalo/envy"
)

type Report struct {
	CIService       ccCIService `json:"ci_service"`
	Environment     Environment `json:"environment"`
	Git             ccGit       `json:"git"`
	CoveredPercent  float64     `json:"covered_percent"`
	CoveredStrength int         `json:"covered_strength"`
	LineCounts      LineCounts  `json:"line_counts"`
	SourceFiles     SourceFiles `json:"source_files"`
	RepoToken       string      `json:"repo_token"`
}

type ccCIService struct {
	Branch          string `json:"branch"`
	BuildIdentifier string `json:"build_identifier"`
	BuildURL        string `json:"build_url"`
	CommitSHA       string `json:"commit_sha"`
	CommittedAt     int    `json:"committed_at"`
	Name            string `json:"name"`
}

type ccGit struct {
	Branch      string `json:"branch" structs:"branch"`
	Head        string `json:"head" structs:"head"`
	CommittedAt int    `json:"committed_at" structs:"committed_at"`
}

type Environment struct {
	GemVersion      string `json:"gem_version"`
	PackageVersion  string `json:"package_version"`
	PWD             string `json:"pwd"`
	RailsRoot       string `json:"rails_root"`
	ReporterVersion string `json:"reporter_version"`
	SimplecovRoot   string `json:"simplecov_root"`
}

func newEnvironment() Environment {
	cc := Environment{
		RailsRoot:       envy.Get("RAILS_ROOT", ""),
		ReporterVersion: version.Version,
	}

	pwd, _ := os.Getwd()
	cc.PWD = pwd

	_, err := exec.LookPath("gem")
	if err == nil {
		cmd := exec.Command("gem", "--version")
		out, err := cmd.Output()
		if err == nil {
			cc.GemVersion = strings.TrimSpace(string(out))
		}
	}

	return cc
}

func NewReport() (Report, error) {
	rep := Report{
		SourceFiles: SourceFiles{},
		LineCounts:  LineCounts{},
		Environment: newEnvironment(),
	}

	e, err := env.New()
	if err != nil {
		return rep, err
	}

	rep.Git = ccGit{
		Branch:      e.Git.Branch,
		Head:        e.Git.CommitSHA,
		CommittedAt: e.Git.CommittedAt,
	}

	rep.CIService = ccCIService{
		Branch:          e.Git.Branch,
		BuildURL:        e.CI.BuildURL,
		BuildIdentifier: e.CI.BuildID,
		CommitSHA:       e.Git.CommitSHA,
		CommittedAt:     e.Git.CommittedAt,
		Name:            e.CI.Name,
	}

	rep.RepoToken = e.RepoToken

	return rep, nil
}

func (a *Report) Merge(reps ...*Report) error {
	for _, r := range reps {
		if a.Git.Head != r.Git.Head {
			return errors.New("git heads do not match")
		}
		for _, sf := range r.SourceFiles {
			err := a.AddSourceFile(sf)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (rep *Report) AddSourceFile(sf SourceFile) error {
	var err error
	if s, ok := rep.SourceFiles[sf.Name]; ok {
		sf, err = s.Merge(sf)
		if err != nil {
			return err
		}
	}
	sf.CalcLineCounts()
	rep.SourceFiles[sf.Name] = sf

	lc := LineCounts{}
	for _, s := range rep.SourceFiles {
		lc.Covered += s.LineCounts.Covered
		lc.Missed += s.LineCounts.Missed
		lc.Total += s.LineCounts.Total
	}
	rep.LineCounts = lc
	rep.CoveredPercent = rep.LineCounts.CoveredPercent()
	return nil
}

func (r Report) Save(w io.Writer) error {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
