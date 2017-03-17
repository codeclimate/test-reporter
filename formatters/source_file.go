package formatters

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/codeclimate/test-reporter/env"
	"github.com/pkg/errors"
)

type SourceFile struct {
	BlobID          string     `json:"blob_id"`
	Coverage        Coverage   `json:"coverage"`
	CoveredPercent  float64    `json:"covered_percent"`
	CoveredStrength int        `json:"covered_strength"`
	LineCounts      LineCounts `json:"line_counts"`
	Name            string     `json:"name"`
}

func (a SourceFile) Merge(b SourceFile) (SourceFile, error) {
	if len(a.Coverage) != len(b.Coverage) {
		return a, errors.Errorf("coverage length mismatch for %s", a.Name)
	}

	for i, c := range b.Coverage {
		if x, ok := c.(int); ok {
			// the secondary is a number
			if y, ok := a.Coverage[i].(int); ok {
				// the primary is also a number
				a.Coverage[i] = x + y
				continue
			}
			// set to the secondary:
			a.Coverage[i] = x
		}
	}
	a.CalcLineCounts()
	return a, nil
}

func (sf *SourceFile) CalcLineCounts() {
	lc := LineCounts{
		Total: len(sf.Coverage),
	}

	for _, c := range sf.Coverage {
		if _, ok := c.(int); ok {
			lc.Covered++
			continue
		}
		lc.Missed++
	}

	sf.LineCounts = lc
}

func NewSourceFile(name string) SourceFile {
	if pwd, err := os.Getwd(); err == nil {
		pwd := fmt.Sprintf("%s%s", pwd, string(os.PathSeparator))
		name = strings.TrimPrefix(name, pwd)
	}

	sf := SourceFile{Name: name}
	sf.BlobID, _ = env.GitSHA(name)
	return sf
}

type SourceFiles map[string]SourceFile

func (sf SourceFiles) MarshalJSON() ([]byte, error) {
	files := []SourceFile{}
	for _, s := range sf {
		files = append(files, s)
	}
	return json.Marshal(files)
}

func (sf SourceFiles) UnmarshalJSON(text []byte) error {
	files := []SourceFile{}
	err := json.Unmarshal(text, &files)
	if err != nil {
		return err
	}
	for _, f := range files {
		if ff, ok := sf[f.Name]; ok {
			f, err = ff.Merge(f)
			if err != nil {
				return err
			}
		}
		sf[f.Name] = f
	}
	return nil
}
