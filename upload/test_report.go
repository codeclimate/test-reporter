package upload

import (
	"time"

	"github.com/codeclimate/test-reporter/formatters"
)

func NewTestReport(rep formatters.Report) *TestReport {
	tr := &TestReport{
		Type: "test_reports",
		Attributes: Attributes{
			CIBranch:          rep.Git.Branch,
			CIBuildIdentifier: rep.CIService.BuildIdentifier,
			CIBuildURL:        rep.CIService.BuildURL,
			CICommitSha:       rep.Git.Head,
			CIServiceName:     rep.CIService.Name,
			CICommittedAt:     rep.CIService.CommittedAt,
			GitBranch:         rep.Git.Branch,
			CommitSha:         rep.Git.Head,
			CommittedAt:       rep.Git.CommittedAt,
			RunAt:             time.Now().Unix(),
			CoveredPercent:    rep.CoveredPercent,
			CoveredStrength:   rep.CoveredStrength,
			LineCounts:        rep.LineCounts,
			Environment:       rep.Environment,
		},
		SourceFiles: []SourceFile{},
	}
	for _, sf := range rep.SourceFiles {
		tr.SourceFiles = append(tr.SourceFiles, SourceFile{
			Type:            "test_file_reports",
			BlobID:          sf.BlobID,
			Coverage:        sf.Coverage,
			CoveredPercent:  sf.CoveredPercent,
			CoveredStrength: sf.CoveredStrength,
			LineCounts:      sf.LineCounts,
			Path:            sf.Name,
		})
	}
	return tr
}

type TestReport struct {
	Type        string       `json:"type"`
	Attributes  Attributes   `json:"attributes"`
	SourceFiles []SourceFile `json:"-"`
}
