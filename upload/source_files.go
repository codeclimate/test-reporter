package upload

import (
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/markbates/pop/nulls"
)

type Coverage []nulls.Int

type SourceFile struct {
	Type            string                `json:"type"`
	BlobID          string                `json:"blob_id"`
	Coverage        Coverage              `json:"coverage"`
	CoveredPercent  float64               `json:"covered_percent"`
	CoveredStrength float64               `json:"covered_strength"`
	LineCounts      formatters.LineCounts `json:"line_counts"`
	Path            string                `json:"path"`
}
