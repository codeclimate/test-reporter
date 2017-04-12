package upload

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func exampleReport() (formatters.Report, error) {
	rep := formatters.Report{}
	f, err := os.Open("../examples/codeclimate.json")
	if err != nil {
		return rep, errors.WithStack(err)
	}

	rep, err = formatters.NewReport()
	if err != nil {
		return rep, errors.WithStack(err)
	}
	err = json.NewDecoder(f).Decode(&rep)
	if err != nil {
		return rep, errors.WithStack(err)
	}
	return rep, nil
}

func Test_NewTestReport(t *testing.T) {
	r := require.New(t)

	rep, err := exampleReport()
	r.NoError(err)
	data := NewTestReport(rep)
	r.Equal("test_reports", data.Type)

	at := data.Attributes
	r.Equal("master", at.CIBranch)
	r.Equal("3", at.CIBuildIdentifier)
	r.Equal("4", at.CIBuildURL)
	r.Equal("700e63d964e0ca1c22fdb11b806109836ca77365", at.CICommitSha)
	r.Equal(1489695537, at.CICommittedAt)
	r.Equal("travis", at.CIServiceName)
	r.Equal("700e63d964e0ca1c22fdb11b806109836ca77365", at.CommitSha)
	r.Equal(1489695537, at.CommittedAt)
	r.NotZero(at.RunAt)
	r.InDelta(88.92, at.CoveredPercent, 1.0)
	r.Equal(0, at.CoveredStrength)
	r.Equal(rep.LineCounts, at.LineCounts)

	env := data.Attributes.Environment
	r.Equal("2.6.10", env.GemVersion)
	r.Equal("42", env.PackageVersion)
	r.Equal("/go/src/github.com/codeclimate/test-reporter/simplecov-test-reporter", env.PWD)
	r.Equal("/rails", env.RailsRoot)
	r.Equal("1", env.ReporterVersion)
	r.Equal("/scov", env.SimplecovRoot)

	r.Len(data.SourceFiles.SourceFiles, 20)
}
