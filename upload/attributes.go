package upload

import "github.com/codeclimate/test-reporter/formatters"

type Attributes struct {
	CIBranch          string                 `json:"ci_branch"`
	CIBuildIdentifier string                 `json:"ci_build_identifier"`
	CIBuildURL        string                 `json:"ci_build_url"`
	CICommitSha       string                 `json:"ci_commit_sha"`
	CICommittedAt     int                    `json:"ci_committed_at"`
	CIServiceName     string                 `json:"ci_service_name"`
	GitBranch         string                 `json:"git_branch"`
	CommitSha         string                 `json:"commit_sha"`
	CommittedAt       int                    `json:"committed_at"`
	RunAt             int64                  `json:"run_at"`
	CoveredPercent    float64                `json:"covered_percent"`
	CoveredStrength   int                    `json:"covered_strength"`
	Environment       formatters.Environment `json:"environment"`
	LineCounts        formatters.LineCounts  `json:"line_counts"`
}
