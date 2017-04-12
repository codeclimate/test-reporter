package upload

import (
	"time"

	"github.com/codeclimate/test-reporter/formatters"
)

type Data struct {
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

func newData(rep formatters.Report) Data {
	data := Data{
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
	}
	return data
}
