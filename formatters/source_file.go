package formatters

import "github.com/codeclimate/test-reporter/env"

type SourceFile struct {
	BlobID          string        `json:"blob_id"`
	Coverage        []interface{} `json:"coverage"`
	CoveredPercent  float64       `json:"covered_percent"`
	CoveredStrength int           `json:"covered_strength"`
	LineCounts      LineCounts    `json:"line_counts"`
	Name            string        `json:"name"`
}

func NewSourceFile(name string) SourceFile {
	sf := SourceFile{Name: name}
	sf.BlobID, _ = env.GitSHA(name)
	return sf
}
