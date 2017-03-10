package formatters

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/codeclimate/test-reporter/env"
	"github.com/gobuffalo/envy"
)

type Report struct {
	CIService       env.Environment `json:"ci_service"`
	Environment     ccEnvironment   `json:"environment"`
	Git             ccGit           `json:"git"`
	CoveredPercent  float64         `json:"covered_percent"`
	CoveredStrength int             `json:"covered_strength"`
	LineCounts      LineCounts      `json:"line_counts"`
	SourceFiles     []SourceFile    `json:"source_files"`
	RepoToken       string          `json:"repo_token"`
	totalCoverage   float64
}

type ccGit struct {
	Branch      string `json:"branch" structs:"branch"`
	Head        string `json:"head" structs:"head"`
	CommittedAt int    `json:"committed_at" structs:"committed_at"`
}

type ccEnvironment struct {
	GemVersion      string `json:"gem_version"`
	PackageVersion  string `json:"package_version"`
	PWD             string `json:"pwd"`
	RailsRoot       string `json:"rails_root"`
	ReporterVersion string `json:"reporter_version"`
	SimplecovRoot   string `json:"simplecov_root"`
}

func newCCEnvironment() ccEnvironment {
	cc := ccEnvironment{
		RailsRoot: envy.Get("RAILS_ROOT", ""),
	}

	pwd, _ := os.Getwd()
	cc.PWD = pwd

	cmd := exec.Command("gem", "--version")
	out, err := cmd.Output()
	if err == nil {
		cc.GemVersion = strings.TrimSpace(string(out))
	}

	return cc
}

func NewReport() (Report, error) {
	rep := Report{
		SourceFiles: []SourceFile{},
		LineCounts:  LineCounts{},
		Environment: newCCEnvironment(),
	}

	e, err := env.New()
	if err != nil {
		return rep, err
	}
	rep.CIService = e
	rep.Git = ccGit{
		Branch:      e.Git.Branch,
		Head:        e.Git.CommitSHA,
		CommittedAt: e.Git.CommittedAt,
	}
	rep.RepoToken = e.RepoToken

	return rep, nil
}

func (rep *Report) AddSourceFile(sf SourceFile) {
	rep.SourceFiles = append(rep.SourceFiles, sf)
	rep.LineCounts.Covered += sf.LineCounts.Covered
	rep.LineCounts.Missed += sf.LineCounts.Missed
	rep.LineCounts.Total += sf.LineCounts.Total
	rep.totalCoverage += sf.CoveredPercent
	rep.CoveredPercent = rep.totalCoverage / float64(len(rep.SourceFiles))
}

func (r Report) Save(w io.Writer) error {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
