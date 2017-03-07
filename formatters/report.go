package formatters

import (
	"encoding/json"
	"io"

	"github.com/codeclimate/test-reporter/env"
)

type Report struct {
	CIService       env.Environment `json:"ci_service"`
	Git             env.Git         `json:"git"`
	CoveredPercent  float64         `json:"covered_percent"`
	CoveredStrength int             `json:"covered_strength"`
	LineCounts      LineCounts      `json:"line_counts"`
	SourceFiles     []SourceFile    `json:"source_files"`
}

func (r Report) Save(w io.Writer) error {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
